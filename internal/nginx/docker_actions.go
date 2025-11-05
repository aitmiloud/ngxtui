package nginx

import (
	"fmt"
	"os/exec"
)

// DockerTestConfig tests NGINX config in Docker container
func DockerTestConfig(containerID string) error {
	cmd := exec.Command("docker", "exec", containerID, "nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config test failed: %s", string(output))
	}
	return nil
}

// DockerReload reloads NGINX in Docker container
func DockerReload(containerID string) error {
	cmd := exec.Command("docker", "exec", containerID, "nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reload failed: %s", string(output))
	}
	return nil
}

// DockerEnableSite - Not applicable for Docker (sites are in config)
func DockerEnableSite(containerID, siteName string) error {
	return fmt.Errorf("enable/disable not supported for Docker NGINX - edit config and restart container")
}

// DockerDisableSite - Not applicable for Docker (sites are in config)
func DockerDisableSite(containerID, siteName string) error {
	return fmt.Errorf("enable/disable not supported for Docker NGINX - edit config and restart container")
}
