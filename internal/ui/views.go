package ui

import (
	"fmt"
	"strings"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/NimbleMarkets/ntcharts/linechart/streamlinechart"
	"github.com/NimbleMarkets/ntcharts/sparkline"
	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/aitmiloud/ngxtui/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

// Renderer handles all view rendering
type Renderer struct{}

// New creates a new renderer
func New() *Renderer {
	return &Renderer{}
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
func (r *Renderer) RenderSitesTable(m *model.Model) string {
	// Use the stickers table renderer with default dimensions
	return r.RenderSitesTableStickers(m, 100, 20)
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

// RenderLogsView renders the logs view with enhanced styling
func (r *Renderer) RenderLogsView(m *model.Model) string {
	title := styles.PanelTitle.Render("📋 Access Logs")

	// Legend for status codes
	legend := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.StatusSuccess.Render("● 2xx Success"),
		"  ",
		styles.StatusInfo.Render("● 3xx Redirect"),
		"  ",
		styles.StatusWarning.Render("● 4xx Client Error"),
		"  ",
		styles.StatusError.Render("● 5xx Server Error"),
	)

	// Helper styles for log formatting
	secondaryStyle := lipgloss.NewStyle().Foreground(styles.TextSecondary)
	primaryStyle := lipgloss.NewStyle().Foreground(styles.TextPrimary)
	mutedStyle := lipgloss.NewStyle().Foreground(styles.TextMuted)

	// Simulated log entries with enhanced color coding
	logs := []string{
		styles.StatusSuccess.Render("● ") + secondaryStyle.Render("192.168.1.100 - [04/Nov/2024:15:42:07] ") + primaryStyle.Render("GET /api/users") + styles.StatusSuccess.Render(" 200") + mutedStyle.Render(" 1.2KB"),
		styles.StatusInfo.Render("● ") + secondaryStyle.Render("192.168.1.101 - [04/Nov/2024:15:42:08] ") + primaryStyle.Render("POST /api/login") + styles.StatusInfo.Render(" 201") + mutedStyle.Render(" 567B"),
		styles.StatusWarning.Render("● ") + secondaryStyle.Render("192.168.1.102 - [04/Nov/2024:15:42:09] ") + primaryStyle.Render("GET /missing") + styles.StatusWarning.Render(" 404") + mutedStyle.Render(" 89B"),
		styles.StatusSuccess.Render("● ") + secondaryStyle.Render("192.168.1.103 - [04/Nov/2024:15:42:10] ") + primaryStyle.Render("GET /api/products") + styles.StatusSuccess.Render(" 200") + mutedStyle.Render(" 2.3KB"),
		styles.StatusError.Render("● ") + secondaryStyle.Render("192.168.1.104 - [04/Nov/2024:15:42:11] ") + primaryStyle.Render("POST /api/order") + styles.StatusError.Render(" 500") + mutedStyle.Render(" 123B"),
		styles.StatusSuccess.Render("● ") + secondaryStyle.Render("192.168.1.105 - [04/Nov/2024:15:42:12] ") + primaryStyle.Render("GET /health") + styles.StatusSuccess.Render(" 200") + mutedStyle.Render(" 45B"),
		styles.StatusInfo.Render("● ") + secondaryStyle.Render("192.168.1.106 - [04/Nov/2024:15:42:13] ") + primaryStyle.Render("PUT /api/user/123") + styles.StatusInfo.Render(" 204") + mutedStyle.Render(" 0B"),
		styles.StatusWarning.Render("● ") + secondaryStyle.Render("192.168.1.107 - [04/Nov/2024:15:42:14] ") + primaryStyle.Render("GET /old-page") + styles.StatusWarning.Render(" 404") + mutedStyle.Render(" 234B"),
		styles.StatusSuccess.Render("● ") + secondaryStyle.Render("192.168.1.108 - [04/Nov/2024:15:42:15] ") + primaryStyle.Render("GET /api/status") + styles.StatusSuccess.Render(" 200") + mutedStyle.Render(" 512B"),
		styles.StatusSuccess.Render("● ") + secondaryStyle.Render("192.168.1.109 - [04/Nov/2024:15:42:16] ") + primaryStyle.Render("GET /assets/logo.png") + styles.StatusSuccess.Render(" 200") + mutedStyle.Render(" 15KB"),
	}

	content := strings.Join(logs, "\n")
	m.Viewport.SetContent(content)

	divider := styles.Divider.Render(strings.Repeat("─", 80))

	return styles.Panel.Render(
		lipgloss.JoinVertical(lipgloss.Left, title, legend, divider, m.Viewport.View()),
	)
}

// RenderStatsView renders the statistics view with enhanced layout
func (r *Renderer) RenderStatsView(m *model.Model, width int) string {
	totalSites := len(m.Sites)
	enabledSites := 0
	for _, site := range m.Sites {
		if site.Enabled {
			enabledSites++
		}
	}

	// Calculate percentage
	percentage := "0%"
	if totalSites > 0 {
		percentage = fmt.Sprintf("%.1f%%", float64(enabledSites)/float64(totalSites)*100)
	}

	// Metric cards with icons
	stats := []string{
		r.RenderMetricCard("🌐 Total Sites", fmt.Sprintf("%d", totalSites), "All configured"),
		r.RenderMetricCard("✓ Active Sites", fmt.Sprintf("%d", enabledSites), percentage+" enabled"),
		r.RenderMetricCard("⚡ Request Rate", "1.2k/s", "Avg per second"),
		r.RenderMetricCard("⏱️  Uptime", "99.9%", "Last 30 days"),
	}

	statsRow := lipgloss.JoinHorizontal(lipgloss.Top, stats...)

	// Add section title
	chartTitle := styles.PanelTitle.Render("📊 Site Distribution")

	chart := r.RenderSiteDistributionChart(m)

	return lipgloss.JoinVertical(lipgloss.Left, statsRow, "", chartTitle, chart)
}

// RenderMetricCard renders a metric card with modern styling
func (r *Renderer) RenderMetricCard(label, value, subtext string) string {
	title := styles.CardTitle.Render(label)

	// Large, centered value
	val := styles.CardValue.
		Width(20).
		Render(value)

	sub := styles.CardSubtext.
		Width(20).
		Render(subtext)

	content := lipgloss.JoinVertical(lipgloss.Center, title, "", val, sub)
	return styles.Card.Width(24).Height(6).Render(content)
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

// RenderMetricsView renders the metrics view with charts using full width
func (r *Renderer) RenderMetricsView(m *model.Model, width, height int) string {
	// Calculate chart dimensions to use full width
	chartWidth := (width / 2) - 6   // Split width in half, account for padding/borders
	chartHeight := (height / 2) - 4 // Split height in half, account for spacing

	if chartWidth < 40 {
		chartWidth = 40
	}
	if chartHeight < 8 {
		chartHeight = 8
	}

	cpuChart := r.RenderLineChart("CPU Usage", m.CPUHistory, styles.AccentPrimary, chartWidth, chartHeight)
	memChart := r.RenderLineChart("Memory Usage", m.MemHistory, styles.AccentSecondary, chartWidth, chartHeight)
	netChart := r.RenderLineChart("Network Traffic", m.NetHistory, styles.AccentSuccess, chartWidth, chartHeight)
	reqChart := r.RenderLineChart("Request Rate", m.RequestHistory, styles.AccentWarning, chartWidth, chartHeight)

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
