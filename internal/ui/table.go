package ui

import (
	"github.com/76creates/stickers/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/aitmiloud/ngxtui/internal/styles"
)

// CreateSitesTable creates a new stickers table for sites
func CreateSitesTable(sites []model.Site, width, height int) *table.Table {
	// Create table headers
	headers := []string{
		"Site Name",
		"Status",
		"Port",
		"SSL",
		"Uptime",
	}

	// Create table
	t := table.NewTable(0, 0, headers)

	// Prepare rows as [][]any
	var rows [][]any
	for _, site := range sites {
		status := "○ Inactive"
		if site.Enabled {
			status = "● Active"
		}

		ssl := "No"
		if site.SSL {
			ssl = "Yes"
		}

		row := []any{
			site.Name,
			status,
			site.Port,
			ssl,
			site.Uptime,
		}
		rows = append(rows, row)
	}

	// Add all rows at once
	if len(rows) > 0 {
		t.MustAddRows(rows)
	}

	// Set table dimensions
	if width > 0 {
		t.SetWidth(width - 4) // Account for borders
	}
	if height > 0 {
		t.SetHeight(height)
	}

	// Apply modern styling
	t.SetStyles(map[table.StyleKey]lipgloss.Style{
		table.StyleKeyHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.AccentPrimary).
			Padding(0, 1),
	})

	return t
}

// RenderSitesTableStickers renders the sites table using Stickers
func (r *Renderer) RenderSitesTableStickers(m *model.Model, width, height int) string {
	if len(m.Sites) == 0 {
		emptyMsg := styles.MutedText.Render("No NGINX sites configured")
		return styles.Panel.
			Width(width - 4).
			Height(height).
			Render(lipgloss.Place(width-8, height-4, lipgloss.Center, lipgloss.Center, emptyMsg))
	}

	// Create or update table
	t := CreateSitesTable(m.Sites, width, height)

	// Set cursor position by moving down m.Selected times
	for i := 0; i < m.Selected && i < len(m.Sites); i++ {
		t.CursorDown()
	}

	title := styles.PanelTitle.Render("🌐 NGINX Sites")

	tableView := t.Render()

	// Add helpful hint
	hint := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render("Press Enter to manage site • ↑↓ to navigate • r to refresh")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		tableView,
		"",
		hint,
	)

	return styles.Panel.
		Width(width - 2).
		Render(content)
}
