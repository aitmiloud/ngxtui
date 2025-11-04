package model

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the Bubble Tea program
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		tickCmd(),
	)
}

// tickCmd returns a command that sends tick messages
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
