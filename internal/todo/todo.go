package todo

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Priority represents task priority levels
type Priority int

const (
	Low Priority = iota + 1
	Medium
	High
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "LOW"
	case Medium:
		return "MEDIUM"
	case High:
		return "HIGH"
	default:
		return "UNKNOWN"
	}
}

// Status represents task completion status
type Status int

const (
	Open Status = iota + 1
	Done
)

func (s Status) String() string {
	switch s {
	case Open:
		return "OPEN"
	case Done:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}

// Todo represents a todo task
type Todo struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Age returns how long ago the todo was created
func (t *Todo) Age() time.Duration {
	return time.Since(t.CreatedAt)
}

// CompletedAge returns how long ago the todo was completed (if completed)
func (t *Todo) CompletedAge() *time.Duration {
	if t.CompletedAt == nil {
		return nil
	}
	age := time.Since(*t.CompletedAt)
	return &age
}

// DB handles all database operations
type DB struct {
	conn *sql.DB
}

// GetDatabasePath returns the path to the database file
func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	tadaDir := filepath.Join(homeDir, ".tada")
	if err := os.MkdirAll(tadaDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create .tada directory: %w", err)
	}

	return filepath.Join(tadaDir, "todos.db"), nil
}

// NewDB creates a new database connection
func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create table and indexes
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &DB{conn: db}, nil
}

// createTables creates the database schema
func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		priority INTEGER NOT NULL DEFAULT 2 CHECK(priority IN (1, 2, 3)),
		status INTEGER NOT NULL DEFAULT 1 CHECK(status IN (1, 2)),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME NULL
	);

	-- Index for common queries
	CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
	CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
	CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);

	-- Quotes table
	CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		author TEXT DEFAULT '',
		category TEXT DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	-- Index for quotes
	CREATE INDEX IF NOT EXISTS idx_quotes_author ON quotes(author);
	CREATE INDEX IF NOT EXISTS idx_quotes_category ON quotes(category);
	CREATE INDEX IF NOT EXISTS idx_quotes_created_at ON quotes(created_at);
	`

	_, err := db.Exec(schema)
	return err
}

// Create creates a new todo task
func (db *DB) Create(description string, priority Priority) (*Todo, error) {
	result, err := db.conn.Exec(`
		INSERT INTO todos (description, priority)
		VALUES (?, ?)
	`, description, int(priority))

	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get todo ID: %w", err)
	}

	return db.Get(int(id))
}

// Get retrieves a todo by ID
func (db *DB) Get(id int) (*Todo, error) {
	row := db.conn.QueryRow(`
		SELECT id, description, priority, status, created_at, updated_at, completed_at
		FROM todos WHERE id = ?
	`, id)

	todo := &Todo{}
	var completedAt sql.NullTime

	err := row.Scan(
		&todo.ID, &todo.Description, &todo.Priority, &todo.Status,
		&todo.CreatedAt, &todo.UpdatedAt, &completedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	if completedAt.Valid {
		todo.CompletedAt = &completedAt.Time
	}

	return todo, nil
}

// List retrieves todos with optional filtering
func (db *DB) List(status *Status, priority *Priority) ([]*Todo, error) {
	query := `
		SELECT id, description, priority, status, created_at, updated_at, completed_at
		FROM todos WHERE 1=1
	`
	args := []interface{}{}

	if status != nil {
		query += " AND status = ?"
		args = append(args, int(*status))
	}

	if priority != nil {
		query += " AND priority = ?"
		args = append(args, int(*priority))
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		var completedAt sql.NullTime

		err := rows.Scan(
			&todo.ID, &todo.Description, &todo.Priority, &todo.Status,
			&todo.CreatedAt, &todo.UpdatedAt, &completedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}

		if completedAt.Valid {
			todo.CompletedAt = &completedAt.Time
		}

		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

// UpdateStatus updates the status of a todo
func (db *DB) UpdateStatus(id int, status Status) error {
	var completedAt interface{}
	if status == Done {
		completedAt = time.Now()
	} else {
		completedAt = nil
	}

	_, err := db.conn.Exec(`
		UPDATE todos 
		SET status = ?, completed_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, int(status), completedAt, id)

	return err
}

// UpdatePriority updates the priority of a todo
func (db *DB) UpdatePriority(id int, priority Priority) error {
	_, err := db.conn.Exec(`
		UPDATE todos 
		SET priority = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, int(priority), id)

	return err
}

// UpdateDescription updates the description of a todo
func (db *DB) UpdateDescription(id int, description string) error {
	_, err := db.conn.Exec(`
		UPDATE todos 
		SET description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, description, id)

	return err
}

// Delete deletes a todo by ID
func (db *DB) Delete(id int) error {
	_, err := db.conn.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
