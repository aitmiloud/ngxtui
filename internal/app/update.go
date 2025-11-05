package app

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aitmiloud/ngxtui/internal/model"
	"github.com/aitmiloud/ngxtui/internal/nginx"
	"github.com/aitmiloud/ngxtui/internal/ui"
)

// Update handles all state updates
func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Viewport.Width = msg.Width - 4
		m.Viewport.Height = msg.Height - 10

	case tea.KeyMsg:
		// Global keys
		if key.Matches(msg, model.Keys.Quit) {
			m.Quitting = true
			return m, tea.Quit
		}

		// Handle menu mode
		if m.MenuMode {
			return handleMenuMode(m, msg)
		}

		// Handle tab navigation
		if key.Matches(msg, model.Keys.Left) {
			if m.ActiveTab > 0 {
				m.ActiveTab--
			}
		} else if key.Matches(msg, model.Keys.Right) {
			if m.ActiveTab < model.MetricsTab {
				m.ActiveTab++
			}
		}

		// Handle refresh
		if key.Matches(msg, model.Keys.Refresh) {
			return m, refreshSites(&m)
		}

		// Tab-specific handling
		switch m.ActiveTab {
		case model.SitesTab:
			return handleSitesTab(m, msg)
		case model.LogsTab:
			m.Viewport, cmd = m.Viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case model.TickMsg:
		// Update REAL metrics when on Metrics tab
		if m.ActiveTab == model.MetricsTab {
			nginxService := nginx.New()
			metrics, err := nginxService.GetMetrics()

			if err == nil {
				// Calculate network rate (MB/s) from change in total bytes
				var networkRate float64
				if m.LastNetworkIn > 0 && m.LastNetworkOut > 0 {
					// Calculate change since last measurement
					deltaIn := metrics.NetworkIn - m.LastNetworkIn
					deltaOut := metrics.NetworkOut - m.LastNetworkOut
					networkRate = (deltaIn + deltaOut) / 1024 // Convert to MB/s (measured per second)
				}

				// Store current values for next calculation
				m.LastNetworkIn = metrics.NetworkIn
				m.LastNetworkOut = metrics.NetworkOut

				// Add new data point (will grow from 0 to 50 points)
				if len(m.CPUHistory) < 50 {
					// Still filling up - just append
					m.CPUHistory = append(m.CPUHistory, metrics.CPU)
					m.MemHistory = append(m.MemHistory, metrics.Memory)
					m.NetHistory = append(m.NetHistory, networkRate)
					m.RequestHistory = append(m.RequestHistory, metrics.RequestRate)
				} else {
					// Full - shift and add new data
					m.CPUHistory = append(m.CPUHistory[1:], metrics.CPU)
					m.MemHistory = append(m.MemHistory[1:], metrics.Memory)
					m.NetHistory = append(m.NetHistory[1:], networkRate)
					m.RequestHistory = append(m.RequestHistory[1:], metrics.RequestRate)
				}
			}
		}
		m.LastUpdate = time.Now()
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return model.TickMsg(t)
		}))

	case model.StatusMsg:
		m.StatusMsg = msg.Message
		m.IsError = msg.IsError
		m.ShowStatus = true
		cmds = append(cmds, clearStatusAfter(2*time.Second))

	case model.ClearStatusMsg:
		m.ShowStatus = false
		m.StatusMsg = ""

	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleSitesTab handles key events in the sites tab
func handleSitesTab(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	if key.Matches(msg, model.Keys.Enter) {
		if m.Cursor < len(m.Sites) {
			// Check if Docker NGINX
			if nginx.IsDockerAvailable() {
				if _, err := nginx.DetectDockerNginx(); err == nil {
					// Docker mode - show info message
					return m, func() tea.Msg {
						return model.StatusMsg{
							Message: "Docker NGINX detected - use 'docker exec' commands to manage. Press 'l' for logs.",
							IsError: false,
						}
					}
				}
			}
			// Native mode - show menu
			m.Selected = m.Cursor
			m.MenuMode = true
			m.Cursor = 0
		}
	} else if key.Matches(msg, model.Keys.Up) {
		if m.Selected > 0 {
			m.Selected--
		}
	} else if key.Matches(msg, model.Keys.Down) {
		if m.Selected < len(m.Sites)-1 {
			m.Selected++
		}
	}
	return m, nil
}

// handleMenuMode handles key events in menu mode
func handleMenuMode(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	if key.Matches(msg, model.Keys.Back) {
		m.MenuMode = false
		m.Selected = -1
		return m, nil
	}

	if key.Matches(msg, model.Keys.Up) {
		if m.Cursor > 0 {
			m.Cursor--
		}
	} else if key.Matches(msg, model.Keys.Down) {
		if m.Cursor < 5 { // 6 menu items (0-5)
			m.Cursor++
		}
	} else if key.Matches(msg, model.Keys.Enter) {
		return m, executeAction(&m)
	}

	return m, nil
}

// refreshSites refreshes the list of sites
func refreshSites(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		nginxService := nginx.New()
		sites, err := nginxService.ListSites()
		if err != nil {
			return model.StatusMsg{
				Message: "Failed to refresh sites: " + err.Error(),
				IsError: true,
			}
		}

		// Update model sites
		m.Sites = sites

		// Recreate table with new data
		m.Table = ui.CreateSitesTable(sites, 100, 15)

		return model.StatusMsg{
			Message: "Sites refreshed successfully",
			IsError: false,
		}
	}
}

// executeAction executes the selected action
func executeAction(m *model.Model) tea.Cmd {
	if m.Selected < 0 || m.Selected >= len(m.Sites) {
		return nil
	}

	site := m.Sites[m.Selected]
	nginxService := nginx.New()

	return func() tea.Msg {
		var err error
		var message string

		switch m.Cursor {
		case 0: // Enable Site
			err = nginxService.EnableSite(site.Name)
			message = "Site enabled successfully"
		case 1: // Disable Site
			err = nginxService.DisableSite(site.Name)
			message = "Site disabled successfully"
		case 2: // Test Configuration
			err = nginxService.TestConfig()
			message = "Configuration test passed"
		case 3: // Reload NGINX
			err = nginxService.Reload()
			message = "NGINX reloaded successfully"
		case 4: // View Logs
			// Switch to logs tab
			m.ActiveTab = model.LogsTab
			m.MenuMode = false
			m.Selected = -1
			return model.StatusMsg{
				Message: "Switched to logs view",
				IsError: false,
			}
		case 5: // Back
			m.MenuMode = false
			m.Selected = -1
			return model.StatusMsg{
				Message: "",
				IsError: false,
			}
		}

		if err != nil {
			return model.StatusMsg{
				Message: err.Error(),
				IsError: true,
			}
		}

		// Refresh sites after action
		sites, _ := nginxService.ListSites()
		m.Sites = sites

		// Recreate table with new data
		m.Table = ui.CreateSitesTable(sites, 100, 15)

		m.MenuMode = false
		m.Selected = -1

		return model.StatusMsg{
			Message: message,
			IsError: false,
		}
	}
}

// clearStatusAfter returns a command that clears the status after a duration
func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return model.ClearStatusMsg{}
	})
}
