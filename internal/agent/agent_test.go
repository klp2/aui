package agent

import (
	"testing"
)

func TestNewAgent(t *testing.T) {
	tests := []struct {
		name      string
		agentName string
		model     string
		provider  string
		want      Agent
	}{
		{
			name:      "create Claude agent",
			agentName: "Claude",
			model:     "claude-3.5-sonnet",
			provider:  "anthropic",
			want: Agent{
				Name:        "Claude",
				Model:       "claude-3.5-sonnet",
				Provider:    "anthropic",
				Status:      StatusReady,
				CurrentTask: "",
			},
		},
		{
			name:      "create Gemini agent",
			agentName: "Gemini",
			model:     "gemini-1.5-pro",
			provider:  "google",
			want: Agent{
				Name:        "Gemini",
				Model:       "gemini-1.5-pro",
				Provider:    "google",
				Status:      StatusReady,
				CurrentTask: "",
			},
		},
		{
			name:      "create GPT-4 agent",
			agentName: "GPT-4",
			model:     "gpt-4",
			provider:  "openai",
			want: Agent{
				Name:        "GPT-4",
				Model:       "gpt-4",
				Provider:    "openai",
				Status:      StatusReady,
				CurrentTask: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAgent(tt.agentName, tt.model, tt.provider)

			if got.ID == "" {
				t.Error("NewAgent().ID should not be empty")
			}
			if got.Name != tt.want.Name {
				t.Errorf("NewAgent().Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Model != tt.want.Model {
				t.Errorf("NewAgent().Model = %v, want %v", got.Model, tt.want.Model)
			}
			if got.Provider != tt.want.Provider {
				t.Errorf("NewAgent().Provider = %v, want %v", got.Provider, tt.want.Provider)
			}
			if got.Status != tt.want.Status {
				t.Errorf("NewAgent().Status = %v, want %v", got.Status, tt.want.Status)
			}
			if got.CurrentTask != tt.want.CurrentTask {
				t.Errorf("NewAgent().CurrentTask = %v, want %v", got.CurrentTask, tt.want.CurrentTask)
			}
		})
	}
}

func TestAgentAssignTask(t *testing.T) {
	agent := NewAgent("Claude", "claude-3.5-sonnet", "anthropic")

	task := "Fix authentication bug"
	agent.AssignTask(task)

	if agent.CurrentTask != task {
		t.Errorf("After AssignTask(), CurrentTask = %v, want %v", agent.CurrentTask, task)
	}

	if agent.Status != StatusWorking {
		t.Errorf("After AssignTask(), Status = %v, want %v", agent.Status, StatusWorking)
	}
}

func TestAgentCompleteTask(t *testing.T) {
	agent := NewAgent("Gemini", "gemini-1.5-pro", "google")
	agent.AssignTask("Review code")

	agent.CompleteTask()

	if agent.CurrentTask != "" {
		t.Errorf("After CompleteTask(), CurrentTask = %v, want empty string", agent.CurrentTask)
	}

	if agent.Status != StatusReady {
		t.Errorf("After CompleteTask(), Status = %v, want %v", agent.Status, StatusReady)
	}

	if agent.LastError != "" {
		t.Errorf("After CompleteTask(), LastError = %v, want empty string", agent.LastError)
	}
}

func TestAgentSetError(t *testing.T) {
	agent := NewAgent("Claude", "claude-3.5-sonnet", "anthropic")

	errorMsg := "API rate limit exceeded"
	agent.SetError(errorMsg)

	if agent.Status != StatusError {
		t.Errorf("After SetError(), Status = %v, want %v", agent.Status, StatusError)
	}

	if agent.LastError != errorMsg {
		t.Errorf("After SetError(), LastError = %v, want %v", agent.LastError, errorMsg)
	}
}

func TestGenerateID(t *testing.T) {
	// Test that IDs are unique
	id1 := generateID()
	id2 := generateID()

	if id1 == id2 {
		t.Error("generateID() should produce unique IDs")
	}

	if len(id1) != 16 {
		t.Errorf("generateID() should produce 16-character hex strings, got %d", len(id1))
	}
}
