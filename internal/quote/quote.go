package quote

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Quote represents a motivational quote
type Quote struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Age returns how long ago the quote was created
func (q *Quote) Age() time.Duration {
	return time.Since(q.CreatedAt)
}

// DB handles all database operations for quotes
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
	CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		author TEXT DEFAULT '',
		category TEXT DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	-- Index for common queries
	CREATE INDEX IF NOT EXISTS idx_quotes_author ON quotes(author);
	CREATE INDEX IF NOT EXISTS idx_quotes_category ON quotes(category);
	CREATE INDEX IF NOT EXISTS idx_quotes_created_at ON quotes(created_at);
	`

	_, err := db.Exec(schema)
	return err
}

// Create creates a new quote
func (db *DB) Create(text, author, category string) (*Quote, error) {
	result, err := db.conn.Exec(`
		INSERT INTO quotes (text, author, category)
		VALUES (?, ?, ?)
	`, text, author, category)

	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get quote ID: %w", err)
	}

	return db.Get(int(id))
}

// Get retrieves a quote by ID
func (db *DB) Get(id int) (*Quote, error) {
	row := db.conn.QueryRow(`
		SELECT id, text, author, category, created_at, updated_at
		FROM quotes WHERE id = ?
	`, id)

	quote := &Quote{}
	err := row.Scan(
		&quote.ID, &quote.Text, &quote.Author, &quote.Category,
		&quote.CreatedAt, &quote.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	return quote, nil
}

// List retrieves all quotes with optional filtering
func (db *DB) List(author, category *string) ([]*Quote, error) {
	query := `
		SELECT id, text, author, category, created_at, updated_at
		FROM quotes WHERE 1=1
	`
	args := []interface{}{}

	if author != nil && *author != "" {
		query += " AND author = ?"
		args = append(args, *author)
	}

	if category != nil && *category != "" {
		query += " AND category = ?"
		args = append(args, *category)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*Quote
	for rows.Next() {
		quote := &Quote{}
		err := rows.Scan(
			&quote.ID, &quote.Text, &quote.Author, &quote.Category,
			&quote.CreatedAt, &quote.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}

		quotes = append(quotes, quote)
	}

	return quotes, rows.Err()
}

// GetRandom retrieves a random quote
func (db *DB) GetRandom() (*Quote, error) {
	row := db.conn.QueryRow(`
		SELECT id, text, author, category, created_at, updated_at
		FROM quotes ORDER BY RANDOM() LIMIT 1
	`)

	quote := &Quote{}
	err := row.Scan(
		&quote.ID, &quote.Text, &quote.Author, &quote.Category,
		&quote.CreatedAt, &quote.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return quote, nil
}

// Update updates a quote
func (db *DB) Update(id int, text, author, category string) error {
	_, err := db.conn.Exec(`
		UPDATE quotes 
		SET text = ?, author = ?, category = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, text, author, category, id)

	return err
}

// Delete deletes a quote by ID
func (db *DB) Delete(id int) error {
	_, err := db.conn.Exec("DELETE FROM quotes WHERE id = ?", id)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
