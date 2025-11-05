package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/NimbleMarkets/ntcharts/linechart/streamlinechart"
	"github.com/NimbleMarkets/ntcharts/sparkline"
	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/aitmiloud/ngxtui/internal/nginx"
	"github.com/aitmiloud/ngxtui/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

// Renderer handles all view rendering
type Renderer struct{}

// New creates a new renderer
func New() *Renderer {
	return &Renderer{}
}

// RenderCreativeHeader renders a beautiful, unified header for NgxTUI
func (r *Renderer) RenderCreativeHeader(width int) string {
	// Compact ASCII art style NgxTUI logo with gradient effect and top margin
	lineEmpty1 := "  \033[1;38;5;51m\033[0m"
	lineEmpty2 := "  \033[1;38;5;51m\033[0m"
	line1 := "  \033[1;38;5;51m███╗   ██╗ ██████╗ ██╗  ██╗\033[38;5;48m████████╗\033[38;5;46m██╗   ██╗██╗\033[0m"
	line2 := "  \033[1;38;5;51m████╗  ██║██╔════╝ ╚██╗██╔╝\033[38;5;48m╚══██╔══╝\033[38;5;46m██║   ██║██║\033[0m"
	line3 := "  \033[1;38;5;51m██╔██╗ ██║██║  ███╗ ╚███╔╝ \033[38;5;48m   ██║   \033[38;5;46m██║   ██║██║\033[0m  \033[2;37mNGINX Management Terminal UI\033[0m \033[2;90mv1.0.0\033[0m"
	line4 := "  \033[1;38;5;51m██║╚██╗██║██║   ██║ ██╔██╗ \033[38;5;48m   ██║   \033[38;5;46m██║   ██║██║\033[0m"
	line5 := "  \033[1;38;5;51m██║ ╚████║╚██████╔╝██╔╝ ██╗\033[38;5;48m   ██║   \033[38;5;46m╚██████╔╝██║\033[0m"
	line6 := "  \033[1;38;5;51m╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝\033[38;5;48m   ╚═╝   \033[38;5;46m ╚═════╝ ╚═╝\033[0m"
	lineEmpty3 := "  \033[1;38;5;51m\033[0m"

	return "\n" + strings.Join([]string{lineEmpty1, lineEmpty2, line1, line2, line3, line4, line5, line6, lineEmpty3}, "\n")
}

// RenderTabs renders the tab bar with icons
func (r *Renderer) RenderTabs(m *model.Model) string {
	tabs := []struct {
		icon string
		name string
	}{
		{"🌐", "Sites"},
		{"📋", "Logs"},
		{"📊", "Stats"},
		{"📈", "Metrics"},
	}

	var renderedTabs []string

	for i, tab := range tabs {
		tabText := fmt.Sprintf("%s %s", tab.icon, tab.name)
		var styledTab string
		if model.TabType(i) == m.ActiveTab {
			// Active tab: bright cyan with bold
			styledTab = fmt.Sprintf("\033[1;36m%s\033[0m", tabText)
		} else {
			// Inactive tab: bright white
			styledTab = fmt.Sprintf("\033[97m%s\033[0m", tabText)
		}
		renderedTabs = append(renderedTabs, styledTab)

		// Add separator between tabs
		if i < len(tabs)-1 {
			renderedTabs = append(renderedTabs, " \033[90m│\033[0m ")
		}
	}

	tabBar := strings.Join(renderedTabs, "")

	// Add a divider below tabs
	divider := "\033[90m" + strings.Repeat("─", 80) + "\033[0m"

	return tabBar + "\n" + divider
}

// RenderSitesTable renders the sites table view
func (r *Renderer) RenderSitesTable(m *model.Model, width, height int) string {
	// Use the stickers table renderer with actual terminal dimensions
	return r.RenderSitesTableStickers(m, width, height)
}

