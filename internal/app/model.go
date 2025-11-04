package app

import (
	"github.com/aitmiloud/ngxtui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

// App wraps the model and implements tea.Model interface
type App struct {
	Model model.Model
}

// New creates a new App instance
func New() *App {
	return &App{
		Model: InitialModel(),
	}
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return a.Model.Init()
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := Update(a.Model, msg)
	a.Model = newModel
	return a, cmd
}

// View implements tea.Model
func (a *App) View() string {
	return View(a.Model)
}
