package app

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/aitmiloud/ngxtui/internal/nginx"
	"github.com/aitmiloud/ngxtui/internal/styles"
)

// InitialModel creates the initial model state
func InitialModel() model.Model {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.AccentPrimary)

	// Initialize progress bar
	prog := progress.New(progress.WithDefaultGradient())

	// Initialize viewport for logs
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.BorderSubtle).
		Padding(1)

	// Load initial sites
	nginxService := nginx.New()
	sites, err := nginxService.ListSites()
	if err != nil {
		sites = []model.Site{}
	}

	// Create table
	columns := []table.Column{
		{Title: "Site Name", Width: 30},
		{Title: "Status", Width: 15},
		{Title: "Port", Width: 10},
		{Title: "SSL", Width: 10},
		{Title: "Uptime", Width: 15},
	}

	rows := []table.Row{}
	for _, site := range sites {
		status := "Disabled"
		if site.Enabled {
			status = "Enabled"
		}
		ssl := "No"
		if site.SSL {
			ssl = "Yes"
		}
		rows = append(rows, table.Row{
			site.Name,
			status,
			site.Port,
			ssl,
			site.Uptime,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderSubtle).
		BorderBottom(true).
		Bold(true).
		Foreground(styles.AccentPrimary)
	tableStyle.Selected = tableStyle.Selected.
		Foreground(styles.TextPrimary).
		Background(styles.AccentSecondary).
		Bold(true)
	t.SetStyles(tableStyle)

	// Initialize metric history with random data
	cpuHistory := make([]float64, 50)
	memHistory := make([]float64, 50)
	netHistory := make([]float64, 50)
	requestHistory := make([]float64, 50)

	for i := 0; i < 50; i++ {
		cpuHistory[i] = 20 + rand.Float64()*30
		memHistory[i] = 40 + rand.Float64()*20
		netHistory[i] = 10 + rand.Float64()*40
		requestHistory[i] = 50 + rand.Float64()*50
	}

	return model.Model{
		Sites:          sites,
		Table:          t,
		Cursor:         0,
		Selected:       -1,
		MenuMode:       false,
		ActiveTab:      model.SitesTab,
		Spinner:        s,
		Loading:        false,
		Viewport:       vp,
		CPUHistory:     cpuHistory,
		MemHistory:     memHistory,
		NetHistory:     netHistory,
		RequestHistory: requestHistory,
		Progress:       prog,
		LastUpdate:     time.Now(),
	}
}

// TickCmd returns a command that sends tick messages
func TickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return model.TickMsg(t)
	})
}
