package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/aui/internal/agent"
	"github.com/yourusername/aui/internal/context"
)

func TestInitialApp(t *testing.T) {
	app := InitialApp()

	if app.ActiveTab != 0 {
		t.Errorf("InitialApp().ActiveTab = %v, want 0", app.ActiveTab)
	}

	if len(app.Tabs) != 3 {
		t.Errorf("InitialApp() should have 3 tabs, got %v", len(app.Tabs))
	}

	expectedTabs := []string{"Agents", "Contexts", "Files"}
	for i, tab := range expectedTabs {
		if app.Tabs[i] != tab {
			t.Errorf("Tab[%d] = %v, want %v", i, app.Tabs[i], tab)
		}
	}

	// Should start with Claude and Gemini agents
	if len(app.Agents) != 2 {
		t.Errorf("InitialApp() should have 2 agents, got %v", len(app.Agents))
	}

	// Check Claude agent
	if app.Agents[0].Name != "Claude" {
		t.Errorf("First agent name = %v, want Claude", app.Agents[0].Name)
	}
	if app.Agents[0].Model != "claude-3.5-sonnet" {
		t.Errorf("First agent model = %v, want claude-3.5-sonnet", app.Agents[0].Model)
	}

	// Check Gemini agent
	if app.Agents[1].Name != "Gemini" {
		t.Errorf("Second agent name = %v, want Gemini", app.Agents[1].Name)
	}
	if app.Agents[1].Model != "gemini-1.5-pro" {
		t.Errorf("Second agent model = %v, want gemini-1.5-pro", app.Agents[1].Model)
	}

	// Should have at least one example context
	if len(app.Contexts) < 1 {
		t.Errorf("InitialApp() should have at least 1 context, got %v", len(app.Contexts))
	}

	if !app.Ready {
		t.Error("InitialApp() should be ready")
	}
}

func TestAppInit(t *testing.T) {
	app := InitialApp()
	cmd := app.Init()

	if cmd != nil {
		t.Errorf("Init() should return nil, got %v", cmd)
	}
}

func TestAppTabSwitching(t *testing.T) {
	app := InitialApp()

	tests := []struct {
		name        string
		key         string
		expectedTab int
	}{
		{"switch to next tab", "tab", 1},
		{"switch with l key", "l", 1},
		{"switch to previous tab", "shift+tab", 2}, // wraps around to last tab
		{"switch with h key", "h", 2},              // wraps around
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ActiveTab = 0 // Reset to first tab

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{}}
			switch tt.key {
			case "tab":
				msg.Type = tea.KeyTab
			case "shift+tab":
				msg.Type = tea.KeyShiftTab
			case "l":
				msg.Runes = []rune{'l'}
			case "h":
				msg.Runes = []rune{'h'}
			}

			newModel, _ := app.Update(msg)
			updatedApp := newModel.(App)

			if updatedApp.ActiveTab != tt.expectedTab {
				t.Errorf("After pressing %s, ActiveTab = %v, want %v",
					tt.key, updatedApp.ActiveTab, tt.expectedTab)
			}
		})
	}
}

