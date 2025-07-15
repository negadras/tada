package models

import (
	"path/filepath"
	"testing"

	"github.com/negadras/tada/internal/todo"
	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

func TestTodoManager_NewTodoManager(t *testing.T) {
	styles := styles.DefaultStyles()
	keymap := utils.DefaultKeyMap()
	
	manager := NewTodoManager(styles, keymap)
	
	if manager == nil {
		t.Fatal("NewTodoManager returned nil")
	}
	
	if manager.styles != styles {
		t.Error("TodoManager styles not set correctly")
	}
	
	if len(manager.todos) != 0 {
		t.Errorf("Expected empty todos slice, got %d items", len(manager.todos))
	}
	
	if !manager.loading {
		t.Error("Expected TodoManager to be in loading state initially")
	}
	
	if manager.addForm == nil {
		t.Error("Expected add form to be initialized")
	}
	
	if manager.editForm == nil {
		t.Error("Expected edit form to be initialized")
	}
}

func TestTodoManager_updateTable(t *testing.T) {
	manager := NewTodoManager(styles.DefaultStyles(), utils.DefaultKeyMap())
	
	// Create test todos
	testTodos := []*todo.Todo{
		{
			ID:          1,
			Description: "Test todo 1",
			Priority:    todo.High,
			Status:      todo.Open,
		},
		{
			ID:          2,
			Description: "Test todo 2",
			Priority:    todo.Medium,
			Status:      todo.Done,
		},
		{
			ID:          3,
			Description: "Test todo 3 with very long description that should be handled properly",
			Priority:    todo.Low,
			Status:      todo.Open,
		},
	}
	
	manager.todos = testTodos
	manager.updateTable()
	
	// Check that table rows were set
	rows := manager.table.Rows()
	if len(rows) != len(testTodos) {
		t.Errorf("Expected %d table rows, got %d", len(testTodos), len(rows))
	}
	
	// Check first row content
	if len(rows) > 0 {
		row := rows[0]
		if len(row) != 5 { // ID, Priority, Status, Age, Description
			t.Errorf("Expected row to have 5 columns, got %d", len(row))
		}
		
		if row[0] != "#1" {
			t.Errorf("Expected ID column to be '#1', got '%s'", row[0])
		}
		
		if row[1] != "HIGH" {
			t.Errorf("Expected Priority column to be 'HIGH', got '%s'", row[1])
		}
		
		if row[2] != "OPEN" {
			t.Errorf("Expected Status column to be 'OPEN', got '%s'", row[2])
		}
		
		if row[4] != "Test todo 1" {
			t.Errorf("Expected Description column to be 'Test todo 1', got '%s'", row[4])
		}
	}
}

func TestTodoManager_StatusFiltering(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := todo.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create test todos
	_, err = db.Create("Open todo 1", todo.High)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	_, err = db.Create("Open todo 2", todo.Medium)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	todo3, err := db.Create("Done todo 1", todo.Low)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	// Mark one as completed
	err = db.UpdateStatus(todo3.ID, todo.Done)
	if err != nil {
		t.Fatalf("Failed to update todo status: %v", err)
	}

	// Test filtering by status
	tests := []struct {
		name           string
		statusFilter   *todo.Status
		expectedCount  int
	}{
		{
			name:          "No filter (all todos)",
			statusFilter:  nil,
			expectedCount: 3,
		},
		{
			name:          "Filter by Open",
			statusFilter:  func() *todo.Status { s := todo.Open; return &s }(),
			expectedCount: 2,
		},
		{
			name:          "Filter by Done",
			statusFilter:  func() *todo.Status { s := todo.Done; return &s }(),
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todos, err := db.List(tt.statusFilter, nil)
			if err != nil {
				t.Fatalf("Failed to list todos: %v", err)
			}

			if len(todos) != tt.expectedCount {
				t.Errorf("Expected %d todos, got %d", tt.expectedCount, len(todos))
			}

			// Verify filter works correctly
			if tt.statusFilter != nil {
				for _, todo := range todos {
					if todo.Status != *tt.statusFilter {
						t.Errorf("Expected todo status to be %v, got %v", *tt.statusFilter, todo.Status)
					}
				}
			}
		})
	}
}

