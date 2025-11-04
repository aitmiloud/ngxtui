package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/NimbleMarkets/ntcharts/linechart/streamlinechart"
	"github.com/NimbleMarkets/ntcharts/sparkline"
)

// Styles
var (
	// Color palette - dark theme
	accentPrimary   = lipgloss.Color("#00D9FF")
	accentSecondary = lipgloss.Color("#7C3AED")
	accentSuccess   = lipgloss.Color("#10B981")
	accentWarning   = lipgloss.Color("#F59E0B")
	accentDanger    = lipgloss.Color("#EF4444")
	textPrimary     = lipgloss.Color("#F9FAFB")
	textSecondary   = lipgloss.Color("#9CA3AF")
	textMuted       = lipgloss.Color("#6B7280")
	bgPrimary       = lipgloss.Color("#111827")
	bgSecondary     = lipgloss.Color("#1F2937")
	bgTertiary      = lipgloss.Color("#374151")
	borderSubtle    = lipgloss.Color("#4B5563")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Foreground(textPrimary)

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentPrimary).
			Background(bgSecondary).
			Padding(0, 2).
			MarginBottom(1)

	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textPrimary).
			Background(accentSecondary).
			Padding(0, 3).
			MarginRight(1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(textSecondary).
				Background(bgTertiary).
				Padding(0, 3).
				MarginRight(1)

	// Card/Panel styles
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderSubtle).
			Padding(1, 2).
			MarginBottom(1)

	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentPrimary).
			MarginBottom(1)

	// Status badge styles
	enabledBadge = lipgloss.NewStyle().
			Foreground(textPrimary).
			Background(accentSuccess).
			Padding(0, 1).
			Bold(true)

	disabledBadge = lipgloss.NewStyle().
			Foreground(textPrimary).
			Background(textMuted).
			Padding(0, 1)

	// Action styles
	actionItemStyle = lipgloss.NewStyle().
			Foreground(textSecondary).
			Padding(0, 2).
			MarginBottom(0)

	selectedActionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(textPrimary).
				Background(accentSecondary).
				Padding(0, 2).
				MarginBottom(0)

	// Info styles
	labelStyle = lipgloss.NewStyle().
			Foreground(textSecondary).
			Width(16)

	valueStyle = lipgloss.NewStyle().
			Foreground(textPrimary).
			Bold(true)

	// Status message styles
	successMsgStyle = lipgloss.NewStyle().
			Foreground(accentSuccess).
			Bold(true).
			Padding(0, 1)

	errorMsgStyle = lipgloss.NewStyle().
			Foreground(accentDanger).
			Bold(true).
			Padding(0, 1)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			MarginTop(1)

	// Metric card styles
	metricCardStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderSubtle).
			Padding(1, 2).
			Width(24).
			Height(6)

	metricLabelStyle = lipgloss.NewStyle().
				Foreground(textSecondary).
				MarginBottom(1)

	metricValueStyle = lipgloss.NewStyle().
				Foreground(accentPrimary).
				Bold(true).
				Width(0)

	metricSubtextStyle = lipgloss.NewStyle().
				Foreground(textMuted).
				Width(0)
)

type site struct {
	Name    string
	Enabled bool
	Port    string
	SSL     bool
	Uptime  string
}

type tabType int

const (
	sitesTab tabType = iota
	logsTab
	statsTab
	metricsTab
)

type statusMsg struct {
	message string
	isError bool
}

type tickMsg time.Time

