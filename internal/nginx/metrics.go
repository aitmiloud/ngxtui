package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Metrics represents real-time NGINX metrics
type Metrics struct {
	CPU            float64
	Memory         float64
	NetworkIn      float64 // Total bytes in (cumulative)
	NetworkOut     float64 // Total bytes out (cumulative)
	NetworkInRate  float64 // MB/s
	NetworkOutRate float64 // MB/s
	RequestRate    float64
	ActiveConns    int
	TotalConns     int64
	Timestamp      time.Time
}

// MetricsHistory stores historical metrics
type MetricsHistory struct {
	CPU        []float64
	Memory     []float64
	Network    []float64
	Requests   []float64
	Timestamps []time.Time
	MaxPoints  int
}

// NewMetricsHistory creates a new metrics history tracker
func NewMetricsHistory(maxPoints int) *MetricsHistory {
	return &MetricsHistory{
		CPU:        make([]float64, 0, maxPoints),
		Memory:     make([]float64, 0, maxPoints),
		Network:    make([]float64, 0, maxPoints),
		Requests:   make([]float64, 0, maxPoints),
		Timestamps: make([]time.Time, 0, maxPoints),
		MaxPoints:  maxPoints,
	}
}

// AddMetrics adds a metrics snapshot to history
func (h *MetricsHistory) AddMetrics(m *Metrics) {
	// Add new data
	h.CPU = append(h.CPU, m.CPU)
	h.Memory = append(h.Memory, m.Memory)
	h.Network = append(h.Network, m.NetworkIn+m.NetworkOut)
	h.Requests = append(h.Requests, m.RequestRate)
	h.Timestamps = append(h.Timestamps, m.Timestamp)

	// Keep only last MaxPoints
	if len(h.CPU) > h.MaxPoints {
		h.CPU = h.CPU[1:]
		h.Memory = h.Memory[1:]
		h.Network = h.Network[1:]
		h.Requests = h.Requests[1:]
		h.Timestamps = h.Timestamps[1:]
	}
}

// GetMetrics collects current NGINX metrics
func (s *Service) GetMetrics() (*Metrics, error) {
	metrics := &Metrics{
		Timestamp: time.Now(),
	}

	// Get CPU usage
	cpu, err := s.getNginxCPU()
	if err == nil {
		metrics.CPU = cpu
	}

	// Get memory usage
	mem, err := s.getNginxMemory()
	if err == nil {
		metrics.Memory = mem
	}

	// Get network stats
	netIn, netOut, err := s.getNetworkStats()
	if err == nil {
		metrics.NetworkIn = netIn
		metrics.NetworkOut = netOut
	}

	// Get connection stats
	activeConns, totalConns, err := s.getConnectionStats()
	if err == nil {
		metrics.ActiveConns = activeConns
		metrics.TotalConns = totalConns
	}

	// Get request rate
	rate, _, err := s.calculateRequestRate()
	if err == nil {
		metrics.RequestRate = rate
	}

	return metrics, nil
}

