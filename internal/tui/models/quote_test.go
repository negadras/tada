package models

import (
	"path/filepath"
	"testing"

	"github.com/negadras/tada/internal/quote"
	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

func TestQuoteManager_NewQuoteManager(t *testing.T) {
	styles := styles.DefaultStyles()
	keymap := utils.DefaultKeyMap()
	
	manager := NewQuoteManager(styles, keymap)
	
	if manager == nil {
		t.Fatal("NewQuoteManager returned nil")
	}
	
	if manager.styles != styles {
		t.Error("QuoteManager styles not set correctly")
	}
	
	if len(manager.quotes) != 0 {
		t.Errorf("Expected empty quotes slice, got %d items", len(manager.quotes))
	}
	
	if !manager.loading {
		t.Error("Expected QuoteManager to be in loading state initially")
	}
	
	if manager.addForm == nil {
		t.Error("Expected add form to be initialized")
	}
	
	if manager.editForm == nil {
		t.Error("Expected edit form to be initialized")
	}
}

func TestQuoteManager_updateTable(t *testing.T) {
	manager := NewQuoteManager(styles.DefaultStyles(), utils.DefaultKeyMap())
	
	// Create test quotes
	testQuotes := []*quote.Quote{
		{
			ID:       1,
			Text:     "Test quote 1",
			Author:   "Author 1",
			Category: "motivation",
		},
		{
			ID:       2,
			Text:     "Test quote 2",
			Author:   "",
			Category: "",
		},
		{
			ID:       3,
			Text:     "Test quote 3 with very long text that should be handled properly",
			Author:   "Author 3",
			Category: "inspiration",
		},
	}
	
	manager.quotes = testQuotes
	manager.updateTable()
	
	// Check that table rows were set
	rows := manager.table.Rows()
	if len(rows) != len(testQuotes) {
		t.Errorf("Expected %d table rows, got %d", len(testQuotes), len(rows))
	}
	
	// Check first row content
	if len(rows) > 0 {
		row := rows[0]
		if len(row) != 5 { // ID, Author, Category, Age, Quote
			t.Errorf("Expected row to have 5 columns, got %d", len(row))
		}
		
		if row[0] != "#1" {
			t.Errorf("Expected ID column to be '#1', got '%s'", row[0])
		}
		
		if row[1] != "Author 1" {
			t.Errorf("Expected Author column to be 'Author 1', got '%s'", row[1])
		}
		
		if row[2] != "motivation" {
			t.Errorf("Expected Category column to be 'motivation', got '%s'", row[2])
		}
		
		if row[4] != "Test quote 1" {
			t.Errorf("Expected Quote column to be 'Test quote 1', got '%s'", row[4])
		}
	}
	
	// Check second row (empty author and category)
	if len(rows) > 1 {
		row := rows[1]
		if row[1] != "Unknown" {
			t.Errorf("Expected empty author to be 'Unknown', got '%s'", row[1])
		}
		
		if row[2] != "None" {
			t.Errorf("Expected empty category to be 'None', got '%s'", row[2])
		}
	}
}

func TestQuoteManager_CategoryFiltering(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := quote.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create test quotes with different categories
	_, err = db.Create("Quote 1", "Author 1", "motivation")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Quote 2", "Author 2", "inspiration")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Quote 3", "Author 3", "motivation")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Quote 4", "Author 4", "")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	// Test filtering by category
	tests := []struct {
		name             string
		categoryFilter   *string
		expectedCount    int
	}{
		{
			name:           "No filter (all quotes)",
			categoryFilter: nil,
			expectedCount:  4,
		},
		{
			name:           "Filter by motivation",
			categoryFilter: func() *string { s := "motivation"; return &s }(),
			expectedCount:  2,
		},
		{
			name:           "Filter by inspiration",
			categoryFilter: func() *string { s := "inspiration"; return &s }(),
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes, err := db.List(nil, tt.categoryFilter)
			if err != nil {
				t.Fatalf("Failed to list quotes: %v", err)
			}

			if len(quotes) != tt.expectedCount {
				t.Errorf("Expected %d quotes, got %d", tt.expectedCount, len(quotes))
			}

			// Verify filter works correctly
			if tt.categoryFilter != nil {
				for _, quote := range quotes {
					if quote.Category != *tt.categoryFilter {
						t.Errorf("Expected quote category to be %v, got %v", *tt.categoryFilter, quote.Category)
					}
				}
			}
		})
	}
}

