package app

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/aitmiloud/ngxtui/internal/styles"
	"github.com/aitmiloud/ngxtui/internal/ui"
)

// View renders the entire application view
func View(m model.Model) string {
	if m.Quitting {
		return styles.SuccessText.Render("Thanks for using NgxTUI! 👋\n")
	}

	renderer := ui.New()

	// Header
	header := styles.Header.Render("🚀 NGINX Terminal UI Manager")

	// Tabs
	tabs := renderer.RenderTabs(&m)

	// Content based on active tab
	var content string
	if m.MenuMode {
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			renderer.RenderSitesTable(&m),
			"  ",
			renderer.RenderActionMenu(&m),
		)
	} else {
		switch m.ActiveTab {
		case model.SitesTab:
			content = renderer.RenderSitesTable(&m)
		case model.LogsTab:
			content = renderer.RenderLogsView(&m)
		case model.StatsTab:
			content = renderer.RenderStatsView(&m)
		case model.MetricsTab:
			content = renderer.RenderMetricsView(&m)
		}
	}

	// Status bar
	statusBar := renderer.RenderStatusBar(&m)

	// Help
	help := renderer.RenderHelp(&m)

	// Combine all sections
	sections := []string{header, tabs, content}
	if statusBar != "" {
		sections = append(sections, statusBar)
	}
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
