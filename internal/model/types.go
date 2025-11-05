package model

import (
	"time"

	"github.com/76creates/stickers/table"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
)

// TabType represents the different tabs in the application
type TabType int

const (
	SitesTab TabType = iota
	LogsTab
	StatsTab
	MetricsTab
)

// Site represents an NGINX site configuration
type Site struct {
	Name    string
	Enabled bool
	Port    string
	SSL     bool
	Uptime  string
}

// StatusMsg represents a status message to display to the user
type StatusMsg struct {
	Message string
	IsError bool
}

// TickMsg is sent on each tick for animations and updates
type TickMsg time.Time

// ClearStatusMsg is sent to clear the status message
type ClearStatusMsg struct{}

// Model represents the application state
type Model struct {
	Sites          []Site
	Table          *table.Table
	Cursor         int
	Selected       int
	MenuMode       bool
	ActiveTab      TabType
	Quitting       bool
	Err            error
	StatusMsg      string
	ShowStatus     bool
	IsError        bool
	Spinner        spinner.Model
	Loading        bool
	Width          int
	Height         int
	Viewport       viewport.Model
	CPUHistory     []float64
	MemHistory     []float64
	NetHistory     []float64
	RequestHistory []float64
	Progress       progress.Model
	LastUpdate     time.Time
	MetricsHistory interface{} // Will store *nginx.MetricsHistory
	LastNetworkIn  float64     // Track last network in for rate calculation
	LastNetworkOut float64     // Track last network out for rate calculation
}

// KeyMap defines the keybindings for the application
type KeyMap struct {
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

// Keys is the default keymap
var Keys = KeyMap{
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
		key.WithHelp("←/h", "previous tab"),
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

// ShortHelp returns a short help text
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Enter, k.Back, k.Refresh, k.Quit}
}

// FullHelp returns the full help text
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Back, k.Refresh, k.Quit},
	}
}
