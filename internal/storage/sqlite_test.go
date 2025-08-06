package storage

import (
	"path/filepath"
	"testing"

	"github.com/yourusername/aui/internal/agent"
	"github.com/yourusername/aui/internal/context"
)

func TestNewSQLiteStore(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	if store == nil {
		t.Fatal("Expected store, got nil")
	}
}

func TestSQLiteStoreInitializeSchema(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Verify tables exist by attempting to query them
	var count int
	err = store.db.QueryRow("SELECT COUNT(*) FROM agents").Scan(&count)
	if err != nil {
		t.Errorf("agents table not created: %v", err)
	}

	err = store.db.QueryRow("SELECT COUNT(*) FROM contexts").Scan(&count)
	if err != nil {
		t.Errorf("contexts table not created: %v", err)
	}

	err = store.db.QueryRow("SELECT COUNT(*) FROM files").Scan(&count)
	if err != nil {
		t.Errorf("files table not created: %v", err)
	}

	err = store.db.QueryRow("SELECT COUNT(*) FROM context_files").Scan(&count)
	if err != nil {
		t.Errorf("context_files table not created: %v", err)
	}
}

func TestSQLiteStoreAgentCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Create agent
	a := agent.NewAgent("test-agent", "Claude", "anthropic")

	// Save agent
	err = store.SaveAgent(a)
	if err != nil {
		t.Fatalf("Failed to save agent: %v", err)
	}

	// Get agent
	retrieved, err := store.GetAgent(a.ID)
	if err != nil {
		t.Fatalf("Failed to get agent: %v", err)
	}

	if retrieved.ID != a.ID {
		t.Errorf("Expected ID %s, got %s", a.ID, retrieved.ID)
	}
	if retrieved.Name != a.Name {
		t.Errorf("Expected name %s, got %s", a.Name, retrieved.Name)
	}
	if retrieved.Model != a.Model {
		t.Errorf("Expected model %s, got %s", a.Model, retrieved.Model)
	}

	// Update agent
	retrieved.AssignTask("Test task")
	err = store.SaveAgent(retrieved)
	if err != nil {
		t.Fatalf("Failed to update agent: %v", err)
	}

	// List agents
	agents, err := store.ListAgents()
	if err != nil {
		t.Fatalf("Failed to list agents: %v", err)
	}

	if len(agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agents))
	}

	// Delete agent
	err = store.DeleteAgent(a.ID)
	if err != nil {
		t.Fatalf("Failed to delete agent: %v", err)
	}

	// Verify deletion
	_, err = store.GetAgent(a.ID)
	if err == nil {
		t.Error("Expected error getting deleted agent")
	}
}

func TestSQLiteStoreContextCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Create context
	ctx := context.NewContext("test-context", "Test context description")

	// Add files to context
	file1 := context.NewFile("/path/to/file1.go", "file1.go")
	file2 := context.NewFile("/path/to/file2.go", "file2.go")
	ctx.AddFile(file1)
	ctx.AddFile(file2)

	// Save context
	err = store.SaveContext(ctx)
	if err != nil {
		t.Fatalf("Failed to save context: %v", err)
	}

	// Get context
	retrieved, err := store.GetContext(ctx.ID)
	if err != nil {
		t.Fatalf("Failed to get context: %v", err)
	}

	if retrieved.ID != ctx.ID {
		t.Errorf("Expected ID %s, got %s", ctx.ID, retrieved.ID)
	}
	if retrieved.Name != ctx.Name {
		t.Errorf("Expected name %s, got %s", ctx.Name, retrieved.Name)
	}
	if len(retrieved.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(retrieved.Files))
	}

	// List contexts
	contexts, err := store.ListContexts()
	if err != nil {
		t.Fatalf("Failed to list contexts: %v", err)
	}

	if len(contexts) != 1 {
		t.Errorf("Expected 1 context, got %d", len(contexts))
	}

	// Delete context
	err = store.DeleteContext(ctx.ID)
	if err != nil {
		t.Fatalf("Failed to delete context: %v", err)
	}

	// Verify deletion
	_, err = store.GetContext(ctx.ID)
	if err == nil {
		t.Error("Expected error getting deleted context")
	}
}

func TestSQLiteStoreTransaction(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Start transaction
	tx, err := store.BeginTx()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Create and save agent in transaction
	a := agent.NewAgent("tx-agent", "Claude", "anthropic")
	err = store.SaveAgentTx(tx, a)
	if err != nil {
		t.Fatalf("Failed to save agent in transaction: %v", err)
	}

	// Rollback transaction
	err = tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to rollback transaction: %v", err)
	}

	// Verify agent was not saved
	_, err = store.GetAgent(a.ID)
	if err == nil {
		t.Error("Expected error getting agent after rollback")
	}

	// Test commit
	tx, err = store.BeginTx()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	a2 := agent.NewAgent("tx-agent-2", "Claude", "anthropic")
	err = store.SaveAgentTx(tx, a2)
	if err != nil {
		t.Fatalf("Failed to save agent in transaction: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// Verify agent was saved
	retrieved, err := store.GetAgent(a2.ID)
	if err != nil {
		t.Fatalf("Failed to get agent after commit: %v", err)
	}

	if retrieved.ID != a2.ID {
		t.Error("Agent not properly saved after commit")
	}
}

func TestSQLiteStoreMigration(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Close and reopen to test migration handling
	store.Close()

	store2, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to reopen store: %v", err)
	}
	defer store2.Close()

	// Verify schema still exists
	var count int
	err = store2.db.QueryRow("SELECT COUNT(*) FROM agents").Scan(&count)
	if err != nil {
		t.Errorf("Schema not preserved after reopen: %v", err)
	}
}
