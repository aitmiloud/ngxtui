package nginx

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	crossplane "github.com/nginxinc/nginx-go-crossplane"
)

// Stats represents NGINX statistics
type Stats struct {
	ActiveConnections int
	RequestRate       float64
	TotalRequests     int64
	Uptime            time.Duration
	WorkerProcesses   int
}

// GetStats retrieves real NGINX statistics
func (s *Service) GetStats() (*Stats, error) {
	stats := &Stats{}

	// Get active connections from nginx status
	connections, err := s.getActiveConnections()
	if err == nil {
		stats.ActiveConnections = connections
	}

	// Get worker processes
	workers, err := s.getWorkerProcesses()
	if err == nil {
		stats.WorkerProcesses = workers
	}

	// Get uptime
	uptime, err := s.getNginxUptime()
	if err == nil {
		stats.Uptime = uptime
	}

	// Calculate request rate from access logs
	rate, total, err := s.calculateRequestRate()
	if err == nil {
		stats.RequestRate = rate
		stats.TotalRequests = total
	}

	return stats, nil
}

// getActiveConnections gets the number of active connections
func (s *Service) getActiveConnections() (int, error) {
	// Get all ports NGINX is listening on
	ports, err := s.GetListeningPorts()
	if err != nil || len(ports) == 0 {
		// Fallback to port 80 if detection fails
		ports = []string{"80"}
	}

	// Build grep pattern for all ports: ":80\|:443\|:8080"
	var portPattern string
	for i, port := range ports {
		if i > 0 {
			portPattern += "\\|"
		}
		portPattern += ":" + port
	}

	// Count connections on all detected ports
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ss -tn | grep -E '(%s)' | wc -l", portPattern))
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// getWorkerProcesses counts nginx worker processes
func (s *Service) getWorkerProcesses() (int, error) {
	cmd := exec.Command("sh", "-c", "ps aux | grep 'nginx: worker process' | grep -v grep | wc -l")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// getNginxUptime gets nginx process uptime
func (s *Service) getNginxUptime() (time.Duration, error) {
	// Get nginx master process start time
	cmd := exec.Command("sh", "-c", "ps -eo pid,etime,cmd | grep 'nginx: master process' | grep -v grep | awk '{print $2}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	etimeStr := strings.TrimSpace(string(output))
	if etimeStr == "" {
		return 0, fmt.Errorf("nginx not running")
	}

	// Parse etime format (can be: ss, mm:ss, hh:mm:ss, or dd-hh:mm:ss)
	duration, err := parseElapsedTime(etimeStr)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

// parseElapsedTime parses ps etime format
func parseElapsedTime(etime string) (time.Duration, error) {
	var days, hours, minutes, seconds int

	// Check for days
	if strings.Contains(etime, "-") {
		parts := strings.Split(etime, "-")
		days, _ = strconv.Atoi(parts[0])
		etime = parts[1]
	}

	// Parse time components
	timeParts := strings.Split(etime, ":")
	switch len(timeParts) {
	case 1: // ss
		seconds, _ = strconv.Atoi(timeParts[0])
	case 2: // mm:ss
		minutes, _ = strconv.Atoi(timeParts[0])
		seconds, _ = strconv.Atoi(timeParts[1])
	case 3: // hh:mm:ss
		hours, _ = strconv.Atoi(timeParts[0])
		minutes, _ = strconv.Atoi(timeParts[1])
		seconds, _ = strconv.Atoi(timeParts[2])
	}

	totalSeconds := days*86400 + hours*3600 + minutes*60 + seconds
	return time.Duration(totalSeconds) * time.Second, nil
}

// calculateRequestRate calculates requests per second from access log
func (s *Service) calculateRequestRate() (float64, int64, error) {
	// Check if Docker NGINX (with caching)
	if IsDockerAvailable() {
		containerID, err := getCachedContainerID()
		if err == nil {
			return CalculateDockerRequestRate(containerID)
		}
	}

	// Native NGINX
	logPath := "/var/log/nginx/access.log"

	file, err := os.Open(logPath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var totalRequests int64
	var recentRequests int
	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalRequests++
		line := scanner.Text()

		// Try to parse timestamp from log line
		// Typical format: 127.0.0.1 - - [05/Nov/2024:10:24:57 +0100] ...
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			start := strings.Index(line, "[") + 1
			end := strings.Index(line, "]")
			if start > 0 && end > start {
				timestampStr := line[start:end]
				// Parse nginx timestamp format
				timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", timestampStr)
				if err == nil && timestamp.After(oneMinuteAgo) {
					recentRequests++
				}
			}
		}
	}

	// Calculate rate (requests per second over last minute)
	rate := float64(recentRequests) / 60.0

	return rate, totalRequests, scanner.Err()
}

// GetConfigErrors validates nginx configuration and returns errors
func (s *Service) GetConfigErrors() ([]string, error) {
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)
	var errors []string

	if err != nil {
		// Parse error messages from output
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			if strings.Contains(line, "error") || strings.Contains(line, "failed") {
				errors = append(errors, strings.TrimSpace(line))
			}
		}
		return errors, fmt.Errorf("configuration has errors")
	}

	return nil, nil
}

// GetServerNames extracts all server names from configurations
func (s *Service) GetServerNames() (map[string][]string, error) {
	serverNames := make(map[string][]string)

	entries, err := os.ReadDir(sitesAvailableDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		siteName := entry.Name()
		sitePath := filepath.Join(sitesAvailableDir, siteName)

		payload, err := crossplane.Parse(sitePath, &crossplane.ParseOptions{
			SingleFile: true,
		})
		if err != nil {
			continue
		}

		var names []string
		if len(payload.Config) > 0 {
			for _, directive := range payload.Config[0].Parsed {
				if directive.Directive == "server" {
					names = append(names, s.extractServerNames(directive.Block)...)
				}
			}
		}

		if len(names) > 0 {
			serverNames[siteName] = names
		}
	}

	return serverNames, nil
}

// extractServerNames extracts server_name directives from a server block
func (s *Service) extractServerNames(block crossplane.Directives) []string {
	var names []string
	for _, directive := range block {
		if directive.Directive == "server_name" {
			names = append(names, directive.Args...)
		}
	}
	return names
}