type model struct {
	sites          []site
	table          table.Model
	cursor         int
	selected       int
	menuMode       bool
	activeTab      tabType
	quitting       bool
	err            error
	statusMsg      string
	showStatus     bool
	isError        bool
	spinner        spinner.Model
	loading        bool
	width          int
	height         int
	viewport       viewport.Model
	cpuHistory     []float64
	memHistory     []float64
	netHistory     []float64
	requestHistory []float64
	progress       progress.Model
	lastUpdate     time.Time
}

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Enter   key.Binding
	Back    key.Binding
	Quit    key.Binding
	Tab     key.Binding
	Refresh key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev tab"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next tab"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch tab"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func initialModel() model {
	sites, err := listNginxSites()

	// Create table
	columns := []table.Column{
		{Title: "SITE", Width: 30},
		{Title: "STATUS", Width: 12},
		{Title: "PORT", Width: 8},
		{Title: "SSL", Width: 6},
		{Title: "UPTIME", Width: 12},
	}

	rows := []table.Row{}
	for _, s := range sites {
		status := "DISABLED"
		if s.Enabled {
			status = "ENABLED"
		}
		ssl := "NO"
		if s.SSL {
			ssl = "YES"
		}
		rows = append(rows, table.Row{s.Name, status, s.Port, ssl, s.Uptime})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderSubtle).
		BorderBottom(true).
		Bold(true).
		Foreground(accentPrimary)
	s.Selected = s.Selected.
		Foreground(textPrimary).
		Background(accentSecondary).
		Bold(true)
	t.SetStyles(s)

	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(accentPrimary)

	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	// Initialize history data
	cpuHistory := make([]float64, 60)
	memHistory := make([]float64, 60)
	netHistory := make([]float64, 60)
	requestHistory := make([]float64, 60)

	// Generate initial data
	for i := 0; i < 60; i++ {
		cpuHistory[i] = 30 + rand.Float64()*40
		memHistory[i] = 40 + rand.Float64()*30
		netHistory[i] = 20 + rand.Float64()*50
		requestHistory[i] = 100 + rand.Float64()*200
	}

	var errMsg error
	if err != nil {
		errMsg = err
	}

	return model{
		sites:          sites,
		table:          t,
		cursor:         0,
		selected:       -1,
		menuMode:       false,
		activeTab:      sitesTab,
		err:            errMsg,
		spinner:        sp,
		progress:       prog,
		cpuHistory:     cpuHistory,
		memHistory:     memHistory,
		netHistory:     netHistory,
		requestHistory: requestHistory,
		lastUpdate:     time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tickCmd())
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func listNginxSites() ([]site, error) {
	// Uncomment for real implementation
	/*
		dir := "/etc/nginx/sites-available"
		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		var sites []site
		for _, f := range files {
			enabled := false
			link := "/etc/nginx/sites-enabled/" + f.Name()
			if _, err := os.Lstat(link); err == nil {
				enabled = true
			}
			sites = append(sites, site{Name: f.Name(), Enabled: enabled})
		}
		return sites, nil
	*/

	sites := []site{
		{Name: "example.com", Enabled: true, Port: "80", SSL: true, Uptime: "45d 12h"},
		{Name: "test.com", Enabled: false, Port: "8080", SSL: false, Uptime: "N/A"},
		{Name: "myapp.local", Enabled: true, Port: "3000", SSL: true, Uptime: "12d 3h"},
		{Name: "demo.site", Enabled: false, Port: "80", SSL: false, Uptime: "N/A"},
		{Name: "api.service.com", Enabled: true, Port: "8000", SSL: true, Uptime: "89d 23h"},
		{Name: "staging.app", Enabled: true, Port: "4000", SSL: true, Uptime: "2d 8h"},
		{Name: "dev.local", Enabled: false, Port: "5000", SSL: false, Uptime: "N/A"},
	}
	return sites, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetHeight(m.height - 18)
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = m.height - 10

	case tickMsg:
		// Update metrics data
		m.cpuHistory = append(m.cpuHistory[1:], 30+rand.Float64()*40)
		m.memHistory = append(m.memHistory[1:], 40+rand.Float64()*30)
		m.netHistory = append(m.netHistory[1:], 20+rand.Float64()*50)
		m.requestHistory = append(m.requestHistory[1:], 100+rand.Float64()*200)
		m.lastUpdate = time.Time(msg)
		return m, tickCmd()

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case statusMsg:
		m.statusMsg = msg.message
		m.isError = msg.isError
		m.showStatus = true
		m.loading = false
		return m, tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
			return clearStatusMsg{}
		})

	case clearStatusMsg:
		m.showStatus = false

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.Back):
			if m.menuMode {
				m.menuMode = false
				m.cursor = 0
			}

		case key.Matches(msg, keys.Tab), key.Matches(msg, keys.Right):
			if !m.menuMode {
				m.activeTab = (m.activeTab + 1) % 4
			}

		case key.Matches(msg, keys.Left):
			if !m.menuMode {
				m.activeTab = (m.activeTab + 3) % 4
			}

		case key.Matches(msg, keys.Refresh):
			if !m.menuMode && m.activeTab == sitesTab {
				cmds = append(cmds, m.refreshSites())
			}

		case key.Matches(msg, keys.Enter):
			if !m.menuMode && m.activeTab == sitesTab {
				m.selected = m.table.Cursor()
				m.cursor = 0
				m.menuMode = true
			} else if m.menuMode {
				cmds = append(cmds, m.executeAction())
			}

		case key.Matches(msg, keys.Up):
			if m.menuMode {
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.activeTab == sitesTab {
				m.table, cmd = m.table.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.activeTab == logsTab {
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}

		case key.Matches(msg, keys.Down):
			if m.menuMode {
				if m.cursor < 3 {
					m.cursor++
				}
			} else if m.activeTab == sitesTab {
				m.table, cmd = m.table.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.activeTab == logsTab {
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	if !m.menuMode && m.activeTab == sitesTab {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

type clearStatusMsg struct{}

func (m model) refreshSites() tea.Cmd {
	return func() tea.Msg {
		sites, err := listNginxSites()
		if err != nil {
			return statusMsg{message: fmt.Sprintf("Error: %v", err), isError: true}
		}

		rows := []table.Row{}
		for _, s := range sites {
			status := "DISABLED"
			if s.Enabled {
				status = "ENABLED"
			}
			ssl := "NO"
			if s.SSL {
				ssl = "YES"
			}
			rows = append(rows, table.Row{s.Name, status, s.Port, ssl, s.Uptime})
		}

		return statusMsg{message: "Sites refreshed successfully", isError: false}
	}
}

func (m model) executeAction() tea.Cmd {
	return func() tea.Msg {
		if m.selected < 0 {
			return statusMsg{message: "No site selected", isError: true}
		}

		site := m.sites[m.selected]
		var message string
		var err error

		switch m.cursor {
		case 0: // Enable
			link := "/etc/nginx/sites-enabled/" + site.Name
			cmd := exec.Command("sudo", "ln", "-s", "/etc/nginx/sites-available/"+site.Name, link)
			err = cmd.Run()
			if err != nil {
				message = fmt.Sprintf("Failed to enable %s", site.Name)
			} else {
				message = fmt.Sprintf("Successfully enabled %s", site.Name)
			}

		case 1: // Disable
			link := "/etc/nginx/sites-enabled/" + site.Name
			cmd := exec.Command("sudo", "rm", link)
			err = cmd.Run()
			if err != nil {
				message = fmt.Sprintf("Failed to disable %s", site.Name)
			} else {
				message = fmt.Sprintf("Successfully disabled %s", site.Name)
			}

		case 2: // Test config
			cmd := exec.Command("sudo", "nginx", "-t")
			out, cmdErr := cmd.CombinedOutput()
			if cmdErr != nil {
				message = fmt.Sprintf("Config test failed: %s", string(out))
				err = cmdErr
			} else {
				message = "Configuration test passed"
			}

		case 3: // Reload
			cmd := exec.Command("sudo", "nginx", "-s", "reload")
			err = cmd.Run()
			if err != nil {
				message = "Failed to reload Nginx"
			} else {
				message = "Nginx reloaded successfully"
			}
		}

		return statusMsg{message: message, isError: err != nil}
	}
}

func (m model) View() string {
	if m.err != nil {
		return errorMsgStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
	}

	if m.quitting {
		return lipgloss.NewStyle().
			Foreground(accentPrimary).
			Bold(true).
			Render("\nGoodbye!\n\n")
	}

	var content strings.Builder

	// Header
	header := headerStyle.Render("NGINX TERMINAL MANAGER")
	content.WriteString(header + "\n")

	// Tabs
	tabs := m.renderTabs()
	content.WriteString(tabs + "\n\n")

	// Status message
	if m.showStatus {
		var statusStyle lipgloss.Style
		if m.isError {
			statusStyle = errorMsgStyle
		} else {
			statusStyle = successMsgStyle
		}
		statusBar := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderSubtle).
			Padding(0, 1).
			Width(m.width - 8)
		content.WriteString(statusBar.Render(statusStyle.Render(m.statusMsg)) + "\n\n")
	}

	// Content based on active tab
	switch m.activeTab {
	case sitesTab:
		if m.menuMode {
			content.WriteString(m.renderActionMenu())
		} else {
			content.WriteString(m.renderSitesTable())
		}
	case logsTab:
		content.WriteString(m.renderLogsView())
	case statsTab:
		content.WriteString(m.renderStatsView())
	case metricsTab:
		content.WriteString(m.renderMetricsView())
	}

	// Help
	content.WriteString("\n" + m.renderHelp())

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content.String())
}

func (m model) renderTabs() string {
	tabs := []string{"SITES", "LOGS", "STATS", "METRICS"}
	var renderedTabs []string

	for i, tab := range tabs {
		if tabType(i) == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(tab))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m model) renderSitesTable() string {
	return m.table.View()
}

func (m model) renderActionMenu() string {
	if m.selected < 0 || m.selected >= len(m.sites) {
		return "No site selected"
	}

	site := m.sites[m.selected]
	var b strings.Builder

	// Site info panel
	infoLines := []string{
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Site:"), valueStyle.Render(site.Name)),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Status:"), m.renderStatusBadge(site.Enabled)),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Port:"), valueStyle.Render(site.Port)),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("SSL:"), valueStyle.Render(fmt.Sprintf("%v", site.SSL))),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Uptime:"), valueStyle.Render(site.Uptime)),
	}

	infoContent := lipgloss.JoinVertical(lipgloss.Left, infoLines...)
	b.WriteString(panelStyle.Render(infoContent) + "\n\n")

	// Actions
	actions := []string{
		"Enable Site",
		"Disable Site",
		"Test Configuration",
		"Reload Nginx",
	}

	b.WriteString(panelTitleStyle.Render("ACTIONS") + "\n\n")

	for i, action := range actions {
		indicator := "  "
		if m.cursor == i {
			indicator = "▶ "
		}

		if m.cursor == i {
			b.WriteString(selectedActionStyle.Render(indicator+action) + "\n")
		} else {
			b.WriteString(actionItemStyle.Render(indicator+action) + "\n")
		}
	}

	return b.String()
}

