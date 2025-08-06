package agent

import (
	"testing"
)

func TestNewAgent(t *testing.T) {
	tests := []struct {
		name      string
		agentName string
		model     string
		want      Agent
	}{
		{
			name:      "create Claude agent",
			agentName: "Claude",
			model:     "claude-3.5-sonnet",
			want: Agent{
				Name:   "Claude",
				Model:  "claude-3.5-sonnet",
				Status: StatusReady,
				Task:   "",
			},
		},
		{
			name:      "create Gemini agent",
			agentName: "Gemini",
			model:     "gemini-1.5-pro",
			want: Agent{
				Name:   "Gemini",
				Model:  "gemini-1.5-pro",
				Status: StatusReady,
				Task:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAgent(tt.agentName, tt.model)

			if got.Name != tt.want.Name {
				t.Errorf("NewAgent().Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Model != tt.want.Model {
				t.Errorf("NewAgent().Model = %v, want %v", got.Model, tt.want.Model)
			}
			if got.Status != tt.want.Status {
				t.Errorf("NewAgent().Status = %v, want %v", got.Status, tt.want.Status)
			}
			if got.Task != tt.want.Task {
				t.Errorf("NewAgent().Task = %v, want %v", got.Task, tt.want.Task)
			}
		})
	}
}

func TestAgentAssignTask(t *testing.T) {
	agent := NewAgent("Claude", "claude-3.5-sonnet")

	task := "Fix authentication bug"
	agent.AssignTask(task)

	if agent.Task != task {
		t.Errorf("After AssignTask(), Task = %v, want %v", agent.Task, task)
	}

	if agent.Status != StatusWorking {
		t.Errorf("After AssignTask(), Status = %v, want %v", agent.Status, StatusWorking)
	}
}

func TestAgentCompleteTask(t *testing.T) {
	agent := NewAgent("Gemini", "gemini-1.5-pro")
	agent.AssignTask("Review code")

	agent.CompleteTask()

	if agent.Task != "" {
		t.Errorf("After CompleteTask(), Task = %v, want empty string", agent.Task)
	}

	if agent.Status != StatusReady {
		t.Errorf("After CompleteTask(), Status = %v, want %v", agent.Status, StatusReady)
	}
}

func TestAgentSetError(t *testing.T) {
	agent := NewAgent("Claude", "claude-3.5-sonnet")

	errorMsg := "API rate limit exceeded"
	agent.SetError(errorMsg)

	if agent.Status != StatusError {
		t.Errorf("After SetError(), Status = %v, want %v", agent.Status, StatusError)
	}

	if agent.LastError != errorMsg {
		t.Errorf("After SetError(), LastError = %v, want %v", agent.LastError, errorMsg)
	}
}
