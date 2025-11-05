package nginx

import (
	"os"
	"path/filepath"
	"strings"

	crossplane "github.com/nginxinc/nginx-go-crossplane"
)

// GetListeningPorts returns all ports that NGINX is configured to listen on
func (s *Service) GetListeningPorts() ([]string, error) {
	// First, try to detect Docker NGINX (with caching)
	if IsDockerAvailable() {
		containerID, err := getCachedContainerID()
		if err == nil {
			// NGINX is running in Docker
			ports, err := GetDockerNginxPorts(containerID)
			if err == nil && len(ports) > 0 {
				return ports, nil
			}
		}
	}

	// Fall back to parsing config files (native NGINX)
	portsMap := make(map[string]bool)

	// Parse main nginx.conf
	payload, err := crossplane.Parse(nginxConfPath, &crossplane.ParseOptions{
		SingleFile:         false,
		StopParsingOnError: false,
	})
	if err != nil {
		return nil, err
	}

	// Extract ports from all server blocks
	for _, config := range payload.Config {
		for _, directive := range config.Parsed {
			if directive.Directive == "server" {
				extractPortsFromBlock(directive.Block, portsMap)
			}
			// Also check http blocks
			if directive.Directive == "http" && directive.Block != nil {
				for _, httpDir := range directive.Block {
					if httpDir.Directive == "server" {
						extractPortsFromBlock(httpDir.Block, portsMap)
					}
				}
			}
		}
	}

	// Parse sites-available
	entries, err := os.ReadDir(sitesAvailableDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			sitePath := filepath.Join(sitesAvailableDir, entry.Name())
			sitePayload, err := crossplane.Parse(sitePath, &crossplane.ParseOptions{
				SingleFile: true,
			})
			if err != nil {
				continue
			}

			for _, config := range sitePayload.Config {
				for _, directive := range config.Parsed {
					if directive.Directive == "server" {
						extractPortsFromBlock(directive.Block, portsMap)
					}
				}
			}
		}
	}

	// Convert map to slice
	var ports []string
	for port := range portsMap {
		ports = append(ports, port)
	}

	// If no ports found, default to common ports
	if len(ports) == 0 {
		ports = []string{"80", "443"}
	}

	return ports, nil
}

// extractPortsFromBlock extracts port numbers from listen directives
func extractPortsFromBlock(block crossplane.Directives, portsMap map[string]bool) {
	for _, directive := range block {
		if directive.Directive == "listen" && len(directive.Args) > 0 {
			port := directive.Args[0]
			// Remove options like "default_server", "ssl", etc.
			parts := strings.Fields(port)
			if len(parts) > 0 {
				portStr := parts[0]
				// Handle formats like "0.0.0.0:8080" or "[::]:8080"
				if strings.Contains(portStr, ":") {
					portStr = portStr[strings.LastIndex(portStr, ":")+1:]
				}
				// Remove any trailing semicolons or options
				portStr = strings.TrimRight(portStr, ";")
				portsMap[portStr] = true
			}
		}
	}
}
