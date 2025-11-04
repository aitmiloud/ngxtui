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

// RenderTabs renders the tab bar
func (r *Renderer) RenderTabs(m *model.Model) string {
	tabs := []string{"Sites", "Logs", "Stats", "Metrics"}
	var renderedTabs []string

	for i, tab := range tabs {
		if model.TabType(i) == m.ActiveTab {
			renderedTabs = append(renderedTabs, styles.ActiveTab.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, styles.InactiveTab.Render(tab))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// RenderSitesTable renders the sites table view
func (r *Renderer) RenderSitesTable(m *model.Model) string {
	return m.Table.View()
}

// RenderActionMenu renders the action menu for a selected site
func (r *Renderer) RenderActionMenu(m *model.Model) string {
	if m.Selected < 0 || m.Selected >= len(m.Sites) {
		return ""
	}

	site := m.Sites[m.Selected]
	title := styles.PanelTitle.Render(fmt.Sprintf("Actions for: %s", site.Name))

	actions := []string{
		"Enable Site",
		"Disable Site",
		"Test Configuration",
		"Reload NGINX",
		"View Logs",
		"Back",
	}

	var items []string
	for i, action := range actions {
		if i == m.Cursor {
			items = append(items, styles.SelectedAction.Render("▸ "+action))
		} else {
			items = append(items, styles.ActionItem.Render("  "+action))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return styles.Panel.Render(lipgloss.JoinVertical(lipgloss.Left, title, content))
}

// RenderStatusBadge renders a status badge for enabled/disabled state
func (r *Renderer) RenderStatusBadge(enabled bool) string {
	if enabled {
		return styles.EnabledBadge.Render("ENABLED")
	}
	return styles.DisabledBadge.Render("DISABLED")
}

// RenderLogsView renders the logs view
func (r *Renderer) RenderLogsView(m *model.Model) string {
	title := styles.PanelTitle.Render("📝 Access Logs")

	// Simulated log entries with color coding
	logs := []string{
		styles.SuccessText.Render("192.168.1.100 - - [04/Nov/2024:15:42:07 +0100] \"GET /api/users HTTP/1.1\" 200 1234"),
		styles.InfoText.Render("192.168.1.101 - - [04/Nov/2024:15:42:08 +0100] \"POST /api/login HTTP/1.1\" 201 567"),
		styles.WarningText.Render("192.168.1.102 - - [04/Nov/2024:15:42:09 +0100] \"GET /missing HTTP/1.1\" 404 89"),
		styles.SuccessText.Render("192.168.1.103 - - [04/Nov/2024:15:42:10 +0100] \"GET /api/products HTTP/1.1\" 200 2345"),
		styles.ErrorText.Render("192.168.1.104 - - [04/Nov/2024:15:42:11 +0100] \"POST /api/order HTTP/1.1\" 500 123"),
		styles.SuccessText.Render("192.168.1.105 - - [04/Nov/2024:15:42:12 +0100] \"GET /health HTTP/1.1\" 200 45"),
		styles.InfoText.Render("192.168.1.106 - - [04/Nov/2024:15:42:13 +0100] \"PUT /api/user/123 HTTP/1.1\" 204 0"),
		styles.WarningText.Render("192.168.1.107 - - [04/Nov/2024:15:42:14 +0100] \"GET /old-page HTTP/1.1\" 404 234"),
	}

	content := strings.Join(logs, "\n")
	m.Viewport.SetContent(content)

	return styles.Panel.Render(lipgloss.JoinVertical(lipgloss.Left, title, m.Viewport.View()))
}

// RenderStatsView renders the statistics view
func (r *Renderer) RenderStatsView(m *model.Model) string {
	totalSites := len(m.Sites)
	enabledSites := 0
	for _, site := range m.Sites {
		if site.Enabled {
			enabledSites++
		}
	}

	stats := []string{
		r.RenderMetricCard("Total Sites", fmt.Sprintf("%d", totalSites), "All configured sites"),
		r.RenderMetricCard("Active Sites", fmt.Sprintf("%d", enabledSites), "Currently enabled"),
		r.RenderMetricCard("Request Rate", "1.2k/s", "Average requests per second"),
		r.RenderMetricCard("Uptime", "99.9%", "Last 30 days"),
	}

	statsRow := lipgloss.JoinHorizontal(lipgloss.Top, stats...)

	chart := r.RenderSiteDistributionChart(m)

	return lipgloss.JoinVertical(lipgloss.Left, statsRow, chart)
}

// RenderMetricCard renders a metric card
func (r *Renderer) RenderMetricCard(label, value, subtext string) string {
	title := styles.PanelTitle.Render(label)
	val := styles.InfoText.Render(value)
	sub := styles.MutedText.Render(subtext)

	content := lipgloss.JoinVertical(lipgloss.Left, title, val, sub)
	return styles.Panel.Width(25).Render(content)
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

// RenderMetricsView renders the metrics view with charts
func (r *Renderer) RenderMetricsView(m *model.Model) string {
	cpuChart := r.RenderLineChart("CPU Usage", m.CPUHistory, styles.AccentPrimary)
	memChart := r.RenderLineChart("Memory Usage", m.MemHistory, styles.AccentSecondary)
	netChart := r.RenderLineChart("Network Traffic", m.NetHistory, styles.AccentSuccess)
	reqChart := r.RenderLineChart("Request Rate", m.RequestHistory, styles.AccentWarning)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, cpuChart, memChart)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, netChart, reqChart)

	return lipgloss.JoinVertical(lipgloss.Left, row1, row2)
}

// RenderLineChart renders a line chart
func (r *Renderer) RenderLineChart(title string, data []float64, color lipgloss.Color) string {
	titleText := styles.PanelTitle.Render(title)

	if len(data) == 0 {
		return styles.Panel.Width(40).Render(
			lipgloss.JoinVertical(lipgloss.Left, titleText, styles.MutedText.Render("No data available")),
		)
	}

	slc := streamlinechart.New(35, 8)
	for _, val := range data {
		slc.Push(val)
	}
	slc.Draw()

	// Get current value
	currentVal := data[len(data)-1]
	currentText := styles.InfoText.Render(fmt.Sprintf("Current: %.1f%%", currentVal))

	// Create sparkline for mini view
	sparklineText := r.RenderSparkline(data, 35)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleText,
		slc.View(),
		sparklineText,
		currentText,
	)

	return styles.Panel.Width(42).Render(content)
}

// RenderSparkline creates a sparkline using ntcharts
func (r *Renderer) RenderSparkline(data []float64, width int) string {
	sl := sparkline.New(width, 3)
	for _, val := range data {
		sl.Push(val)
	}
	return sl.View()
}

// RenderHelp renders the help text
func (r *Renderer) RenderHelp(m *model.Model) string {
	keys := []string{
		styles.HelpKey.Render("←/→") + styles.HelpSeparator.Render(" • ") + styles.HelpDesc.Render("switch tabs"),
		styles.HelpKey.Render("↑/↓") + styles.HelpSeparator.Render(" • ") + styles.HelpDesc.Render("navigate"),
		styles.HelpKey.Render("enter") + styles.HelpSeparator.Render(" • ") + styles.HelpDesc.Render("select"),
		styles.HelpKey.Render("esc") + styles.HelpSeparator.Render(" • ") + styles.HelpDesc.Render("back"),
		styles.HelpKey.Render("r") + styles.HelpSeparator.Render(" • ") + styles.HelpDesc.Render("refresh"),
		styles.HelpKey.Render("q") + styles.HelpSeparator.Render(" • ") + styles.HelpDesc.Render("quit"),
	}

	return styles.MutedText.Render(strings.Join(keys, "  "))
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