func TestAppQuit(t *testing.T) {
	app := InitialApp()

	tests := []string{"q", "ctrl+c"}

	for _, key := range tests {
		t.Run(key, func(t *testing.T) {
			app.Quitting = false // Reset

			var msg tea.KeyMsg
			if key == "q" {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			newModel, cmd := app.Update(msg)
			updatedApp := newModel.(App)

			if !updatedApp.Quitting {
				t.Errorf("After pressing %s, Quitting should be true", key)
			}

			// Check that Quit command is returned
			if cmd == nil {
				t.Errorf("After pressing %s, should return Quit command", key)
			}
		})
	}
}

func TestAppWindowResize(t *testing.T) {
	app := InitialApp()

	newWidth := 120
	newHeight := 40

	msg := tea.WindowSizeMsg{
		Width:  newWidth,
		Height: newHeight,
	}

	newModel, _ := app.Update(msg)
	updatedApp := newModel.(App)

	if updatedApp.Width != newWidth {
		t.Errorf("After resize, Width = %v, want %v", updatedApp.Width, newWidth)
	}

	if updatedApp.Height != newHeight {
		t.Errorf("After resize, Height = %v, want %v", updatedApp.Height, newHeight)
	}
}

func TestAppAddAgent(t *testing.T) {
	app := InitialApp()
	initialCount := len(app.Agents)

	app.AddAgent("GPT-4", "gpt-4-turbo")

	if len(app.Agents) != initialCount+1 {
		t.Errorf("After AddAgent, agent count = %v, want %v",
			len(app.Agents), initialCount+1)
	}

	lastAgent := app.Agents[len(app.Agents)-1]
	if lastAgent.Name != "GPT-4" {
		t.Errorf("New agent name = %v, want GPT-4", lastAgent.Name)
	}

	if lastAgent.Model != "gpt-4-turbo" {
		t.Errorf("New agent model = %v, want gpt-4-turbo", lastAgent.Model)
	}
}

func TestAppAddContext(t *testing.T) {
	app := InitialApp()
	initialCount := len(app.Contexts)

	app.AddContext("test-context", "Test context for unit tests")

	if len(app.Contexts) != initialCount+1 {
		t.Errorf("After AddContext, context count = %v, want %v",
			len(app.Contexts), initialCount+1)
	}

	lastContext := app.Contexts[len(app.Contexts)-1]
	if lastContext.Name != "test-context" {
		t.Errorf("New context name = %v, want test-context", lastContext.Name)
	}

	if lastContext.Description != "Test context for unit tests" {
		t.Errorf("New context description = %v, want 'Test context for unit tests'",
			lastContext.Description)
	}
}

func TestAppView(t *testing.T) {
	app := InitialApp()

	// Test normal view
	view := app.View()
	if view == "" {
		t.Error("View() should not return empty string")
	}

	// Test that view contains expected elements
	expectedElements := []string{
		"aui",               // Title
		"Agents",            // Tab name
		"Contexts",          // Tab name
		"Files",             // Tab name
		"Claude",            // Default agent
		"Gemini",            // Default agent
		"claude-3.5-sonnet", // Model name
		"gemini-1.5-pro",    // Model name
		"tab/l: next",       // Help text
		"q: quit",           // Help text
	}

	for _, expected := range expectedElements {
		if !strings.Contains(view, expected) {
			t.Errorf("View() should contain '%s'", expected)
		}
	}

	// Test quitting view
	app.Quitting = true
	view = app.View()
	if view != "Goodbye!\n" {
		t.Errorf("View() when quitting = %v, want 'Goodbye!\\n'", view)
	}
}

func TestAppViewActiveTabIndicator(t *testing.T) {
	app := InitialApp()

	// Test that active tab is indicated differently
	view := app.View()

	// For Agents tab (index 0, active by default)
	if !strings.Contains(view, "[Agents]") {
		t.Error("Active tab 'Agents' should be marked with brackets")
	}

	// Switch to Contexts tab
	app.ActiveTab = 1
	view = app.View()

	if !strings.Contains(view, "[Contexts]") {
		t.Error("Active tab 'Contexts' should be marked with brackets")
	}

	// Agents should no longer have brackets
	if strings.Contains(view, "[Agents]") {
		t.Error("Inactive tab 'Agents' should not have brackets")
	}
}

func TestAppViewContentChangesWithTab(t *testing.T) {
	app := InitialApp()

	tests := []struct {
		tabIndex int
		expected string
	}{
		{0, "Claude"},       // Agents tab shows agents
		{1, "bug-fix-auth"}, // Contexts tab shows contexts
		{2, "File browser"}, // Files tab shows placeholder
	}

	for _, tt := range tests {
		app.ActiveTab = tt.tabIndex
		view := app.View()

		if !strings.Contains(view, tt.expected) {
			t.Errorf("Tab %d view should contain '%s'", tt.tabIndex, tt.expected)
		}
	}
}

func TestAppViewEmptyStates(t *testing.T) {
	app := InitialApp()

	// Clear agents and contexts to test empty states
	app.Agents = []*agent.Agent{}
	app.Contexts = []*context.Context{}

	// Test empty agents view
	app.ActiveTab = 0
	view := app.View()
	if !strings.Contains(view, "No agents") || !strings.Contains(view, "Press 'a'") {
		t.Error("Empty agents view should show helpful message")
	}

	// Test empty contexts view
	app.ActiveTab = 1
	view = app.View()
	if !strings.Contains(view, "No contexts") || !strings.Contains(view, "Press 'c'") {
		t.Error("Empty contexts view should show helpful message")
	}
}
