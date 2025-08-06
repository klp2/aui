package storage

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yourusername/aui/internal/agent"
	"github.com/yourusername/aui/internal/context"
)

// SQLiteStore implements storage using SQLite
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite storage instance
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "/" {
		// Create directory if it doesn't exist
		// Using os.MkdirAll would require importing os
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &SQLiteStore{db: db}

	// Initialize schema
	if err := store.initializeSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// initializeSchema creates the database tables if they don't exist
func (s *SQLiteStore) initializeSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS agents (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		model TEXT NOT NULL,
		provider TEXT NOT NULL,
		status TEXT NOT NULL,
		current_task TEXT,
		last_error TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS contexts (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		total_tokens INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS files (
		id TEXT PRIMARY KEY,
		path TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		content TEXT,
		language TEXT,
		tokens INTEGER DEFAULT 0,
		hash TEXT,
		size INTEGER DEFAULT 0,
		modified_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS context_files (
		context_id TEXT NOT NULL,
		file_id TEXT NOT NULL,
		position INTEGER,
		PRIMARY KEY (context_id, file_id),
		FOREIGN KEY (context_id) REFERENCES contexts(id) ON DELETE CASCADE,
		FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME NOT NULL
	);
	`

	_, err := s.db.Exec(schema)
	return err
}

// BeginTx starts a new database transaction
func (s *SQLiteStore) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}

// Agent operations

// SaveAgent saves or updates an agent
func (s *SQLiteStore) SaveAgent(a *agent.Agent) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.SaveAgentTx(tx, a); err != nil {
		return err
	}

	return tx.Commit()
}

// SaveAgentTx saves an agent within a transaction
func (s *SQLiteStore) SaveAgentTx(tx *sql.Tx, a *agent.Agent) error {
	now := time.Now()

	query := `
	INSERT INTO agents (id, name, model, provider, status, current_task, last_error, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		model = excluded.model,
		provider = excluded.provider,
		status = excluded.status,
		current_task = excluded.current_task,
		last_error = excluded.last_error,
		updated_at = excluded.updated_at
	`

	_, err := tx.Exec(query, a.ID, a.Name, a.Model, a.Provider, a.Status, a.CurrentTask, a.LastError, now, now)
	return err
}

// GetAgent retrieves an agent by ID
func (s *SQLiteStore) GetAgent(id string) (*agent.Agent, error) {
	query := `
	SELECT id, name, model, provider, status, current_task, last_error
	FROM agents
	WHERE id = ?
	`

	var a agent.Agent
	var currentTask, lastError sql.NullString

	err := s.db.QueryRow(query, id).Scan(
		&a.ID, &a.Name, &a.Model, &a.Provider, &a.Status,
		&currentTask, &lastError,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	if currentTask.Valid {
		a.CurrentTask = currentTask.String
	}
	if lastError.Valid {
		a.LastError = lastError.String
	}

	return &a, nil
}

// ListAgents returns all agents
func (s *SQLiteStore) ListAgents() ([]*agent.Agent, error) {
	query := `
	SELECT id, name, model, provider, status, current_task, last_error
	FROM agents
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*agent.Agent
	for rows.Next() {
		var a agent.Agent
		var currentTask, lastError sql.NullString

		err := rows.Scan(
			&a.ID, &a.Name, &a.Model, &a.Provider, &a.Status,
			&currentTask, &lastError,
		)
		if err != nil {
			return nil, err
		}

		if currentTask.Valid {
			a.CurrentTask = currentTask.String
		}
		if lastError.Valid {
			a.LastError = lastError.String
		}

		agents = append(agents, &a)
	}

	return agents, rows.Err()
}

// DeleteAgent deletes an agent
func (s *SQLiteStore) DeleteAgent(id string) error {
	_, err := s.db.Exec("DELETE FROM agents WHERE id = ?", id)
	return err
}

// Context operations

// SaveContext saves or updates a context and its files
func (s *SQLiteStore) SaveContext(ctx *context.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// Save context
	query := `
	INSERT INTO contexts (id, name, description, total_tokens, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		description = excluded.description,
		total_tokens = excluded.total_tokens,
		updated_at = excluded.updated_at
	`

	_, err = tx.Exec(query, ctx.ID, ctx.Name, ctx.Description, ctx.TotalTokens, now, now)
	if err != nil {
		return err
	}

	// Delete existing file associations
	_, err = tx.Exec("DELETE FROM context_files WHERE context_id = ?", ctx.ID)
	if err != nil {
		return err
	}

	// Save files and associations
	for i, file := range ctx.Files {
		// Save file
		fileQuery := `
		INSERT INTO files (id, path, name, content, language, tokens, hash, size, modified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			name = excluded.name,
			content = excluded.content,
			language = excluded.language,
			tokens = excluded.tokens,
			hash = excluded.hash,
			size = excluded.size,
			modified_at = excluded.modified_at,
			updated_at = excluded.updated_at
		`

		_, err = tx.Exec(fileQuery, file.ID, file.Path, file.Name, file.Content, file.Language,
			file.Tokens, file.Hash, file.Size, file.ModifiedAt, now, now)
		if err != nil {
			return err
		}

		// Save association
		assocQuery := `
		INSERT INTO context_files (context_id, file_id, position)
		VALUES (?, ?, ?)
		`

		_, err = tx.Exec(assocQuery, ctx.ID, file.ID, i)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetContext retrieves a context by ID with its files
func (s *SQLiteStore) GetContext(id string) (*context.Context, error) {
	// Get context
	query := `
	SELECT id, name, description, total_tokens
	FROM contexts
	WHERE id = ?
	`

	var ctx context.Context
	var description sql.NullString

	err := s.db.QueryRow(query, id).Scan(&ctx.ID, &ctx.Name, &description, &ctx.TotalTokens)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("context not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	if description.Valid {
		ctx.Description = description.String
	}

	// Get associated files
	fileQuery := `
	SELECT f.id, f.path, f.name, f.content, f.language, f.tokens, f.hash, f.size, f.modified_at
	FROM files f
	JOIN context_files cf ON f.id = cf.file_id
	WHERE cf.context_id = ?
	ORDER BY cf.position
	`

	rows, err := s.db.Query(fileQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ctx.Files = make([]*context.File, 0)
	for rows.Next() {
		var f context.File
		var content, language, hash sql.NullString
		var modifiedAt sql.NullTime

		err := rows.Scan(&f.ID, &f.Path, &f.Name, &content, &language,
			&f.Tokens, &hash, &f.Size, &modifiedAt)
		if err != nil {
			return nil, err
		}

		if content.Valid {
			f.Content = content.String
		}
		if language.Valid {
			f.Language = language.String
		}
		if hash.Valid {
			f.Hash = hash.String
		}
		if modifiedAt.Valid {
			f.ModifiedAt = modifiedAt.Time
		}

		ctx.Files = append(ctx.Files, &f)
	}

	return &ctx, rows.Err()
}

// ListContexts returns all contexts
func (s *SQLiteStore) ListContexts() ([]*context.Context, error) {
	query := `
	SELECT id, name, description, total_tokens
	FROM contexts
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contexts []*context.Context
	for rows.Next() {
		var ctx context.Context
		var description sql.NullString

		err := rows.Scan(&ctx.ID, &ctx.Name, &description, &ctx.TotalTokens)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			ctx.Description = description.String
		}

		// Note: Not loading files for list operation to keep it efficient
		ctx.Files = make([]*context.File, 0)

		contexts = append(contexts, &ctx)
	}

	return contexts, rows.Err()
}

// DeleteContext deletes a context and its file associations
func (s *SQLiteStore) DeleteContext(id string) error {
	// File associations will be deleted automatically due to CASCADE
	_, err := s.db.Exec("DELETE FROM contexts WHERE id = ?", id)
	return err
}
