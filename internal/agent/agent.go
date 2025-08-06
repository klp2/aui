package agent

import (
	"crypto/rand"
	"encoding/hex"
)

// Status represents the current state of an agent
type Status string

const (
	StatusReady   Status = "ready"
	StatusWorking Status = "working"
	StatusError   Status = "error"
)

// Agent represents an AI agent that can perform tasks
type Agent struct {
	ID          string
	Name        string
	Model       string
	Provider    string
	Status      Status
	CurrentTask string
	LastError   string
}

// NewAgent creates a new agent with the given name, model, and provider
func NewAgent(name, model, provider string) *Agent {
	return &Agent{
		ID:          generateID(),
		Name:        name,
		Model:       model,
		Provider:    provider,
		Status:      StatusReady,
		CurrentTask: "",
	}
}

// AssignTask assigns a task to the agent and sets status to working
func (a *Agent) AssignTask(task string) {
	a.CurrentTask = task
	a.Status = StatusWorking
}

// CompleteTask marks the current task as complete and returns to ready status
func (a *Agent) CompleteTask() {
	a.CurrentTask = ""
	a.Status = StatusReady
	a.LastError = ""
}

// SetError sets the agent to error state with the given error message
func (a *Agent) SetError(err string) {
	a.Status = StatusError
	a.LastError = err
}

// generateID generates a random ID for an agent
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