func TestQuoteManager_AuthorFiltering(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := quote.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create test quotes with different authors
	_, err = db.Create("Quote 1", "Einstein", "science")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Quote 2", "Jobs", "innovation")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Quote 3", "Einstein", "science")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Quote 4", "", "motivation")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	// Test filtering by author
	tests := []struct {
		name           string
		authorFilter   *string
		expectedCount  int
	}{
		{
			name:          "No filter (all quotes)",
			authorFilter:  nil,
			expectedCount: 4,
		},
		{
			name:          "Filter by Einstein",
			authorFilter:  func() *string { s := "Einstein"; return &s }(),
			expectedCount: 2,
		},
		{
			name:          "Filter by Jobs",
			authorFilter:  func() *string { s := "Jobs"; return &s }(),
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes, err := db.List(tt.authorFilter, nil)
			if err != nil {
				t.Fatalf("Failed to list quotes: %v", err)
			}

			if len(quotes) != tt.expectedCount {
				t.Errorf("Expected %d quotes, got %d", tt.expectedCount, len(quotes))
			}

			// Verify filter works correctly
			if tt.authorFilter != nil {
				for _, quote := range quotes {
					if quote.Author != *tt.authorFilter {
						t.Errorf("Expected quote author to be %v, got %v", *tt.authorFilter, quote.Author)
					}
				}
			}
		})
	}
}

