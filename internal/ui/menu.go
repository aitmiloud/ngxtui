package ui

import (
	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/charmbracelet/lipgloss"
)

// RenderSitesWithMenu renders the sites table with action menu side by side
func (r *Renderer) RenderSitesWithMenu(m *model.Model, width, height int) string {
	table := r.RenderSitesTable(m)
	menu := r.RenderActionMenu(m)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		table,
		"  ",
		menu,
	)
}
