package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/aitmiloud/ngxtui/internal/styles"
	"github.com/aitmiloud/ngxtui/internal/ui"
)

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// View renders the entire application view
func View(m model.Model) string {
	if m.Quitting {
		return styles.SuccessText.Render("Thanks for using NgxTUI! 👋\n")
	}

	renderer := ui.New()

	// Calculate dimensions
	width := m.Width
	height := m.Height
	if width == 0 {
		width = 100
	}
	if height == 0 {
		height = 30
	}

	// Header with version and status (fixed height: 3 lines)
	headerLeft := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.TextPrimary).
		Render("🚀 NGINX TUI")

	headerRight := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render("v1.0.0")

	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerLeft,
		strings.Repeat(" ", max(0, width-lipgloss.Width(headerLeft)-lipgloss.Width(headerRight)-6)),
		headerRight,
	)

	header := styles.Header.Width(width - 2).Render(headerContent)

	// Tabs (fixed height: 2 lines)
	tabs := renderer.RenderTabs(&m)

	// Footer/Help (fixed height: 2 lines)
	help := renderer.RenderHelp(&m)

	// Status bar (1 line if present)
	statusBar := renderer.RenderStatusBar(&m)
	statusHeight := 0
	if statusBar != "" {
		statusHeight = 1
	}

	// Calculate available height for content
	// Total: header(3) + tabs(2) + status(0-1) + help(2) = 7-8 lines
	contentHeight := height - 7 - statusHeight
	if contentHeight < 10 {
		contentHeight = 10
	}

	// Content based on active tab (fills remaining space)
	var content string
	if m.MenuMode {
		content = renderer.RenderSitesWithMenu(&m, width, contentHeight)
	} else {
		switch m.ActiveTab {
		case model.SitesTab:
			content = renderer.RenderSitesTable(&m)
		case model.LogsTab:
			content = renderer.RenderLogsView(&m)
		case model.StatsTab:
			content = renderer.RenderStatsView(&m, width)
		case model.MetricsTab:
			content = renderer.RenderMetricsView(&m, width, contentHeight)
		}
	}

	// Ensure content doesn't exceed available height and fills it
	contentLines := strings.Split(content, "\n")
	currentHeight := len(contentLines)

	// Truncate if too tall
	if currentHeight > contentHeight {
		contentLines = contentLines[:contentHeight]
		content = strings.Join(contentLines, "\n")
	} else if currentHeight < contentHeight {
		// Add padding to fill space
		padding := strings.Repeat("\n", contentHeight-currentHeight)
		content = content + padding
	}

	// Combine all sections in fixed order
	sections := []string{header, tabs, content}
	if statusBar != "" {
		sections = append(sections, statusBar)
	}
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
