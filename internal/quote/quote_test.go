package quote

import (
	"path/filepath"
	"testing"
	"time"
)

func TestQuote_Age(t *testing.T) {
	now := time.Now()
	quote := &Quote{
		CreatedAt: now.Add(-1 * time.Hour),
	}

	age := quote.Age()
	if age < 59*time.Minute || age > 61*time.Minute {
		t.Errorf("Quote.Age() = %v, want approximately 1 hour", age)
	}
}

func TestDB_Integration(t *testing.T) {
	// Create temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer db.Close()

	t.Run("Create quote", func(t *testing.T) {
		quote, err := db.Create("Test quote", "Test Author", "test")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if quote.ID == 0 {
			t.Error("Create() returned quote with ID 0")
		}
		if quote.Text != "Test quote" {
			t.Errorf("Create() text = %v, want %v", quote.Text, "Test quote")
		}
		if quote.Author != "Test Author" {
			t.Errorf("Create() author = %v, want %v", quote.Author, "Test Author")
		}
		if quote.Category != "test" {
			t.Errorf("Create() category = %v, want %v", quote.Category, "test")
		}
	})

	t.Run("Create quote with empty author and category", func(t *testing.T) {
		quote, err := db.Create("Simple quote", "", "")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if quote.Author != "" {
			t.Errorf("Create() author = %v, want empty string", quote.Author)
		}
		if quote.Category != "" {
			t.Errorf("Create() category = %v, want empty string", quote.Category)
		}
	})

	t.Run("Get quote", func(t *testing.T) {
		// Create a quote first
		created, err := db.Create("Get test", "Get Author", "get")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Get the quote
		retrieved, err := db.Get(created.ID)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if retrieved.ID != created.ID {
			t.Errorf("Get() ID = %v, want %v", retrieved.ID, created.ID)
		}
		if retrieved.Text != created.Text {
			t.Errorf("Get() text = %v, want %v", retrieved.Text, created.Text)
		}
		if retrieved.Author != created.Author {
			t.Errorf("Get() author = %v, want %v", retrieved.Author, created.Author)
		}
		if retrieved.Category != created.Category {
			t.Errorf("Get() category = %v, want %v", retrieved.Category, created.Category)
		}
	})

	t.Run("Get non-existent quote", func(t *testing.T) {
		_, err := db.Get(99999)
		if err == nil {
			t.Error("Get() with non-existent ID should return error")
		}
	})

	t.Run("List quotes", func(t *testing.T) {
		// Create multiple quotes
		_, err := db.Create("List test 1", "Author 1", "category1")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		_, err = db.Create("List test 2", "Author 2", "category2")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// List all quotes
		quotes, err := db.List(nil, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(quotes) < 2 {
			t.Errorf("List() returned %d quotes, want at least 2", len(quotes))
		}

		// List by author
		author := "Author 1"
		authorQuotes, err := db.List(&author, nil)
		if err != nil {
			t.Fatalf("List() with author filter error = %v", err)
		}

		for _, quote := range authorQuotes {
			if quote.Author != author {
				t.Errorf("List() with author filter returned quote with author %v, want %v", quote.Author, author)
			}
		}

		// List by category
		category := "category2"
		categoryQuotes, err := db.List(nil, &category)
		if err != nil {
			t.Fatalf("List() with category filter error = %v", err)
		}

		for _, quote := range categoryQuotes {
			if quote.Category != category {
				t.Errorf("List() with category filter returned quote with category %v, want %v", quote.Category, category)
			}
		}

		// List by both author and category
		author2 := "Author 2"
		bothQuotes, err := db.List(&author2, &category)
		if err != nil {
			t.Fatalf("List() with both filters error = %v", err)
		}

		for _, quote := range bothQuotes {
			if quote.Author != author2 || quote.Category != category {
				t.Errorf("List() with both filters returned quote with author %v, category %v, want %v, %v", quote.Author, quote.Category, author2, category)
			}
		}
	})

	t.Run("GetRandom quote", func(t *testing.T) {
		// Create a quote first
		_, err := db.Create("Random test", "Random Author", "random")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Get random quote
		quote, err := db.GetRandom()
		if err != nil {
			t.Fatalf("GetRandom() error = %v", err)
		}

		if quote == nil {
			t.Error("GetRandom() returned nil quote")
		}
		if quote.ID == 0 {
			t.Error("GetRandom() returned quote with ID 0")
		}
	})

	t.Run("GetRandom from empty database", func(t *testing.T) {
		// Create new empty database
		emptyDbPath := filepath.Join(tempDir, "empty.db")
		emptyDb, err := NewDB(emptyDbPath)
		if err != nil {
			t.Fatalf("NewDB() for empty database error = %v", err)
		}
		defer emptyDb.Close()

		// Try to get random quote from empty database
		_, err = emptyDb.GetRandom()
		if err == nil {
			t.Error("GetRandom() from empty database should return error")
		}
	})

	t.Run("Update quote", func(t *testing.T) {
		// Create a quote
		quote, err := db.Create("Original text", "Original Author", "original")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Update quote
		newText := "Updated text"
		newAuthor := "Updated Author"
		newCategory := "updated"
		err = db.Update(quote.ID, newText, newAuthor, newCategory)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		// Verify update
		updated, err := db.Get(quote.ID)
		if err != nil {
			t.Fatalf("Get() after Update() error = %v", err)
		}

		if updated.Text != newText {
			t.Errorf("Update() text = %v, want %v", updated.Text, newText)
		}
		if updated.Author != newAuthor {
			t.Errorf("Update() author = %v, want %v", updated.Author, newAuthor)
		}
		if updated.Category != newCategory {
			t.Errorf("Update() category = %v, want %v", updated.Category, newCategory)
		}
	})

	t.Run("Update non-existent quote", func(t *testing.T) {
		err := db.Update(99999, "New text", "New Author", "new")
		// SQLite UPDATE doesn't return error for non-existent rows, it just affects 0 rows
		// This is the expected behavior for SQLite
		if err != nil {
			t.Errorf("Update() with non-existent ID should not return error, got %v", err)
		}
	})

	t.Run("Delete quote", func(t *testing.T) {
		// Create a quote
		quote, err := db.Create("Delete test", "Delete Author", "delete")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Delete the quote
		err = db.Delete(quote.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// Verify deletion
		_, err = db.Get(quote.ID)
		if err == nil {
			t.Error("Get() after Delete() should return error")
		}
	})

	t.Run("Delete non-existent quote", func(t *testing.T) {
		// This should not return an error in SQLite
		err := db.Delete(99999)
		if err != nil {
			t.Errorf("Delete() with non-existent ID should not return error, got %v", err)
		}
	})
}

func TestMigrateHardcodedQuotes(t *testing.T) {
	// Create temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "migrate_test.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer db.Close()

	// Test initial migration
	err = db.MigrateHardcodedQuotes()
	if err != nil {
		t.Fatalf("MigrateHardcodedQuotes() error = %v", err)
	}

	// Verify quotes were created
	quotes, err := db.List(nil, nil)
	if err != nil {
		t.Fatalf("List() after migration error = %v", err)
	}

	if len(quotes) == 0 {
		t.Error("MigrateHardcodedQuotes() should create quotes")
	}

	// Verify some quotes have authors
	hasAuthor := false
	for _, quote := range quotes {
		if quote.Author != "" {
			hasAuthor = true
			break
		}
	}
	if !hasAuthor {
		t.Error("MigrateHardcodedQuotes() should parse some quotes with authors")
	}

	// Test that migration doesn't duplicate quotes
	initialCount := len(quotes)
	err = db.MigrateHardcodedQuotes()
	if err != nil {
		t.Fatalf("Second MigrateHardcodedQuotes() error = %v", err)
	}

	quotes, err = db.List(nil, nil)
	if err != nil {
		t.Fatalf("List() after second migration error = %v", err)
	}

	if len(quotes) != initialCount {
		t.Errorf("Second MigrateHardcodedQuotes() changed quote count from %d to %d", initialCount, len(quotes))
	}
}

func TestParseQuote(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantText   string
		wantAuthor string
	}{
		{
			name:       "quote with author",
			input:      "Success is not final - Winston Churchill",
			wantText:   "Success is not final",
			wantAuthor: "Winston Churchill",
		},
		{
			name:       "quote without author",
			input:      "Just a simple quote",
			wantText:   "Just a simple quote",
			wantAuthor: "",
		},
		{
			name:       "quote with multiple dashes",
			input:      "Life is what happens - when you're busy - making other plans",
			wantText:   "Life is what happens",
			wantAuthor: "when you're busy",
		},
		{
			name:       "quote with emoji",
			input:      "🚀 The way to get started is to quit talking and begin doing. - Walt Disney",
			wantText:   "🚀 The way to get started is to quit talking and begin doing.",
			wantAuthor: "Walt Disney",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, gotAuthor := parseQuote(tt.input)
			if gotText != tt.wantText {
				t.Errorf("parseQuote() text = %v, want %v", gotText, tt.wantText)
			}
			if gotAuthor != tt.wantAuthor {
				t.Errorf("parseQuote() author = %v, want %v", gotAuthor, tt.wantAuthor)
			}
		})
	}
}