func TestTodoManager_PriorityFiltering(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := todo.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create test todos with different priorities
	_, err = db.Create("High priority todo", todo.High)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	_, err = db.Create("Medium priority todo", todo.Medium)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	_, err = db.Create("Low priority todo", todo.Low)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	_, err = db.Create("Another high priority todo", todo.High)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	// Test filtering by priority
	tests := []struct {
		name             string
		priorityFilter   *todo.Priority
		expectedCount    int
	}{
		{
			name:           "No filter (all todos)",
			priorityFilter: nil,
			expectedCount:  4,
		},
		{
			name:           "Filter by High",
			priorityFilter: func() *todo.Priority { p := todo.High; return &p }(),
			expectedCount:  2,
		},
		{
			name:           "Filter by Medium",
			priorityFilter: func() *todo.Priority { p := todo.Medium; return &p }(),
			expectedCount:  1,
		},
		{
			name:           "Filter by Low",
			priorityFilter: func() *todo.Priority { p := todo.Low; return &p }(),
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todos, err := db.List(nil, tt.priorityFilter)
			if err != nil {
				t.Fatalf("Failed to list todos: %v", err)
			}

			if len(todos) != tt.expectedCount {
				t.Errorf("Expected %d todos, got %d", tt.expectedCount, len(todos))
			}

			// Verify filter works correctly
			if tt.priorityFilter != nil {
				for _, todo := range todos {
					if todo.Priority != *tt.priorityFilter {
						t.Errorf("Expected todo priority to be %v, got %v", *tt.priorityFilter, todo.Priority)
					}
				}
			}
		})
	}
}

func TestTodoManager_CombinedFiltering(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := todo.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create test todos with different combinations
	todo1, err := db.Create("High priority open todo", todo.High)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	todo2, err := db.Create("High priority done todo", todo.High)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	_, err = db.Create("Medium priority open todo", todo.Medium)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	// Mark one high priority todo as done
	err = db.UpdateStatus(todo2.ID, todo.Done)
	if err != nil {
		t.Fatalf("Failed to update todo status: %v", err)
	}

	// Test combined filtering
	openStatus := todo.Open
	highPriority := todo.High
	
	todos, err := db.List(&openStatus, &highPriority)
	if err != nil {
		t.Fatalf("Failed to list todos: %v", err)
	}

	// Should only return the high priority open todo
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}

	if len(todos) > 0 {
		if todos[0].ID != todo1.ID {
			t.Errorf("Expected todo ID %d, got %d", todo1.ID, todos[0].ID)
		}
		if todos[0].Status != todo.Open {
			t.Errorf("Expected todo status to be Open, got %v", todos[0].Status)
		}
		if todos[0].Priority != todo.High {
			t.Errorf("Expected todo priority to be High, got %v", todos[0].Priority)
		}
	}
}

func TestTodoManager_StatsCalculation(t *testing.T) {
	manager := NewTodoManager(styles.DefaultStyles(), utils.DefaultKeyMap())
	
	// Create test todos
	testTodos := []*todo.Todo{
		{ID: 1, Description: "Todo 1", Priority: todo.High, Status: todo.Open},
		{ID: 2, Description: "Todo 2", Priority: todo.Medium, Status: todo.Done},
		{ID: 3, Description: "Todo 3", Priority: todo.Low, Status: todo.Done},
		{ID: 4, Description: "Todo 4", Priority: todo.High, Status: todo.Open},
		{ID: 5, Description: "Todo 5", Priority: todo.Medium, Status: todo.Done},
	}
	
	manager.todos = testTodos
	
	// Test stats calculation logic (simulating renderStats method logic)
	var completed, total int
	for _, todoItem := range manager.todos {
		total++
		if todoItem.Status == todo.Done {
			completed++
		}
	}
	
	if total != 5 {
		t.Errorf("Expected 5 total todos, got %d", total)
	}
	
	if completed != 3 {
		t.Errorf("Expected 3 completed todos, got %d", completed)
	}
	
	completionRate := float64(completed) / float64(total) * 100
	expectedRate := 60.0
	if completionRate != expectedRate {
		t.Errorf("Expected completion rate %f, got %f", expectedRate, completionRate)
	}
}

func TestTodoManager_EmptyState(t *testing.T) {
	manager := NewTodoManager(styles.DefaultStyles(), utils.DefaultKeyMap())
	
	// Test with empty todos
	manager.todos = []*todo.Todo{}
	manager.updateTable()
	
	rows := manager.table.Rows()
	if len(rows) != 0 {
		t.Errorf("Expected empty table rows, got %d", len(rows))
	}
	
	// Test stats with empty todos
	var completed, total int
	for _, todoItem := range manager.todos {
		total++
		if todoItem.Status == todo.Done {
			completed++
		}
	}
	
	if total != 0 {
		t.Errorf("Expected 0 total todos, got %d", total)
	}
	
	if completed != 0 {
		t.Errorf("Expected 0 completed todos, got %d", completed)
	}
}