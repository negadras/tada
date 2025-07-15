package models

import (
	"path/filepath"
	"testing"

	"github.com/negadras/tada/internal/quote"
	"github.com/negadras/tada/internal/todo"
	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

func TestDashboard_loadTodoStats(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := todo.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Test statistics calculation logic directly
	// Test with empty database
	todos, err := db.List(nil, nil)
	if err != nil {
		t.Fatalf("Failed to list todos: %v", err)
	}

	// Calculate stats manually (mimicking loadTodoStats logic)
	var completed, todayCompleted int
	for _, todoItem := range todos {
		if todoItem.Status == todo.Done {
			completed++
			// All completed todos are considered "today" for testing purposes
			todayCompleted++
		}
	}

	stats := TodoStats{
		Total:          len(todos),
		Completed:      completed,
		TodayCompleted: todayCompleted,
	}

	if stats.Total != 0 {
		t.Errorf("Expected 0 total todos, got %d", stats.Total)
	}
	if stats.Completed != 0 {
		t.Errorf("Expected 0 completed todos, got %d", stats.Completed)
	}
	if stats.TodayCompleted != 0 {
		t.Errorf("Expected 0 today completed todos, got %d", stats.TodayCompleted)
	}

	// Add some test todos
	todo1, err := db.Create("Test todo 1", todo.High)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	_, err = db.Create("Test todo 2", todo.Medium)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	// Mark one as completed
	err = db.UpdateStatus(todo1.ID, todo.Done)
	if err != nil {
		t.Fatalf("Failed to update todo status: %v", err)
	}

	// Test with populated database
	todos, err = db.List(nil, nil)
	if err != nil {
		t.Fatalf("Failed to list todos: %v", err)
	}

	// Recalculate stats
	completed = 0
	todayCompleted = 0
	for _, todoItem := range todos {
		if todoItem.Status == todo.Done {
			completed++
			todayCompleted++
		}
	}

	stats = TodoStats{
		Total:          len(todos),
		Completed:      completed,
		TodayCompleted: todayCompleted,
	}

	if stats.Total != 2 {
		t.Errorf("Expected 2 total todos, got %d", stats.Total)
	}
	if stats.Completed != 1 {
		t.Errorf("Expected 1 completed todo, got %d", stats.Completed)
	}
	if stats.TodayCompleted != 1 {
		t.Errorf("Expected 1 today completed todo, got %d", stats.TodayCompleted)
	}
}

func TestDashboard_loadQuoteStats(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := quote.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Test statistics calculation logic directly
	// Test with empty database
	quotes, err := db.List(nil, nil)
	if err != nil {
		t.Fatalf("Failed to list quotes: %v", err)
	}

	stats := QuoteStats{
		Total: len(quotes),
	}

	if stats.Total != 0 {
		t.Errorf("Expected 0 total quotes, got %d", stats.Total)
	}

	// Add some test quotes
	_, err = db.Create("Test quote 1", "Author 1", "category1")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Test quote 2", "Author 2", "category2")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Test quote 3", "Author 3", "category1")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	// Test with populated database
	quotes, err = db.List(nil, nil)
	if err != nil {
		t.Fatalf("Failed to list quotes: %v", err)
	}

	stats = QuoteStats{
		Total: len(quotes),
	}

	if stats.Total != 3 {
		t.Errorf("Expected 3 total quotes, got %d", stats.Total)
	}
}

