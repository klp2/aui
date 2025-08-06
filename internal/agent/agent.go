package agent

// Status represents the current state of an agent
type Status string

const (
	StatusReady   Status = "ready"
	StatusWorking Status = "working"
	StatusError   Status = "error"
)

// Agent represents an AI agent that can perform tasks
type Agent struct {
	Name      string
	Model     string
	Status    Status
	Task      string
	LastError string
}

// NewAgent creates a new agent with the given name and model
func NewAgent(name, model string) *Agent {
	return &Agent{
		Name:   name,
		Model:  model,
		Status: StatusReady,
		Task:   "",
	}
}

// AssignTask assigns a task to the agent and sets status to working
func (a *Agent) AssignTask(task string) {
	a.Task = task
	a.Status = StatusWorking
}

// CompleteTask marks the current task as complete and returns to ready status
func (a *Agent) CompleteTask() {
	a.Task = ""
	a.Status = StatusReady
	a.LastError = ""
}

// SetError sets the agent to error state with the given error message
func (a *Agent) SetError(err string) {
	a.Status = StatusError
	a.LastError = err
}