// getNginxCPU gets CPU usage percentage for nginx processes
func (s *Service) getNginxCPU() (float64, error) {
	cmd := exec.Command("sh", "-c", "ps aux | grep nginx | grep -v grep | awk '{sum+=$3} END {print sum}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	cpuStr := strings.TrimSpace(string(output))
	if cpuStr == "" {
		return 0, nil
	}

	cpu, err := strconv.ParseFloat(cpuStr, 64)
	if err != nil {
		return 0, err
	}

	return cpu, nil
}

// getNginxMemory gets memory usage percentage for nginx processes
func (s *Service) getNginxMemory() (float64, error) {
	cmd := exec.Command("sh", "-c", "ps aux | grep nginx | grep -v grep | awk '{sum+=$4} END {print sum}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	memStr := strings.TrimSpace(string(output))
	if memStr == "" {
		return 0, nil
	}

	mem, err := strconv.ParseFloat(memStr, 64)
	if err != nil {
		return 0, err
	}

	return mem, nil
}

// getNetworkStats gets network I/O statistics
func (s *Service) getNetworkStats() (float64, float64, error) {
	// Read network stats from /proc/net/dev
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(string(data), "\n")
	var totalRx, totalTx int64

	for _, line := range lines {
		// Skip header lines
		if !strings.Contains(line, ":") {
			continue
		}

		// Parse interface stats
		parts := strings.Fields(line)
		if len(parts) < 10 {
			continue
		}

		// Skip loopback
		if strings.HasPrefix(parts[0], "lo:") {
			continue
		}

		// Parse received and transmitted bytes
		rx, _ := strconv.ParseInt(parts[1], 10, 64)
		tx, _ := strconv.ParseInt(parts[9], 10, 64)

		totalRx += rx
		totalTx += tx
	}

	// Convert to MB
	rxMB := float64(totalRx) / (1024 * 1024)
	txMB := float64(totalTx) / (1024 * 1024)

	return rxMB, txMB, nil
}

// getConnectionStats gets connection statistics
func (s *Service) getConnectionStats() (int, int64, error) {
	// Get all ports NGINX is listening on
	ports, err := s.GetListeningPorts()
	if err != nil || len(ports) == 0 {
		// Fallback to port 80 if detection fails
		ports = []string{"80"}
	}

	// Build grep pattern for all ports
	var portPattern string
	for i, port := range ports {
		if i > 0 {
			portPattern += "\\|"
		}
		portPattern += ":" + port
	}

	// Get active connections on all detected ports
	activeCmd := exec.Command("sh", "-c", fmt.Sprintf("ss -tn | grep -E '(%s)' | grep ESTAB | wc -l", portPattern))
	activeOutput, err := activeCmd.Output()
	if err != nil {
		return 0, 0, err
	}

	activeConns, _ := strconv.Atoi(strings.TrimSpace(string(activeOutput)))

	// Get total connections from logs (approximate)
	totalCmd := exec.Command("sh", "-c", "wc -l < /var/log/nginx/access.log")
	totalOutput, err := totalCmd.Output()
	if err != nil {
		return activeConns, 0, nil
	}

	totalConns, _ := strconv.ParseInt(strings.TrimSpace(string(totalOutput)), 10, 64)

	return activeConns, totalConns, nil
}

// GetSystemMetrics gets overall system metrics
func (s *Service) GetSystemMetrics() (*SystemMetrics, error) {
	sysMetrics := &SystemMetrics{}

	// Get load average
	loadData, err := os.ReadFile("/proc/loadavg")
	if err == nil {
		parts := strings.Fields(string(loadData))
		if len(parts) >= 3 {
			sysMetrics.LoadAvg1, _ = strconv.ParseFloat(parts[0], 64)
			sysMetrics.LoadAvg5, _ = strconv.ParseFloat(parts[1], 64)
			sysMetrics.LoadAvg15, _ = strconv.ParseFloat(parts[2], 64)
		}
	}

	// Get memory info
	memData, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		lines := strings.Split(string(memData), "\n")
		memInfo := make(map[string]int64)

		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				key := strings.TrimSuffix(parts[0], ":")
				value, _ := strconv.ParseInt(parts[1], 10, 64)
				memInfo[key] = value
			}
		}

		total := memInfo["MemTotal"]
		available := memInfo["MemAvailable"]
		if total > 0 {
			used := total - available
			sysMetrics.MemoryUsedPercent = float64(used) / float64(total) * 100
			sysMetrics.MemoryTotal = total * 1024 // Convert to bytes
			sysMetrics.MemoryUsed = used * 1024
		}
	}

	// Get disk usage
	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) >= 2 {
			parts := strings.Fields(lines[1])
			if len(parts) >= 5 {
				sysMetrics.DiskUsage = strings.TrimSuffix(parts[4], "%")
			}
		}
	}

	return sysMetrics, nil
}

// SystemMetrics represents system-wide metrics
type SystemMetrics struct {
	LoadAvg1          float64
	LoadAvg5          float64
	LoadAvg15         float64
	MemoryTotal       int64
	MemoryUsed        int64
	MemoryUsedPercent float64
	DiskUsage         string
}
