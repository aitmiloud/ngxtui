package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateSiteConfig creates a new NGINX site configuration file
func (s *Service) CreateSiteConfig(filename, content string) error {
	// Check if Docker mode
	if IsDockerAvailable() {
		if containerID, err := DetectDockerNginx(); err == nil {
			return s.createDockerSiteConfig(containerID, filename, content)
		}
	}

	// Native mode
	return s.createNativeSiteConfig(filename, content)
}

// createNativeSiteConfig creates a site config in native NGINX
func (s *Service) createNativeSiteConfig(filename, content string) error {
	// Determine which directory structure to use
	var configPath string
	var symlinkPath string
	useSitesAvailable := false

	// Check if sites-available exists (Debian/Ubuntu style)
	if _, err := os.Stat("/etc/nginx/sites-available"); err == nil {
		useSitesAvailable = true
		sitesAvailable := "/etc/nginx/sites-available"
		sitesEnabled := "/etc/nginx/sites-enabled"
		
		// Ensure directories exist
		if err := os.MkdirAll(sitesAvailable, 0755); err != nil {
			return fmt.Errorf("failed to create sites-available directory: %w", err)
		}
		if err := os.MkdirAll(sitesEnabled, 0755); err != nil {
			return fmt.Errorf("failed to create sites-enabled directory: %w", err)
		}
		
		configPath = filepath.Join(sitesAvailable, filename)
		symlinkPath = filepath.Join(sitesEnabled, filename)
	} else {
		// Use conf.d (RHEL/CentOS/Fedora style)
		confD := "/etc/nginx/conf.d"
		
		// Ensure directory exists
		if err := os.MkdirAll(confD, 0755); err != nil {
			return fmt.Errorf("failed to create conf.d directory: %w", err)
		}
		
		configPath = filepath.Join(confD, filename)
	}

	// Write configuration file
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Create symlink if using sites-available/sites-enabled
	if useSitesAvailable {
		if err := os.Symlink(configPath, symlinkPath); err != nil {
			// If symlink already exists, that's okay
			if !os.IsExist(err) {
				// Rollback: remove config file
				os.Remove(configPath)
				return fmt.Errorf("failed to create symlink: %w", err)
			}
		}
	}

	// Test configuration
	if err := s.TestConfig(); err != nil {
		// Rollback: remove the files
		if useSitesAvailable {
			os.Remove(symlinkPath)
		}
		os.Remove(configPath)
		return fmt.Errorf("configuration test failed: %w", err)
	}

	// Reload NGINX
	if err := s.Reload(); err != nil {
		return fmt.Errorf("failed to reload NGINX: %w", err)
	}

	return nil
}

// createDockerSiteConfig creates a site config in Docker NGINX
func (s *Service) createDockerSiteConfig(containerID, filename, content string) error {
	// For Docker, we need to:
	// 1. Create a temp file with the config
	// 2. Copy it into the container
	// 3. Test the config
	// 4. Reload NGINX

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "nginx-config-*.conf")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to temp file
	if _, err := tmpFile.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close()

	// Determine the target path in container
	// Check which directory structure the container uses
	checkCmd := exec.Command("docker", "exec", containerID, "test", "-d", "/etc/nginx/sites-enabled")
	useSitesEnabled := checkCmd.Run() == nil
	
	var targetPath string
	if useSitesEnabled {
		// Debian/Ubuntu style - need to create in sites-available and link to sites-enabled
		targetPath = fmt.Sprintf("/etc/nginx/sites-available/%s", filename)
	} else {
		// RHEL/CentOS style
		targetPath = fmt.Sprintf("/etc/nginx/conf.d/%s", filename)
	}

	// Copy file into container
	cmd := exec.Command("docker", "cp", tmpFile.Name(), fmt.Sprintf("%s:%s", containerID, targetPath))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to copy config to container: %s", string(output))
	}

	// Fix file permissions in container (make it readable by nginx)
	cmd = exec.Command("docker", "exec", containerID, "chmod", "644", targetPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set file permissions: %s", string(output))
	}

	// Create symlink if using sites-enabled
	if useSitesEnabled {
		symlinkPath := fmt.Sprintf("/etc/nginx/sites-enabled/%s", filename)
		cmd = exec.Command("docker", "exec", containerID, "ln", "-sf", targetPath, symlinkPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to create symlink: %s", string(output))
		}
	}

	// Test configuration
	cmd = exec.Command("docker", "exec", containerID, "nginx", "-t")
	if output, err := cmd.CombinedOutput(); err != nil {
		// Rollback: remove the file from container
		exec.Command("docker", "exec", containerID, "rm", targetPath).Run()
		return fmt.Errorf("configuration test failed: %s", string(output))
	}

	// Reload NGINX in container
	cmd = exec.Command("docker", "exec", containerID, "nginx", "-s", "reload")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload NGINX: %s", string(output))
	}

	return nil
}

// GetSiteConfigPath returns the path where site configs are stored
func (s *Service) GetSiteConfigPath() (string, error) {
	if IsDockerAvailable() {
		if _, err := DetectDockerNginx(); err == nil {
			return "/etc/nginx/conf.d", nil
		}
	}
	return "/etc/nginx/sites-available", nil
}