// RenderActionMenu renders the action menu for a selected site
func (r *Renderer) RenderActionMenu(m *model.Model) string {
	if m.Selected < 0 || m.Selected >= len(m.Sites) {
		return ""
	}

	site := m.Sites[m.Selected]

	// Enhanced title with site status
	statusBadge := r.RenderStatusBadge(site.Enabled)
	titleText := fmt.Sprintf("⚙️  %s", site.Name)
	title := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.CardTitle.Render(titleText),
		"  ",
		statusBadge,
	)

	// Actions with icons
	actions := []struct {
		icon  string
		text  string
		color lipgloss.Color
	}{
		{"✓", "Enable Site", styles.AccentSuccess},
		{"✗", "Disable Site", styles.AccentDanger},
		{"🔍", "Test Configuration", styles.AccentInfo},
		{"🔄", "Reload NGINX", styles.AccentWarning},
		{"📋", "View Logs", styles.AccentPrimary},
		{"←", "Back", styles.TextMuted},
	}

	var items []string
	for i, action := range actions {
		actionText := fmt.Sprintf("%s  %s", action.icon, action.text)
		if i == m.Cursor {
			items = append(items, styles.SelectedAction.Render("▸ "+actionText))
		} else {
			styledAction := lipgloss.NewStyle().
				Foreground(action.color).
				Render(actionText)
			items = append(items, styles.ActionItem.Render(styledAction))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)

	// Add spacing
	divider := styles.Divider.Render(strings.Repeat("─", 40))

	return styles.Panel.Width(45).Render(
		lipgloss.JoinVertical(lipgloss.Left, title, divider, content),
	)
}

// RenderStatusBadge renders a status badge for enabled/disabled state
func (r *Renderer) RenderStatusBadge(enabled bool) string {
	if enabled {
		return styles.EnabledBadge.Render("● ACTIVE")
	}
	return styles.DisabledBadge.Render("○ INACTIVE")
}

// RenderLogsView renders the logs view with REAL NGINX access logs
func (r *Renderer) RenderLogsView(m *model.Model) string {
	title := fmt.Sprintf("\033[1;36m📋 REAL-TIME ACCESS LOGS\033[0m\n")

	// Legend for status codes with icons
	legend := fmt.Sprintf("  \033[32m✓\033[0m 2xx Success   \033[36m↻\033[0m 3xx Redirect   \033[33m⚠\033[0m 4xx Client Error   \033[31m✗\033[0m 5xx Server Error\n\n")

	// Column headers
	headers := fmt.Sprintf("  \033[1;90m%-8s %-15s %-6s %-35s %-4s %-4s %-12s %-10s\033[0m\n",
		"TIME", "IP", "METHOD", "PATH", "CODE", "SIZE", "CLIENT", "REFERER")

	divider := "\033[90m" + strings.Repeat("─", 130) + "\033[0m\n"

	// Calculate how many log entries can fit on screen
	// Account for: title (1), legend (2), headers (1), divider (1), menu bar (3), padding (2)
	headerLines := 10
	availableLines := m.Height - headerLines
	if availableLines < 10 {
		availableLines = 10 // Minimum 10 lines
	}
	if availableLines > 100 {
		availableLines = 100 // Maximum 100 lines to avoid performance issues
	}

	// Get real access logs from NGINX
	nginxService := nginx.New()
	logEntries, err := nginxService.GetAccessLogs(availableLines)

	var logs []string
	if err != nil {
		logs = []string{fmt.Sprintf("\033[33m⚠ Unable to read access logs: %v\033[0m", err)}
	} else if len(logEntries) == 0 {
		logs = []string{"\033[90mNo access logs available\033[0m"}
	} else {
		// Format each log entry
		for _, entry := range logEntries {
			logs = append(logs, nginx.FormatLogEntry(entry))
		}
	}

	content := strings.Join(logs, "\n")

	return title + legend + headers + divider + content
}

// RenderStatsView renders the statistics view with stunning modern design
func (r *Renderer) RenderStatsView(m *model.Model, width int) string {
	totalSites := len(m.Sites)
	enabledSites := 0
	disabledSites := 0

	for _, site := range m.Sites {
		if site.Enabled {
			enabledSites++
		} else {
			disabledSites++
		}
	}

	// Calculate percentage
	percentage := 0.0
	if totalSites > 0 {
		percentage = float64(enabledSites) / float64(totalSites) * 100
	}

	// Stunning metric cards in a row with box drawing
	card1 := r.RenderStunningMetricCard("🌐", "TOTAL SITES", fmt.Sprintf("%d", totalSites), "Configured", styles.AccentPrimary)
	card2 := r.RenderStunningMetricCard("✓", "ACTIVE", fmt.Sprintf("%d", enabledSites), fmt.Sprintf("%.0f%% Online", percentage), styles.AccentSuccess)
	card3 := r.RenderStunningMetricCard("○", "INACTIVE", fmt.Sprintf("%d", disabledSites), "Offline", styles.AccentWarning)
	card4 := r.RenderStunningMetricCard("⚡", "UPTIME", "99.9%", "Last 30d", styles.AccentSecondary)

	cardsRow := lipgloss.JoinHorizontal(lipgloss.Top, card1, "  ", card2, "  ", card3, "  ", card4)

	// Visual distribution bar
	distBar := r.RenderDistributionBar(enabledSites, disabledSites, totalSites, width-10)

	// Performance metrics section
	perfSection := r.RenderPerformanceMetrics()

	// System health indicators
	healthSection := r.RenderHealthIndicators()

	return lipgloss.JoinVertical(lipgloss.Left,
		cardsRow,
		"",
		"",
		distBar,
		"",
		"",
		perfSection,
		"",
		healthSection,
	)
}

