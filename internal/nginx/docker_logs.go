package nginx

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// GetDockerAccessLogs reads access logs from Docker container
func GetDockerAccessLogs(containerID string, lines int) ([]LogEntry, error) {
	// Try to get logs from container with timeout
	cmd := exec.Command("docker", "logs", "--tail", fmt.Sprintf("%d", lines), containerID)
	output, err := cmd.Output()
	if err != nil {
		// Fallback: try reading from log file inside container
		cmd = exec.Command("docker", "exec", containerID, "sh", "-c", fmt.Sprintf("tail -n %d /var/log/nginx/access.log 2>/dev/null || echo ''", lines))
		output, err = cmd.Output()
		if err != nil {
			// Return empty if can't read logs
			return []LogEntry{}, nil
		}
	}

	// Parse log lines
	var entries []LogEntry
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseAccessLogLine(line)
		if err == nil {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// CalculateDockerRequestRate calculates request rate from Docker container logs
func CalculateDockerRequestRate(containerID string) (float64, int64, error) {
	// Try docker logs first (faster)
	cmd := exec.Command("docker", "logs", "--tail", "1000", containerID)
	output, err := cmd.Output()
	if err != nil {
		// Fallback: try reading from log file
		cmd = exec.Command("docker", "exec", containerID, "sh", "-c", "tail -n 1000 /var/log/nginx/access.log 2>/dev/null || echo ''")
		output, err = cmd.Output()
		if err != nil || len(output) == 0 {
			// Return 0 if can't read logs
			return 0, 0, nil
		}
	}

	var totalRequests int64
	var recentRequests int
	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
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

	return rate, totalRequests, nil
}
