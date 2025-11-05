package ui

import (
	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/charmbracelet/lipgloss"
)

// RenderSitesWithMenu renders the sites table with action menu side by side
func (r *Renderer) RenderSitesWithMenu(m *model.Model, width, height int) string {
	// Split width between table and menu (2/3 for table)
	tableWidth := (width * 2) / 3
	
	table := r.RenderSitesTable(m, tableWidth, height)
	menu := r.RenderActionMenu(m)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		table,
		"  ",
		menu,
	)
}
