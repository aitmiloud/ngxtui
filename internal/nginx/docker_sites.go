package nginx

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aitmiloud/ngxtui/internal/model"
)

// GetDockerNginxSites gets sites from Docker NGINX container
func GetDockerNginxSites(containerID string) ([]model.Site, error) {
	var sites []model.Site

	// Get Docker container uptime
	containerUptime := getDockerContainerUptime(containerID)

	// Get nginx configuration from container
	cmd := exec.Command("docker", "exec", containerID, "nginx", "-T")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get nginx config from container: %w", err)
	}

	configStr := string(output)

	// Parse server blocks from the config
	serverBlocks := parseServerBlocks(configStr)

	for i, block := range serverBlocks {
		site := model.Site{
			Name:    fmt.Sprintf("server-%d", i+1),
			Enabled: true, // All servers in running config are enabled
			Port:    "80",
			SSL:     false,
			Uptime:  containerUptime,
		}

		// Extract server_name
		if serverName := extractDirective(block, "server_name"); serverName != "" {
			site.Name = serverName
		}

		// Extract listen port
		if listenPort := extractDirective(block, "listen"); listenPort != "" {
			// Parse port from "listen 80;" or "listen 443 ssl;"
			parts := strings.Fields(listenPort)
			if len(parts) > 0 {
				port := parts[0]
				// Remove IP if present (e.g., "0.0.0.0:80" -> "80")
				if strings.Contains(port, ":") {
					port = port[strings.LastIndex(port, ":")+1:]
				}
				site.Port = port

				// Check for SSL
				for _, part := range parts {
					if part == "ssl" {
						site.SSL = true
						break
					}
				}
			}
		}

		sites = append(sites, site)
	}

	// If no server blocks found, create a default entry
	if len(sites) == 0 {
		sites = append(sites, model.Site{
			Name:    "default",
			Enabled: true,
			Port:    "80",
			SSL:     false,
			Uptime:  containerUptime,
		})
	}

	return sites, nil
}

// parseServerBlocks extracts server blocks from nginx config
func parseServerBlocks(config string) []string {
	var blocks []string
	var currentBlock strings.Builder
	var braceCount int
	inServerBlock := false

	lines := strings.Split(config, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this line starts a server block
		if strings.HasPrefix(trimmed, "server") && strings.Contains(trimmed, "{") {
			inServerBlock = true
			braceCount = 1
			currentBlock.Reset()
			currentBlock.WriteString(line + "\n")
			continue
		}

		if inServerBlock {
			currentBlock.WriteString(line + "\n")

			// Count braces
			braceCount += strings.Count(line, "{")
			braceCount -= strings.Count(line, "}")

			// End of server block
			if braceCount == 0 {
				blocks = append(blocks, currentBlock.String())
				inServerBlock = false
			}
		}
	}

	return blocks
}

// extractDirective extracts a directive value from a server block
func extractDirective(block, directive string) string {
	lines := strings.Split(block, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, directive) {
			// Extract value after directive name
			parts := strings.SplitN(trimmed, directive, 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				// Remove trailing semicolon
				value = strings.TrimSuffix(value, ";")
				return strings.TrimSpace(value)
			}
		}
	}
	return ""
}

// getDockerContainerUptime gets the uptime of a Docker container
func getDockerContainerUptime(containerID string) string {
	// Get container status with uptime info
	cmd := exec.Command("docker", "inspect", "--format={{.State.Status}} {{.State.StartedAt}}", containerID)
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}

	// Parse output: "running 2024-11-05T12:00:00.000000000Z"
	parts := strings.Fields(string(output))
	if len(parts) < 2 || parts[0] != "running" {
		return "N/A"
	}
	
	// Use docker ps to get formatted uptime
	cmd = exec.Command("docker", "ps", "--filter", fmt.Sprintf("id=%s", containerID), "--format", "{{.Status}}")
	output, err = cmd.Output()
	if err != nil {
		return "N/A"
	}

	status := strings.TrimSpace(string(output))
	// Status format: "Up 2 hours" or "Up 3 days"
	if strings.HasPrefix(status, "Up ") {
		// Extract the uptime part: "Up 2 hours" -> "2h"
		uptime := strings.TrimPrefix(status, "Up ")
		uptime = strings.TrimSpace(uptime)
		
		// Convert to shorter format
		uptime = strings.Replace(uptime, " seconds", "s", 1)
		uptime = strings.Replace(uptime, " second", "s", 1)
		uptime = strings.Replace(uptime, " minutes", "m", 1)
		uptime = strings.Replace(uptime, " minute", "m", 1)
		uptime = strings.Replace(uptime, " hours", "h", 1)
		uptime = strings.Replace(uptime, " hour", "h", 1)
		uptime = strings.Replace(uptime, " days", "d", 1)
		uptime = strings.Replace(uptime, " day", "d", 1)
		uptime = strings.Replace(uptime, " weeks", "w", 1)
		uptime = strings.Replace(uptime, " week", "w", 1)
		
		return uptime
	}

	return "N/A"
}
