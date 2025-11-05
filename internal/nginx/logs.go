package nginx

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

// LogEntry represents a parsed NGINX log entry
type LogEntry struct {
	IP          string
	Timestamp   time.Time
	Method      string
	Path        string
	StatusCode  int
	BytesSent   int
	UserAgent   string
	Referer     string
	StatusClass string // "2xx", "3xx", "4xx", "5xx"
}

// GetAccessLogs returns the last N lines of the access log
func (s *Service) GetAccessLogs(maxLines int) ([]LogEntry, error) {
	// Check if Docker NGINX (with caching)
	if IsDockerAvailable() {
		containerID, err := getCachedContainerID()
		if err == nil {
			return GetDockerAccessLogs(containerID, maxLines)
		}
	}

	// Native NGINX
	logPath := "/var/log/nginx/access.log"

	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open access log: %w", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	// Read all lines first, then take last N
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Take last maxLines
	startIdx := 0
	if len(lines) > maxLines {
		startIdx = len(lines) - maxLines
	}

	// Parse each line
	for i := startIdx; i < len(lines); i++ {
		entry, err := parseAccessLogLine(lines[i])
		if err == nil {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// parseAccessLogLine parses a standard NGINX access log line
// Format: $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"
func parseAccessLogLine(line string) (LogEntry, error) {
	// Regex pattern for standard combined log format
	pattern := `^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) ([^"]+) \S+" (\d+) (\d+) "([^"]*)" "([^"]*)"`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(line)
	if len(matches) < 9 {
		return LogEntry{}, fmt.Errorf("failed to parse log line")
	}

	// Parse timestamp
	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		timestamp = time.Now()
	}

	// Parse status code
	var statusCode int
	fmt.Sscanf(matches[5], "%d", &statusCode)

	// Parse bytes sent
	var bytesSent int
	fmt.Sscanf(matches[6], "%d", &bytesSent)

	// Determine status class
	statusClass := "5xx"
	if statusCode >= 200 && statusCode < 300 {
		statusClass = "2xx"
	} else if statusCode >= 300 && statusCode < 400 {
		statusClass = "3xx"
	} else if statusCode >= 400 && statusCode < 500 {
		statusClass = "4xx"
	}

	return LogEntry{
		IP:          matches[1],
		Timestamp:   timestamp,
		Method:      matches[3],
		Path:        matches[4],
		StatusCode:  statusCode,
		BytesSent:   bytesSent,
		Referer:     matches[7],
		UserAgent:   matches[8],
		StatusClass: statusClass,
	}, nil
}

// GetErrorLogs returns recent error log entries
func (s *Service) GetErrorLogs(maxLines int) ([]string, error) {
	logPath := "/var/log/nginx/error.log"

	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open error log: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Take last maxLines
	startIdx := 0
	if len(lines) > maxLines {
		startIdx = len(lines) - maxLines
	}

	return lines[startIdx:], nil
}

// GetLogStats returns statistics from access logs
func (s *Service) GetLogStats() (*LogStats, error) {
	entries, err := s.GetAccessLogs(1000) // Analyze last 1000 entries
	if err != nil {
		return nil, err
	}

	stats := &LogStats{
		TotalRequests: len(entries),
		StatusCounts:  make(map[string]int),
		MethodCounts:  make(map[string]int),
		TopPaths:      make(map[string]int),
	}

	var totalBytes int64
	uniqueIPs := make(map[string]bool)

	for _, entry := range entries {
		// Count by status class
		stats.StatusCounts[entry.StatusClass]++

		// Count by method
		stats.MethodCounts[entry.Method]++

		// Track top paths
		stats.TopPaths[entry.Path]++

		// Track unique IPs
		uniqueIPs[entry.IP] = true

		// Sum bytes
		totalBytes += int64(entry.BytesSent)
	}

	stats.UniqueIPs = len(uniqueIPs)
	stats.TotalBytes = totalBytes

	// Calculate average bytes per request
	if stats.TotalRequests > 0 {
		stats.AvgBytesPerRequest = totalBytes / int64(stats.TotalRequests)
	}

	return stats, nil
}

// LogStats represents aggregated log statistics
type LogStats struct {
	TotalRequests      int
	UniqueIPs          int
	StatusCounts       map[string]int // "2xx", "3xx", "4xx", "5xx"
	MethodCounts       map[string]int // "GET", "POST", etc.
	TopPaths           map[string]int
	TotalBytes         int64
	AvgBytesPerRequest int64
}

// FormatLogEntry formats a log entry for display with colors
func FormatLogEntry(entry LogEntry) string {
	// Color codes based on status
	var statusColor string
	switch entry.StatusClass {
	case "2xx":
		statusColor = "\033[32m" // Green
	case "3xx":
		statusColor = "\033[36m" // Cyan
	case "4xx":
		statusColor = "\033[33m" // Yellow
	case "5xx":
		statusColor = "\033[31m" // Red
	default:
		statusColor = "\033[37m" // White
	}

	// Format timestamp
	timeStr := entry.Timestamp.Format("15:04:05")

	// Format bytes
	bytesStr := formatBytes(entry.BytesSent)

	// Build formatted line
	line := fmt.Sprintf("%s●\033[0m \033[90m%s\033[0m \033[37m%-15s\033[0m \033[36m%-6s\033[0m \033[97m%-40s\033[0m %s%3d\033[0m \033[90m%s\033[0m",
		statusColor,
		timeStr,
		entry.IP,
		entry.Method,
		truncateString(entry.Path, 40),
		statusColor,
		entry.StatusCode,
		bytesStr,
	)

	return line
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
	}
}

// truncateString truncates a string to maxLen with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