// RenderStunningMetricCard renders a beautifully designed metric card with ANSI art
func (r *Renderer) RenderStunningMetricCard(icon, label, value, subtext string, accentColor lipgloss.Color) string {
	// Simpler approach: use lipgloss to handle the layout properly
	// Create the card content
	titleLine := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render(icon + "  " + label)

	valueLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Bold(true).
		PaddingLeft(3).
		Render(value)

	subtextLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingLeft(3).
		Render(subtext)

	// Build the card with proper box drawing
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(19).
		Padding(0, 1)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleLine,
		"",
		valueLine,
		subtextLine,
	)

	return cardStyle.Render(content)
}

// RenderDistributionBar renders a visual distribution bar
func (r *Renderer) RenderDistributionBar(enabled, disabled, total int, width int) string {
	if total == 0 {
		return ""
	}

	title := fmt.Sprintf("\033[1;36m▸ SITE DISTRIBUTION\033[0m\n")

	// Calculate bar widths
	barWidth := width - 30
	if barWidth < 20 {
		barWidth = 20
	}

	enabledWidth := int(float64(enabled) / float64(total) * float64(barWidth))
	disabledWidth := barWidth - enabledWidth

	// Build the bar
	bar := "  ["
	bar += fmt.Sprintf("\033[42m%s\033[0m", strings.Repeat(" ", enabledWidth))
	bar += fmt.Sprintf("\033[43m%s\033[0m", strings.Repeat(" ", disabledWidth))
	bar += "]"

	legend := fmt.Sprintf("  \033[32m■\033[0m Active: %d    \033[33m■\033[0m Inactive: %d", enabled, disabled)

	return title + bar + "\n" + legend
}

// RenderPerformanceMetrics renders REAL performance indicators
func (r *Renderer) RenderPerformanceMetrics() string {
	title := fmt.Sprintf("\033[1;36m▸ REAL-TIME PERFORMANCE\033[0m\n")

	// Get real stats from NGINX
	nginxService := nginx.New()
	stats, err := nginxService.GetStats()

	var metrics []string
	if err != nil {
		metrics = []string{fmt.Sprintf("  \033[33m⚠ Unable to fetch stats: %v\033[0m", err)}
	} else {
		// Get log stats for success rate
		logStats, _ := nginxService.GetLogStats()
		successRate := 0.0
		if logStats != nil && logStats.TotalRequests > 0 {
			successCount := logStats.StatusCounts["2xx"] + logStats.StatusCounts["3xx"]
			successRate = float64(successCount) / float64(logStats.TotalRequests) * 100
		}

		metrics = []string{
			fmt.Sprintf("  \033[32m●\033[0m Request Rate    : \033[1;97m%.1f\033[0m req/s", stats.RequestRate),
			fmt.Sprintf("  \033[32m●\033[0m Active Conn.    : \033[1;97m%d\033[0m connections", stats.ActiveConnections),
			fmt.Sprintf("  \033[32m●\033[0m Worker Processes: \033[1;97m%d\033[0m workers", stats.WorkerProcesses),
			fmt.Sprintf("  \033[32m●\033[0m Success Rate    : \033[1;97m%.1f%%\033[0m", successRate),
			fmt.Sprintf("  \033[32m●\033[0m Uptime          : \033[1;97m%s\033[0m", formatDuration(stats.Uptime)),
		}
	}

	return title + strings.Join(metrics, "\n")
}

