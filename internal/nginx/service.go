package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aitmiloud/ngxtui/internal/model"
	crossplane "github.com/nginxinc/nginx-go-crossplane"
)

const (
	sitesAvailableDir = "/etc/nginx/sites-available"
	sitesEnabledDir   = "/etc/nginx/sites-enabled"
	nginxConfPath     = "/etc/nginx/nginx.conf"
)

// Service handles NGINX operations using crossplane for real config parsing
type Service struct {
	payload *crossplane.Payload
}

// New creates a new NGINX service
func New() *Service {
	return &Service{}
}

// parseConfig parses the NGINX configuration using crossplane
func (s *Service) parseConfig() error {
	payload, err := crossplane.Parse(nginxConfPath, &crossplane.ParseOptions{
		SingleFile:         false,
		StopParsingOnError: false,
	})
	if err != nil {
		return fmt.Errorf("failed to parse nginx config: %w", err)
	}
	s.payload = payload
	return nil
}

// ListSites returns a list of all NGINX sites with real config parsing
func (s *Service) ListSites() ([]model.Site, error) {
	// First, try to detect Docker NGINX (with caching)
	if IsDockerAvailable() {
		containerID, err := getCachedContainerID()
		if err == nil {
			// NGINX is running in Docker
			sites, err := GetDockerNginxSites(containerID)
			if err == nil && len(sites) > 0 {
				return sites, nil
			}
		}
	}

	// Fall back to reading config files (native NGINX)
	sites := []model.Site{}

	// Read sites-available directory
	entries, err := os.ReadDir(sitesAvailableDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read sites-available: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		siteName := entry.Name()
		if siteName == "default" {
			continue
		}

		// Parse the site configuration
		sitePath := filepath.Join(sitesAvailableDir, siteName)
		siteInfo, err := s.parseSiteConfig(sitePath)
		if err != nil {
			// If parsing fails, create basic site info
			siteInfo = &model.Site{
				Name:   siteName,
				Port:   "unknown",
				SSL:    false,
				Uptime: "N/A",
			}
		}

		// Check if site is enabled
		enabledPath := filepath.Join(sitesEnabledDir, siteName)
		_, err = os.Stat(enabledPath)
		siteInfo.Enabled = err == nil
		siteInfo.Name = siteName

		// Get uptime if enabled
		if siteInfo.Enabled {
			siteInfo.Uptime = s.getSiteUptime(siteName)
		} else {
			siteInfo.Uptime = "Disabled"
		}

		sites = append(sites, *siteInfo)
	}

	return sites, nil
}

// parseSiteConfig parses a site configuration file and extracts key information
func (s *Service) parseSiteConfig(configPath string) (*model.Site, error) {
	payload, err := crossplane.Parse(configPath, &crossplane.ParseOptions{
		SingleFile:         true,
		StopParsingOnError: false,
	})
	if err != nil {
		return nil, err
	}

	site := &model.Site{
		Port: "80",
		SSL:  false,
	}

	// Parse the configuration to extract server blocks
	if len(payload.Config) > 0 {
		for _, directive := range payload.Config[0].Parsed {
			if directive.Directive == "server" {
				s.parseServerBlock(directive.Block, site)
			}
		}
	}

	return site, nil
}

// parseServerBlock extracts information from a server block
func (s *Service) parseServerBlock(block crossplane.Directives, site *model.Site) {
	for _, directive := range block {
		switch directive.Directive {
		case "listen":
			if len(directive.Args) > 0 {
				port := directive.Args[0]
				// Remove any options like "default_server"
				parts := strings.Fields(port)
				if len(parts) > 0 {
					site.Port = parts[0]
				}
				// Check for SSL
				for _, arg := range directive.Args {
					if arg == "ssl" {
						site.SSL = true
					}
				}
			}
		case "ssl_certificate":
			site.SSL = true
		case "server_name":
			// Could store server names if needed
		}
	}
}

// getSiteUptime calculates how long a site has been enabled
func (s *Service) getSiteUptime(siteName string) string {
	enabledPath := filepath.Join(sitesEnabledDir, siteName)
	info, err := os.Stat(enabledPath)
	if err != nil {
		return "N/A"
	}

	// Get the modification time of the symlink
	modTime := info.ModTime()
	duration := time.Since(modTime)

	// Format duration
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return "< 1h"
}

// EnableSite enables an NGINX site
func (s *Service) EnableSite(siteName string) error {
	// Check if Docker NGINX
	if IsDockerAvailable() {
		containerID, err := DetectDockerNginx()
		if err == nil {
			return DockerEnableSite(containerID, siteName)
		}
	}
	availablePath := filepath.Join(sitesAvailableDir, siteName)
	enabledPath := filepath.Join(sitesEnabledDir, siteName)

	// Check if site exists
	if _, err := os.Stat(availablePath); os.IsNotExist(err) {
		return fmt.Errorf("site %s does not exist", siteName)
	}

	// Create symlink
	if err := os.Symlink(availablePath, enabledPath); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("failed to enable site: %w", err)
		}
	}

	return nil
}

// DisableSite disables an NGINX site
func (s *Service) DisableSite(siteName string) error {
	// Check if Docker NGINX
	if IsDockerAvailable() {
		containerID, err := DetectDockerNginx()
		if err == nil {
			return DockerDisableSite(containerID, siteName)
		}
	}
	enabledPath := filepath.Join(sitesEnabledDir, siteName)

	// Remove symlink
	if err := os.Remove(enabledPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to disable site: %w", err)
		}
	}

	return nil
}

// TestConfig tests the NGINX configuration
func (s *Service) TestConfig() error {
	// Check if Docker NGINX
	if IsDockerAvailable() {
		containerID, err := DetectDockerNginx()
		if err == nil {
			return DockerTestConfig(containerID)
		}
	}
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config test failed: %s", string(output))
	}
	return nil
}

// Reload reloads the NGINX configuration
func (s *Service) Reload() error {
	// Check if Docker NGINX
	if IsDockerAvailable() {
		containerID, err := DetectDockerNginx()
		if err == nil {
			return DockerReload(containerID)
		}
	}
	cmd := exec.Command("systemctl", "reload", "nginx")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}
	return nil
}

// ParseLogLine parses an NGINX access log line and returns a styled version
func ParseLogLine(line string) string {
	if strings.Contains(line, " 200 ") {
		return line // Success - no special styling needed
	} else if strings.Contains(line, " 404 ") {
		return line // Not found
	} else if strings.Contains(line, " 500 ") || strings.Contains(line, " 502 ") || strings.Contains(line, " 503 ") {
		return line // Server error
	}
	return line
}