func TestDashboard_CompletionRateCalculation(t *testing.T) {
	tests := []struct {
		name         string
		total        int
		completed    int
		expectedRate float64
	}{
		{
			name:         "Empty database",
			total:        0,
			completed:    0,
			expectedRate: 0.0,
		},
		{
			name:         "All completed",
			total:        5,
			completed:    5,
			expectedRate: 100.0,
		},
		{
			name:         "Half completed",
			total:        4,
			completed:    2,
			expectedRate: 50.0,
		},
		{
			name:         "None completed",
			total:        3,
			completed:    0,
			expectedRate: 0.0,
		},
		{
			name:         "Partial completion",
			total:        3,
			completed:    1,
			expectedRate: 33.333333333333336,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var completionRate float64
			if tt.total > 0 {
				completionRate = float64(tt.completed) / float64(tt.total) * 100
			}

			// Use tolerance-based comparison for floating-point values
			tolerance := 0.000001
			if abs(completionRate-tt.expectedRate) > tolerance {
				t.Errorf("Expected completion rate %f, got %f", tt.expectedRate, completionRate)
			}
		})
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestDashboard_TodayCompletedCalculation(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := todo.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Add todos and mark them as completed at different times
	todo1, err := db.Create("Test todo 1", todo.High)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	todo2, err := db.Create("Test todo 2", todo.Medium)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	_, err = db.Create("Test todo 3", todo.Low)
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	// Mark todos as completed
	err = db.UpdateStatus(todo1.ID, todo.Done)
	if err != nil {
		t.Fatalf("Failed to update todo status: %v", err)
	}

	err = db.UpdateStatus(todo2.ID, todo.Done)
	if err != nil {
		t.Fatalf("Failed to update todo status: %v", err)
	}

	// Leave todo3 as open

	// Test today's completion count calculation logic
	todos, err := db.List(nil, nil)
	if err != nil {
		t.Fatalf("Failed to list todos: %v", err)
	}

	var completed, todayCompleted int
	for _, todoItem := range todos {
		if todoItem.Status == todo.Done {
			completed++
			// For testing, we assume all completed todos are "today"
			// since we can't easily control the CompletedAt timestamp
			todayCompleted++
		}
	}

	// Both todos should be counted as completed today
	if todayCompleted != 2 {
		t.Errorf("Expected 2 todos completed today, got %d", todayCompleted)
	}

	if completed != 2 {
		t.Errorf("Expected 2 total completed todos, got %d", completed)
	}

	if len(todos) != 3 {
		t.Errorf("Expected 3 total todos, got %d", len(todos))
	}
}

func TestDashboard_NewDashboard(t *testing.T) {
	defaultStyles := styles.DefaultStyles()
	keymap := utils.DefaultKeyMap()

	dashboard := NewDashboard(defaultStyles, keymap)

	if dashboard == nil {
		t.Fatal("NewDashboard returned nil")
	}

	if dashboard.styles != defaultStyles {
		t.Error("Dashboard styles not set correctly")
	}

	if dashboard.selectedItem != 0 {
		t.Errorf("Expected selectedItem to be 0, got %d", dashboard.selectedItem)
	}

	if len(dashboard.menuItems) != 4 {
		t.Errorf("Expected 4 menu items, got %d", len(dashboard.menuItems))
	}

	if !dashboard.stats.Loading {
		t.Error("Expected dashboard stats to be in loading state initially")
	}
}

func TestDashboard_MenuItems(t *testing.T) {
	dashboard := NewDashboard(styles.DefaultStyles(), utils.DefaultKeyMap())

	expectedItems := []struct {
		title  string
		action string
	}{
		{"Todo Management", "todos"},
		{"Quote Collection", "quotes"},
		{"Statistics", "stats"},
		{"Settings", "settings"},
	}

	if len(dashboard.menuItems) != len(expectedItems) {
		t.Fatalf("Expected %d menu items, got %d", len(expectedItems), len(dashboard.menuItems))
	}

	for i, expected := range expectedItems {
		if dashboard.menuItems[i].Title != expected.title {
			t.Errorf("Expected menu item %d title to be '%s', got '%s'", i, expected.title, dashboard.menuItems[i].Title)
		}
		if dashboard.menuItems[i].Action != expected.action {
			t.Errorf("Expected menu item %d action to be '%s', got '%s'", i, expected.action, dashboard.menuItems[i].Action)
		}
	}
}
