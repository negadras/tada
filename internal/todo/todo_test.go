package todo

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPriority_String(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		want     string
	}{
		{"Low priority", Low, "LOW"},
		{"Medium priority", Medium, "MEDIUM"},
		{"High priority", High, "HIGH"},
		{"Unknown priority", Priority(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.priority.String(); got != tt.want {
				t.Errorf("Priority.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   string
	}{
		{"Open status", Open, "OPEN"},
		{"Done status", Done, "DONE"},
		{"Unknown status", Status(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodo_Age(t *testing.T) {
	now := time.Now()
	todo := &Todo{
		CreatedAt: now.Add(-1 * time.Hour),
	}

	age := todo.Age()
	if age < 59*time.Minute || age > 61*time.Minute {
		t.Errorf("Todo.Age() = %v, want approximately 1 hour", age)
	}
}

func TestTodo_CompletedAge(t *testing.T) {
	t.Run("completed todo", func(t *testing.T) {
		completedTime := time.Now().Add(-30 * time.Minute)
		todo := &Todo{
			CompletedAt: &completedTime,
		}

		age := todo.CompletedAge()
		if age == nil {
			t.Fatal("Todo.CompletedAge() returned nil for completed todo")
		}

		if *age < 29*time.Minute || *age > 31*time.Minute {
			t.Errorf("Todo.CompletedAge() = %v, want approximately 30 minutes", *age)
		}
	})

	t.Run("open todo", func(t *testing.T) {
		todo := &Todo{
			CompletedAt: nil,
		}

		age := todo.CompletedAge()
		if age != nil {
			t.Errorf("Todo.CompletedAge() = %v, want nil for open todo", age)
		}
	})
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

	t.Run("Create todo", func(t *testing.T) {
		todo, err := db.Create("Test todo", High)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if todo.ID == 0 {
			t.Error("Create() returned todo with ID 0")
		}
		if todo.Description != "Test todo" {
			t.Errorf("Create() description = %v, want %v", todo.Description, "Test todo")
		}
		if todo.Priority != High {
			t.Errorf("Create() priority = %v, want %v", todo.Priority, High)
		}
		if todo.Status != Open {
			t.Errorf("Create() status = %v, want %v", todo.Status, Open)
		}
	})

	t.Run("Get todo", func(t *testing.T) {
		// Create a todo first
		created, err := db.Create("Get test", Medium)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Get the todo
		retrieved, err := db.Get(created.ID)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if retrieved.ID != created.ID {
			t.Errorf("Get() ID = %v, want %v", retrieved.ID, created.ID)
		}
		if retrieved.Description != created.Description {
			t.Errorf("Get() description = %v, want %v", retrieved.Description, created.Description)
		}
	})

	t.Run("List todos", func(t *testing.T) {
		// Create multiple todos
		_, err := db.Create("List test 1", Low)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		_, err = db.Create("List test 2", High)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// List all todos
		todos, err := db.List(nil, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(todos) < 2 {
			t.Errorf("List() returned %d todos, want at least 2", len(todos))
		}

		// List by status
		openStatus := Open
		openTodos, err := db.List(&openStatus, nil)
		if err != nil {
			t.Fatalf("List() with status filter error = %v", err)
		}

		for _, todo := range openTodos {
			if todo.Status != Open {
				t.Errorf("List() with status filter returned todo with status %v, want %v", todo.Status, Open)
			}
		}

		// List by priority
		highPriority := High
		highTodos, err := db.List(nil, &highPriority)
		if err != nil {
			t.Fatalf("List() with priority filter error = %v", err)
		}

		for _, todo := range highTodos {
			if todo.Priority != High {
				t.Errorf("List() with priority filter returned todo with priority %v, want %v", todo.Priority, High)
			}
		}
	})

	t.Run("Update status", func(t *testing.T) {
		// Create a todo
		todo, err := db.Create("Status test", Medium)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Update status to done
		err = db.UpdateStatus(todo.ID, Done)
		if err != nil {
			t.Fatalf("UpdateStatus() error = %v", err)
		}

		// Verify update
		updated, err := db.Get(todo.ID)
		if err != nil {
			t.Fatalf("Get() after UpdateStatus() error = %v", err)
		}

		if updated.Status != Done {
			t.Errorf("UpdateStatus() status = %v, want %v", updated.Status, Done)
		}
		if updated.CompletedAt == nil {
			t.Error("UpdateStatus() to Done should set CompletedAt")
		}

		// Update back to open
		err = db.UpdateStatus(todo.ID, Open)
		if err != nil {
			t.Fatalf("UpdateStatus() back to Open error = %v", err)
		}

		reopened, err := db.Get(todo.ID)
		if err != nil {
			t.Fatalf("Get() after reopening error = %v", err)
		}

		if reopened.Status != Open {
			t.Errorf("UpdateStatus() back to Open status = %v, want %v", reopened.Status, Open)
		}
		if reopened.CompletedAt != nil {
			t.Error("UpdateStatus() back to Open should clear CompletedAt")
		}
	})

	t.Run("Update priority", func(t *testing.T) {
		// Create a todo
		todo, err := db.Create("Priority test", Low)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Update priority
		err = db.UpdatePriority(todo.ID, High)
		if err != nil {
			t.Fatalf("UpdatePriority() error = %v", err)
		}

		// Verify update
		updated, err := db.Get(todo.ID)
		if err != nil {
			t.Fatalf("Get() after UpdatePriority() error = %v", err)
		}

		if updated.Priority != High {
			t.Errorf("UpdatePriority() priority = %v, want %v", updated.Priority, High)
		}
	})

	t.Run("Update description", func(t *testing.T) {
		// Create a todo
		todo, err := db.Create("Original description", Medium)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Update description
		newDesc := "Updated description"
		err = db.UpdateDescription(todo.ID, newDesc)
		if err != nil {
			t.Fatalf("UpdateDescription() error = %v", err)
		}

		// Verify update
		updated, err := db.Get(todo.ID)
		if err != nil {
			t.Fatalf("Get() after UpdateDescription() error = %v", err)
		}

		if updated.Description != newDesc {
			t.Errorf("UpdateDescription() description = %v, want %v", updated.Description, newDesc)
		}
	})

	t.Run("Delete todo", func(t *testing.T) {
		// Create a todo
		todo, err := db.Create("Delete test", Medium)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Delete the todo
		err = db.Delete(todo.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// Verify deletion
		_, err = db.Get(todo.ID)
		if err == nil {
			t.Error("Get() after Delete() should return error")
		}
	})
}

func TestGetDatabasePath(t *testing.T) {
	// Save original home
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set temporary home
	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)

	path, err := GetDatabasePath()
	if err != nil {
		t.Fatalf("GetDatabasePath() error = %v", err)
	}

	expectedPath := filepath.Join(tempHome, ".tada", "todos.db")
	if path != expectedPath {
		t.Errorf("GetDatabasePath() = %v, want %v", path, expectedPath)
	}

	// Verify directory was created
	tadaDir := filepath.Join(tempHome, ".tada")
	if _, err := os.Stat(tadaDir); os.IsNotExist(err) {
		t.Error("GetDatabasePath() should create .tada directory")
	}
}
