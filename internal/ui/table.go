package ui

import (
	"github.com/76creates/stickers/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/aitmiloud/ngxtui/internal/styles"
)

// CreateSitesTable creates a new stickers table for sites
func CreateSitesTable(sites []model.Site, width, height int) *table.Table {
	// Create table headers with icons
	headers := []string{
		"🌐 Site Name",
		"⚡ Status",
		"🔌 Port",
		"🔒 SSL",
		"⏱️  Uptime",
	}

	// Create table
	t := table.NewTable(0, 0, headers)

	// Prepare rows as [][]any - keep it simple without ANSI codes
	var rows [][]any
	for _, site := range sites {
		// Status indicator
		status := "○ Inactive"
		if site.Enabled {
			status = "● Active"
		}

		// SSL indicator
		ssl := "✗ No"
		if site.SSL {
			ssl = "✓ Yes"
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

	// Set table dimensions to be responsive
	if width > 0 {
		t.SetWidth(width - 8) // Account for panel padding and borders
	}
	if height > 0 {
		t.SetHeight(height - 4) // Account for hint text
	}

	// Apply enhanced styling
	t.SetStyles(map[table.StyleKey]lipgloss.Style{
		table.StyleKeyHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.AccentPrimary).
			Background(lipgloss.Color("#1E293B")).
			Padding(0, 2),
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

	// Create table
	t := CreateSitesTable(m.Sites, width, height)

	// Set cursor position
	for i := 0; i < m.Selected && i < len(m.Sites); i++ {
		t.CursorDown()
	}

	tableView := t.Render()

	// Add helpful hint
	hint := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render("Press Enter to manage site • ↑↓ to navigate • r to refresh")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		tableView,
		"",
		hint,
	)

	return styles.Panel.
		Width(width - 2).
		Render(content)
}