func (m model) renderStatusBadge(enabled bool) string {
	if enabled {
		return enabledBadge.Render(" ENABLED ")
	}
	return disabledBadge.Render(" DISABLED ")
}

func (m model) renderLogsView() string {
	logs := []string{
		"[2024-11-04 10:23:45] 192.168.1.100 - GET /api/users - 200 OK - 45ms",
		"[2024-11-04 10:23:46] 192.168.1.101 - POST /api/auth - 201 Created - 123ms",
		"[2024-11-04 10:23:47] 192.168.1.102 - GET /static/main.css - 304 Not Modified - 12ms",
		"[2024-11-04 10:23:48] 192.168.1.103 - GET /api/posts - 200 OK - 67ms",
		"[2024-11-04 10:23:49] 192.168.1.104 - DELETE /api/posts/123 - 204 No Content - 89ms",
		"[2024-11-04 10:23:50] 192.168.1.105 - GET /api/comments - 200 OK - 34ms",
		"[2024-11-04 10:23:51] 192.168.1.106 - PUT /api/users/456 - 200 OK - 156ms",
		"[2024-11-04 10:23:52] 192.168.1.107 - POST /api/upload - 413 Payload Too Large - 23ms",
		"[2024-11-04 10:23:53] 192.168.1.108 - GET /health - 200 OK - 5ms",
		"[2024-11-04 10:23:54] 192.168.1.109 - GET /api/analytics - 200 OK - 234ms",
		"[2024-11-04 10:23:55] 192.168.1.110 - POST /api/webhook - 200 OK - 78ms",
		"[2024-11-04 10:23:56] 192.168.1.111 - GET /favicon.ico - 304 Not Modified - 3ms",
		"[2024-11-04 10:23:57] 192.168.1.112 - GET /api/notifications - 200 OK - 45ms",
		"[2024-11-04 10:23:58] 192.168.1.113 - DELETE /api/cache - 200 OK - 12ms",
		"[2024-11-04 10:23:59] 192.168.1.114 - GET /metrics - 200 OK - 8ms",
	}

	// Color code logs based on status
	var coloredLogs []string
	for _, log := range logs {
		styledLog := log
		if strings.Contains(log, "200") || strings.Contains(log, "201") || strings.Contains(log, "204") {
			styledLog = lipgloss.NewStyle().Foreground(accentSuccess).Render(log)
		} else if strings.Contains(log, "304") {
			styledLog = lipgloss.NewStyle().Foreground(textSecondary).Render(log)
		} else if strings.Contains(log, "413") {
			styledLog = lipgloss.NewStyle().Foreground(accentWarning).Render(log)
		} else if strings.Contains(log, "500") || strings.Contains(log, "404") {
			styledLog = lipgloss.NewStyle().Foreground(accentDanger).Render(log)
		}
		coloredLogs = append(coloredLogs, styledLog)
	}

	logContent := lipgloss.JoinVertical(lipgloss.Left, coloredLogs...)

	return panelStyle.
		Width(m.width - 8).
		Height(m.height - 12).
		Render(panelTitleStyle.Render("ACCESS LOGS") + "\n" + logContent)
}

