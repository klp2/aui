package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/aui/internal/agent"
	"github.com/yourusername/aui/internal/context"
)

// App represents the main TUI application state
type App struct {
	ActiveTab int
	Tabs      []string
	Agents    []*agent.Agent
	Contexts  []*context.Context
	Width     int
	Height    int
	Ready     bool
	Quitting  bool
}

// InitialApp creates the initial application state
func InitialApp() App {
	// Create default agents
	claude := agent.NewAgent("Claude", "claude-3.5-sonnet")
	gemini := agent.NewAgent("Gemini", "gemini-1.5-pro")

	// Create example context
	exampleCtx := context.NewContext("bug-fix-auth", "Authentication bug context")

	return App{
		ActiveTab: 0,
		Tabs:      []string{"Agents", "Contexts", "Files"},
		Agents:    []*agent.Agent{claude, gemini},
		Contexts:  []*context.Context{exampleCtx},
		Ready:     true,
		Quitting:  false,
	}
}

// Init initializes the Bubble Tea application
func (a App) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the application state
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.Width = msg.Width
		a.Height = msg.Height
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			a.Quitting = true
			return a, tea.Quit

		case "tab", "l":
			a.ActiveTab = (a.ActiveTab + 1) % len(a.Tabs)
			return a, nil

		case "shift+tab", "h":
			a.ActiveTab = (a.ActiveTab - 1 + len(a.Tabs)) % len(a.Tabs)
			return a, nil
		}
	}

	return a, nil
}

// View renders the application UI
func (a App) View() string {
	if a.Quitting {
		return "Goodbye!\n"
	}

	// Simple view for now - will be enhanced with lipgloss styling
	view := "aui - AI Agent Manager\n\n"

	// Render tabs
	for i, tab := range a.Tabs {
		if i == a.ActiveTab {
			view += "[" + tab + "] "
		} else {
			view += " " + tab + "  "
		}
	}
	view += "\n\n"

	// Render content based on active tab
	switch a.ActiveTab {
	case 0: // Agents
		view += "Agents:\n"
		if len(a.Agents) == 0 {
			view += "  No agents configured. Press 'a' to add an agent.\n"
		} else {
			for _, ag := range a.Agents {
				view += fmt.Sprintf("  • %s (%s) - %s\n", ag.Name, ag.Model, ag.Status)
			}
		}

	case 1: // Contexts
		view += "Contexts:\n"
		if len(a.Contexts) == 0 {
			view += "  No contexts saved. Press 'c' to create a context.\n"
		} else {
			for _, ctx := range a.Contexts {
				view += fmt.Sprintf("  • %s - %s\n", ctx.Name, ctx.Description)
			}
		}

	case 2: // Files
		view += "Files:\n  (File browser coming soon)\n"
	}

	view += "\n[tab/l: next tab] [shift+tab/h: prev tab] [q: quit]"

	return view
}

// AddAgent adds a new agent to the application
func (a *App) AddAgent(name, model string) {
	newAgent := agent.NewAgent(name, model)
	a.Agents = append(a.Agents, newAgent)
}

// AddContext adds a new context to the application
func (a *App) AddContext(name, description string) {
	newContext := context.NewContext(name, description)
	a.Contexts = append(a.Contexts, newContext)
}
