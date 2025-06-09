package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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

type Todo struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"` // Only set when status = Done
}

func (t *Todo) Age() time.Duration {
	return time.Since(t.CreatedAt)
}

func (t *Todo) CompletedAge() *time.Duration {
	if t.CompletedAt == nil {
		return nil
	}

	age := time.Since(*t.CompletedAt)
	return &age
}

// Database schema
const createTableSQL = `
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
`

type TodoDB struct {
	db *sql.DB
}

func NewTodoDB(dbPath string) (*TodoDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table and indexes
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return &TodoDB{db: db}, nil
}

func (tdb *TodoDB) CreateTodo(description string, priority Priority) (*Todo, error) {
	result, err := tdb.db.Exec(`
	  INSERT INTO todos (description, priority)
	  VALUES (?, ?)
	`, description, int(priority))

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return tdb.GetTodo(int(id))
}

func (tdb *TodoDB) GetTodo(id int) (*Todo, error) {
	row := tdb.db.QueryRow(`
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
		return nil, err
	}

	if completedAt.Valid {
		todo.CompletedAt = &completedAt.Time
	}

	return todo, nil
}

func (tdb *TodoDB) ListTodos(status *Status, priority *Priority) ([]*Todo, error) {
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

	rows, err := tdb.db.Query(query, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if completedAt.Valid {
			todo.CompletedAt = &completedAt.Time
		}

		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func (tdb *TodoDB) UpdateTodoStatus(id int, status Status) error {
	var completedAt interface{}
	if status == Done {
		completedAt = time.Now()
	} else {
		completedAt = nil
	}

	_, err := tdb.db.Exec(`
		UPDATE todos 
		SET status = ?, completed_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, int(status), completedAt, id)

	return err
}

func (tdb *TodoDB) UpdateTodoPriority(id int, priority Priority) error {
	_, err := tdb.db.Exec(`
		UPDATE todos 
		SET priority = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, int(priority), id)

	return err
}

func (tdb *TodoDB) UpdateTodoDescription(id int, description string) error {
	_, err := tdb.db.Exec(`
		UPDATE todos 
		SET description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, description, id)

	return err
}

func (tdb *TodoDB) DeleteTodo(id int) error {
	_, err := tdb.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

func (tdb *TodoDB) Close() error {
	return tdb.db.Close()
}