func (m model) renderStatsView() string {
	var b strings.Builder

	// Metrics cards row
	cards := []string{
		m.renderMetricCard("TOTAL SITES", fmt.Sprintf("%d", len(m.sites)), ""),
		m.renderMetricCard("ACTIVE", "5", "+2 this week"),
		m.renderMetricCard("REQUESTS/MIN", "1,234", "+15%"),
		m.renderMetricCard("UPTIME", "99.9%", "45d 12h"),
	}

	metricsRow := lipgloss.JoinHorizontal(lipgloss.Top, cards...)
	b.WriteString(metricsRow + "\n\n")

	// Bar chart for site distribution
	b.WriteString(m.renderSiteDistributionChart())

	return b.String()
}

func (m model) renderMetricCard(label, value, subtext string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		metricLabelStyle.Render(label),
		"",
		lipgloss.NewStyle().Foreground(accentPrimary).Bold(true).Render(value),
		"",
		lipgloss.NewStyle().Foreground(textMuted).Render(subtext),
	)

	return metricCardStyle.Render(content)
}

func (m model) renderSiteDistributionChart() string {
	var b strings.Builder

	b.WriteString(panelTitleStyle.Render("SITE STATUS DISTRIBUTION") + "\n\n")

	// Count enabled/disabled sites
	enabled := 0
	disabled := 0
	for _, site := range m.sites {
		if site.Enabled {
			enabled++
		} else {
			disabled++
		}
	}

	// Create barchart data
	enabledStyle := lipgloss.NewStyle().Foreground(accentSuccess)
	disabledStyle := lipgloss.NewStyle().Foreground(accentDanger)

	data := []barchart.BarData{
		{
			Label: "Enabled",
			Values: []barchart.BarValue{
				{Name: "", Value: float64(enabled), Style: enabledStyle},
			},
		},
		{
			Label: "Disabled",
			Values: []barchart.BarValue{
				{Name: "", Value: float64(disabled), Style: disabledStyle},
			},
		},
	}

	// Create and configure bar chart
	bc := barchart.New(50, 8)
	bc.PushAll(data)
	bc.Draw()

	chartContent := bc.View()
	chartPanel := panelStyle.Width(50).Render(chartContent)
	b.WriteString(chartPanel)

	return b.String()
}

