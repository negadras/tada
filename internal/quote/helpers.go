package quote

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

// ValidateQuoteText checks if a quote text is valid
func ValidateQuoteText(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("quote text cannot be empty")
	}
	if len(text) > 1000 {
		return errors.New("quote text too long (max 1000 characters)")
	}
	return nil
}

// ValidateAuthor checks if an author name is valid
func ValidateAuthor(author string) error {
	author = strings.TrimSpace(author)
	if len(author) > 100 {
		return errors.New("author name too long (max 100 characters)")
	}
	return nil
}

// ValidateCategory checks if a category is valid
func ValidateCategory(category string) error {
	category = strings.TrimSpace(category)
	if len(category) > 50 {
		return errors.New("category name too long (max 50 characters)")
	}
	return nil
}

// PrintQuote prints a quote in a formatted way
func PrintQuote(cmd *cobra.Command, quote *Quote) {
	cmd.Printf("\n %s\n", quote.Text)
	if quote.Author != "" {
		cmd.Printf("   - %s\n", quote.Author)
	}
}

// PrintQuoteCreated prints a success message for created quote
func PrintQuoteCreated(cmd *cobra.Command, quote *Quote) {
	cmd.Printf("✅ Added quote #%d\n", quote.ID)
	if quote.Author != "" {
		cmd.Printf("   Author: %s\n", quote.Author)
	}
	if quote.Category != "" {
		cmd.Printf("   Category: %s\n", quote.Category)
	}
}

// PrintError prints an error message
func PrintError(cmd *cobra.Command, err error) {
	cmd.Printf("❌ Error: %v\n", err)
}

// PrintSuccess prints a success message
func PrintSuccess(cmd *cobra.Command, message string) {
	cmd.Printf("✅ %s\n", message)
}

// GetDB returns a database connection with cleanup function
func GetDB(cmd *cobra.Command) (*DB, func(), error) {
	dbPath, err := GetDatabasePath()
	if err != nil {
		PrintError(cmd, err)
		return nil, nil, err
	}

	db, err := NewDB(dbPath)
	if err != nil {
		PrintError(cmd, err)
		return nil, nil, err
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup, nil
}
