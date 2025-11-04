package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aitmiloud/ngxtui/internal/model"
)

const (
	sitesAvailableDir = "/etc/nginx/sites-available"
	sitesEnabledDir   = "/etc/nginx/sites-enabled"
)

// Service handles NGINX operations
type Service struct{}

// New creates a new NGINX service
func New() *Service {
	return &Service{}
}

// ListSites returns a list of all NGINX sites
func (s *Service) ListSites() ([]model.Site, error) {
	sites := []model.Site{}

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

		// Check if site is enabled
		enabledPath := filepath.Join(sitesEnabledDir, siteName)
		_, err := os.Stat(enabledPath)
		enabled := err == nil

		sites = append(sites, model.Site{
			Name:    siteName,
			Enabled: enabled,
			Port:    "80",
			SSL:     false,
			Uptime:  "N/A",
		})
	}

	return sites, nil
}

// EnableSite enables an NGINX site
func (s *Service) EnableSite(siteName string) error {
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
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config test failed: %s", string(output))
	}
	return nil
}

// Reload reloads the NGINX service
func (s *Service) Reload() error {
	cmd := exec.Command("systemctl", "reload", "nginx")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}
	return nil
}

// GetAccessLogs returns the last N lines of the access log
func (s *Service) GetAccessLogs(lines int) (string, error) {
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), "/var/log/nginx/access.log")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read access log: %w", err)
	}
	return string(output), nil
}

// GetErrorLogs returns the last N lines of the error log
func (s *Service) GetErrorLogs(lines int) (string, error) {
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), "/var/log/nginx/error.log")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read error log: %w", err)
	}
	return string(output), nil
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