func TestQuoteManager_RandomQuote(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := quote.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Test with empty database
	_, err = db.GetRandom()
	if err == nil {
		t.Error("Expected error when getting random quote from empty database")
	}

	// Add quotes
	quote1, err := db.Create("Quote 1", "Author 1", "category1")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	quote2, err := db.Create("Quote 2", "Author 2", "category2")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	quote3, err := db.Create("Quote 3", "Author 3", "category3")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	// Test with populated database
	randomQuote, err := db.GetRandom()
	if err != nil {
		t.Fatalf("Failed to get random quote: %v", err)
	}

	if randomQuote == nil {
		t.Fatal("GetRandom returned nil quote")
	}

	// Verify it's one of our quotes
	validIDs := []int{quote1.ID, quote2.ID, quote3.ID}
	found := false
	for _, id := range validIDs {
		if randomQuote.ID == id {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Random quote ID %d not found in expected IDs %v", randomQuote.ID, validIDs)
	}
}

func TestQuoteManager_StatsCalculation(t *testing.T) {
	manager := NewQuoteManager(styles.DefaultStyles(), utils.DefaultKeyMap())
	
	// Create test quotes
	testQuotes := []*quote.Quote{
		{ID: 1, Text: "Quote 1", Author: "Author 1", Category: "motivation"},
		{ID: 2, Text: "Quote 2", Author: "Author 2", Category: "inspiration"},
		{ID: 3, Text: "Quote 3", Author: "Author 3", Category: "motivation"},
		{ID: 4, Text: "Quote 4", Author: "Author 4", Category: "life"},
		{ID: 5, Text: "Quote 5", Author: "Author 5", Category: "motivation"},
	}
	
	manager.quotes = testQuotes
	
	// Test stats calculation logic (simulating renderStats method logic)
	if len(manager.quotes) == 0 {
		t.Error("Expected quotes to be populated")
	}
	
	// Count quotes by category
	categoryCount := make(map[string]int)
	for _, quote := range manager.quotes {
		if quote.Category != "" {
			categoryCount[quote.Category]++
		}
	}
	
	expectedCounts := map[string]int{
		"motivation":  3,
		"inspiration": 1,
		"life":        1,
	}
	
	for category, expectedCount := range expectedCounts {
		if categoryCount[category] != expectedCount {
			t.Errorf("Expected %d quotes in category %s, got %d", expectedCount, category, categoryCount[category])
		}
	}
	
	// Test total count
	totalQuotes := len(manager.quotes)
	if totalQuotes != 5 {
		t.Errorf("Expected 5 total quotes, got %d", totalQuotes)
	}
}

func TestQuoteManager_EmptyState(t *testing.T) {
	manager := NewQuoteManager(styles.DefaultStyles(), utils.DefaultKeyMap())
	
	// Test with empty quotes
	manager.quotes = []*quote.Quote{}
	manager.updateTable()
	
	rows := manager.table.Rows()
	if len(rows) != 0 {
		t.Errorf("Expected empty table rows, got %d", len(rows))
	}
	
	// Test stats with empty quotes
	if len(manager.quotes) != 0 {
		t.Errorf("Expected empty quotes slice, got %d", len(manager.quotes))
	}
}

func TestQuoteManager_CombinedFiltering(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := quote.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create test quotes with different combinations
	einstein1, err := db.Create("Science quote 1", "Einstein", "science")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Science quote 2", "Hawking", "science")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	_, err = db.Create("Life quote 1", "Einstein", "life")
	if err != nil {
		t.Fatalf("Failed to create test quote: %v", err)
	}

	// Test combined filtering (Einstein + science)
	author := "Einstein"
	category := "science"
	
	quotes, err := db.List(&author, &category)
	if err != nil {
		t.Fatalf("Failed to list quotes: %v", err)
	}

	// Should only return Einstein's science quote
	if len(quotes) != 1 {
		t.Errorf("Expected 1 quote, got %d", len(quotes))
	}

	if len(quotes) > 0 {
		if quotes[0].ID != einstein1.ID {
			t.Errorf("Expected quote ID %d, got %d", einstein1.ID, quotes[0].ID)
		}
		if quotes[0].Author != "Einstein" {
			t.Errorf("Expected quote author to be Einstein, got %s", quotes[0].Author)
		}
		if quotes[0].Category != "science" {
			t.Errorf("Expected quote category to be science, got %s", quotes[0].Category)
		}
	}
}

func TestQuoteManager_CategoryCycling(t *testing.T) {
	manager := NewQuoteManager(styles.DefaultStyles(), utils.DefaultKeyMap())
	
	// Create test quotes with different categories
	testQuotes := []*quote.Quote{
		{ID: 1, Text: "Quote 1", Author: "Author 1", Category: "motivation"},
		{ID: 2, Text: "Quote 2", Author: "Author 2", Category: "inspiration"},
		{ID: 3, Text: "Quote 3", Author: "Author 3", Category: "motivation"},
		{ID: 4, Text: "Quote 4", Author: "Author 4", Category: "life"},
		{ID: 5, Text: "Quote 5", Author: "Author 5", Category: ""},
	}
	
	manager.quotes = testQuotes
	
	// Test category collection logic (simulating cycleCategoryFilter method logic)
	categories := make(map[string]bool)
	for _, quote := range manager.quotes {
		if quote.Category != "" {
			categories[quote.Category] = true
		}
	}
	
	// Convert to slice
	var categoryList []string
	for category := range categories {
		categoryList = append(categoryList, category)
	}
	
	// Should have 3 unique categories (empty categories are ignored)
	if len(categoryList) != 3 {
		t.Errorf("Expected 3 unique categories, got %d", len(categoryList))
	}
	
	// Verify all expected categories are present
	expectedCategories := map[string]bool{
		"motivation":  true,
		"inspiration": true,
		"life":        true,
	}
	
	for _, category := range categoryList {
		if !expectedCategories[category] {
			t.Errorf("Unexpected category found: %s", category)
		}
	}
}