func (m model) renderMetricsView() string {
	var b strings.Builder

	b.WriteString(panelTitleStyle.Render("REAL-TIME METRICS") + "\n\n")

	// CPU Usage Chart
	b.WriteString(m.renderLineChart("CPU USAGE (%)", m.cpuHistory, accentPrimary))
	b.WriteString("\n\n")

	// Memory Usage Chart
	b.WriteString(m.renderLineChart("MEMORY USAGE (%)", m.memHistory, accentSuccess))
	b.WriteString("\n\n")

	// Network Traffic Chart
	b.WriteString(m.renderLineChart("NETWORK TRAFFIC (MB/s)", m.netHistory, accentWarning))
	b.WriteString("\n\n")

	// Request Rate Chart
	b.WriteString(m.renderLineChart("REQUEST RATE (req/s)", m.requestHistory, accentSecondary))

	return b.String()
}

func (m model) renderLineChart(title string, data []float64, color lipgloss.Color) string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(title) + "\n")

	// Compute stats
	currentValue := data[len(data)-1]
	avgValue := 0.0
	maxValue := 0.0
	minValue := math.MaxFloat64

	for _, v := range data {
		avgValue += v
		if v > maxValue {
			maxValue = v
		}
		if v < minValue {
			minValue = v
		}
	}
	avgValue /= float64(len(data))

	// Stats panel
	statsLines := []string{
		fmt.Sprintf("Current: %.1f", currentValue),
		fmt.Sprintf("Average: %.1f", avgValue),
		fmt.Sprintf("Peak: %.1f", maxValue),
		fmt.Sprintf("Min: %.1f", minValue),
	}

	// Create streamline chart
	slc := streamlinechart.New(80, 5)
	for _, v := range data {
		slc.Push(v)
	}
	slc.Draw()

	// Render the chart with stats
	chartView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		panelStyle.Width(85).Render(slc.View()),
		panelStyle.Width(20).Render(lipgloss.JoinVertical(lipgloss.Left, statsLines...)),
	)

	b.WriteString(chartView)

	return b.String()
}

// renderSparkline creates a sparkline using ntcharts
func renderSparkline(data []float64, width int) string {
	sl := sparkline.New(width, 1)
	sl.PushAll(data)
	sl.Draw()
	return sl.View()
}

func (m model) renderHelp() string {
	if m.menuMode {
		return helpStyle.Render("↑/↓ navigate • enter execute • esc back • q quit")
	}

	switch m.activeTab {
	case sitesTab:
		return helpStyle.Render("↑/↓ navigate • ←/→ switch tabs • enter select • r refresh • q quit")
	case logsTab:
		return helpStyle.Render("↑/↓ scroll • ←/→ switch tabs • q quit")
	case statsTab:
		return helpStyle.Render("←/→ switch tabs • q quit")
	case metricsTab:
		return helpStyle.Render("←/→ switch tabs • updates every second • q quit")
	default:
		return helpStyle.Render("←/→ switch tabs • q quit")
	}
}