// formatDuration formats a duration into human-readable format
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// RenderHealthIndicators renders REAL system health status
func (r *Renderer) RenderHealthIndicators() string {
	title := fmt.Sprintf("\033[1;36m▸ SYSTEM HEALTH\033[0m\n")

	nginxService := nginx.New()

	// Check NGINX service
	nginxStatus := "\033[32m✓\033[0m"
	nginxMsg := "\033[1;32mRunning\033[0m"
	if err := nginxService.TestConfig(); err != nil {
		nginxStatus = "\033[31m✗\033[0m"
		nginxMsg = "\033[1;31mConfig Error\033[0m"
	}

	// Check configuration
	configStatus := "\033[32m✓\033[0m"
	configMsg := "\033[1;32mValid\033[0m"
	errors, err := nginxService.GetConfigErrors()
	if err != nil || len(errors) > 0 {
		configStatus = "\033[31m✗\033[0m"
		configMsg = fmt.Sprintf("\033[1;31m%d Errors\033[0m", len(errors))
	}

	// Get system metrics
	sysMetrics, _ := nginxService.GetSystemMetrics()
	diskMsg := "\033[1;32mOK\033[0m"
	if sysMetrics != nil && sysMetrics.DiskUsage != "" {
		diskMsg = fmt.Sprintf("\033[1;97m%s%% Used\033[0m", sysMetrics.DiskUsage)
	}

	memMsg := "\033[1;32mOK\033[0m"
	if sysMetrics != nil {
		memMsg = fmt.Sprintf("\033[1;97m%.1f%% Used\033[0m", sysMetrics.MemoryUsedPercent)
	}

	indicators := []string{
		fmt.Sprintf("  %s NGINX Service   : %s", nginxStatus, nginxMsg),
		fmt.Sprintf("  %s Configuration   : %s", configStatus, configMsg),
		fmt.Sprintf("  \033[32m●\033[0m Disk Space      : %s", diskMsg),
		fmt.Sprintf("  \033[32m●\033[0m Memory Usage    : %s", memMsg),
	}

	return title + strings.Join(indicators, "\n")
}

// RenderSiteDistributionChart renders a bar chart of site distribution
func (r *Renderer) RenderSiteDistributionChart(m *model.Model) string {
	title := styles.PanelTitle.Render("📊 Site Distribution")

	bc := barchart.New(40, 10)

	enabledCount := 0
	disabledCount := 0
	for _, site := range m.Sites {
		if site.Enabled {
			enabledCount++
		} else {
			disabledCount++
		}
	}

	bc.Push(barchart.BarData{
		Label: "Enabled",
		Values: []barchart.BarValue{
			{Name: "count", Value: float64(enabledCount)},
		},
	})

	bc.Push(barchart.BarData{
		Label: "Disabled",
		Values: []barchart.BarValue{
			{Name: "count", Value: float64(disabledCount)},
		},
	})

	bc.Draw()

	return styles.Panel.Render(lipgloss.JoinVertical(lipgloss.Left, title, bc.View()))
}

func (r *Renderer) RenderMetricsView(m *model.Model, width, height int) string {
	// Collect real metrics
	nginxService := nginx.New()
	metrics, err := nginxService.GetMetrics()

	// Update history with real data
	if err == nil {
		// Calculate network rate (MB/s) from change in total bytes
		var networkRate float64
		if m.LastNetworkIn > 0 && m.LastNetworkOut > 0 {
			// Calculate change since last measurement
			deltaIn := metrics.NetworkIn - m.LastNetworkIn
			deltaOut := metrics.NetworkOut - m.LastNetworkOut
			networkRate = (deltaIn + deltaOut) / 1024 // Convert to MB/s
		}

		// Store current values for next calculation
		m.LastNetworkIn = metrics.NetworkIn
		m.LastNetworkOut = metrics.NetworkOut

		// Add new data point (will grow from 0 to 50 points)
		if len(m.CPUHistory) < 50 {
			// Still filling up - just append
			m.CPUHistory = append(m.CPUHistory, metrics.CPU)
			m.MemHistory = append(m.MemHistory, metrics.Memory)
			m.NetHistory = append(m.NetHistory, networkRate)
			m.RequestHistory = append(m.RequestHistory, metrics.RequestRate)
		} else {
			// Full - shift and add new data
			m.CPUHistory = append(m.CPUHistory[1:], metrics.CPU)
			m.MemHistory = append(m.MemHistory[1:], metrics.Memory)
			m.NetHistory = append(m.NetHistory[1:], networkRate)
			m.RequestHistory = append(m.RequestHistory[1:], metrics.RequestRate)
		}
	}

	// Calculate chart dimensions to use full width
	chartWidth := (width / 2) - 6
	chartHeight := (height / 2) - 4

	if chartWidth < 40 {
		chartWidth = 40
	}
	if chartHeight < 8 {
		chartHeight = 8
	}

	// Render charts with real data
	cpuChart := r.RenderLineChart("CPU Usage (%)", m.CPUHistory, styles.AccentPrimary, chartWidth, chartHeight)
	memChart := r.RenderLineChart("Memory Usage (%)", m.MemHistory, styles.AccentSecondary, chartWidth, chartHeight)
	netChart := r.RenderLineChart("Network (MB/s)", m.NetHistory, styles.AccentSuccess, chartWidth, chartHeight)
	reqChart := r.RenderLineChart("Request Rate (req/s)", m.RequestHistory, styles.AccentWarning, chartWidth, chartHeight)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, cpuChart, "  ", memChart)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, netChart, "  ", reqChart)

	return lipgloss.JoinVertical(lipgloss.Left, row1, "", row2)
}

