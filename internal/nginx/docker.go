package nginx

import (
	"fmt"
	"os/exec"
	"strings"
)

// DetectDockerNginx checks if NGINX is running in Docker and returns container ID
func DetectDockerNginx() (string, error) {
	// Check if there are nginx processes
	cmd := exec.Command("sh", "-c", "ps aux | grep 'nginx: master process' | grep -v grep | head -1")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		return "", fmt.Errorf("no nginx process found")
	}

	// Get the PID of nginx master process
	fields := strings.Fields(string(output))
	if len(fields) < 2 {
		return "", fmt.Errorf("could not parse nginx process")
	}
	pid := fields[1]

	// Check if this PID is in a Docker container
	cmd = exec.Command("sh", "-c", fmt.Sprintf("cat /proc/%s/cgroup 2>/dev/null | grep docker | head -1", pid))
	cgroupOutput, err := cmd.Output()
	if err != nil || len(cgroupOutput) == 0 {
		return "", fmt.Errorf("nginx not running in docker")
	}

	// Extract container ID from cgroup
	// Format can be:
	// - Old: 12:pids:/docker/CONTAINER_ID
	// - New: 0::/system.slice/docker-CONTAINER_ID.scope
	cgroupLine := string(cgroupOutput)

	// Try new format first: docker-CONTAINER_ID.scope
	if strings.Contains(cgroupLine, "docker-") {
		parts := strings.Split(cgroupLine, "docker-")
		if len(parts) >= 2 {
			containerID := strings.TrimSpace(parts[1])
			// Remove .scope suffix and take first 12 chars (short container ID)
			containerID = strings.Split(containerID, ".")[0]
			if len(containerID) >= 12 {
				containerID = containerID[:12]
			}
			return containerID, nil
		}
	}

	// Try old format: docker/CONTAINER_ID
	if strings.Contains(cgroupLine, "docker/") {
		parts := strings.Split(cgroupLine, "docker/")
		if len(parts) >= 2 {
			containerID := strings.TrimSpace(parts[1])
			// Container ID might have more path components, take first part
			containerID = strings.Split(containerID, "/")[0]
			if len(containerID) >= 12 {
				containerID = containerID[:12]
			}
			return containerID, nil
		}
	}

	return "", fmt.Errorf("could not extract container ID")
}

// GetDockerNginxPorts gets ports from Docker container
func GetDockerNginxPorts(containerID string) ([]string, error) {
	// Get port mappings from docker inspect
	cmd := exec.Command("docker", "inspect", "--format", "{{range $p, $conf := .NetworkSettings.Ports}}{{$p}} {{end}}", containerID)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	portsStr := strings.TrimSpace(string(output))
	if portsStr == "" {
		return nil, fmt.Errorf("no ports found")
	}

	// Parse ports (format: "80/tcp 443/tcp")
	var ports []string
	portEntries := strings.Fields(portsStr)
	for _, entry := range portEntries {
		// Remove /tcp or /udp suffix
		port := strings.Split(entry, "/")[0]
		ports = append(ports, port)
	}

	// Also get host port mappings
	cmd = exec.Command("docker", "port", containerID)
	portOutput, err := cmd.Output()
	if err == nil {
		// Parse output like: "80/tcp -> 0.0.0.0:8083"
		lines := strings.Split(string(portOutput), "\n")
		for _, line := range lines {
			if strings.Contains(line, "->") {
				parts := strings.Split(line, "->")
				if len(parts) >= 2 {
					// Extract host port from "0.0.0.0:8083"
					hostPart := strings.TrimSpace(parts[1])
					if strings.Contains(hostPart, ":") {
						hostPort := strings.Split(hostPart, ":")[1]
						ports = append(ports, hostPort)
					}
				}
			}
		}
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports detected")
	}

	return ports, nil
}

// IsDockerAvailable checks if docker command is available
func IsDockerAvailable() bool {
	cmd := exec.Command("docker", "--version")
	err := cmd.Run()
	return err == nil
}