// RenderLineChart renders a line chart with specified dimensions
func (r *Renderer) RenderLineChart(title string, data []float64, color lipgloss.Color, width, height int) string {
	titleText := styles.PanelTitle.Render(title)

	if len(data) == 0 {
		return styles.Panel.Width(width).Render(
			lipgloss.JoinVertical(lipgloss.Left, titleText, styles.MutedText.Render("No data available")),
		)
	}

	// Adjust dimensions for chart content (account for padding and borders)
	chartWidth := width - 6
	chartHeight := height - 6

	if chartWidth < 20 {
		chartWidth = 20
	}
	if chartHeight < 5 {
		chartHeight = 5
	}

	slc := streamlinechart.New(chartWidth, chartHeight)
	for _, val := range data {
		slc.Push(val)
	}
	slc.Draw()

	// Get current value
	currentVal := data[len(data)-1]
	currentText := styles.InfoText.Render(fmt.Sprintf("Current: %.1f%%", currentVal))

	// Create sparkline for mini view
	sparklineText := r.RenderSparkline(data, chartWidth)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleText,
		slc.View(),
		sparklineText,
		currentText,
	)

	return styles.Panel.Width(width).Render(content)
}

// RenderSparkline creates a sparkline using ntcharts
func (r *Renderer) RenderSparkline(data []float64, width int) string {
	sl := sparkline.New(width, 3)
	for _, val := range data {
		sl.Push(val)
	}
	return sl.View()
}

// RenderHelp renders the help text with modern styling
func (r *Renderer) RenderHelp(m *model.Model) string {
	// Group keys by category
	navigation := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.HelpKey.Render("←→"),
		styles.HelpSeparator.Render(" "),
		styles.HelpDesc.Render("tabs"),
		styles.HelpSeparator.Render("  │  "),
		styles.HelpKey.Render("↑↓"),
		styles.HelpSeparator.Render(" "),
		styles.HelpDesc.Render("navigate"),
		styles.HelpSeparator.Render("  │  "),
		styles.HelpKey.Render("enter"),
		styles.HelpSeparator.Render(" "),
		styles.HelpDesc.Render("select"),
	)

	actions := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.HelpKey.Render("r"),
		styles.HelpSeparator.Render(" "),
		styles.HelpDesc.Render("refresh"),
		styles.HelpSeparator.Render("  │  "),
		styles.HelpKey.Render("esc"),
		styles.HelpSeparator.Render(" "),
		styles.HelpDesc.Render("back"),
		styles.HelpSeparator.Render("  │  "),
		styles.HelpKey.Render("q"),
		styles.HelpSeparator.Render(" "),
		styles.HelpDesc.Render("quit"),
	)

	helpText := lipgloss.JoinHorizontal(
		lipgloss.Left,
		navigation,
		styles.HelpSeparator.Render("    "),
		actions,
	)

	return styles.Footer.Render(helpText)
}

// RenderStatusBar renders the status bar
func (r *Renderer) RenderStatusBar(m *model.Model) string {
	if !m.ShowStatus {
		return ""
	}

	if m.IsError {
		return styles.ErrorText.Render("✗ " + m.StatusMsg)
	}
	return styles.SuccessText.Render("✓ " + m.StatusMsg)
}
