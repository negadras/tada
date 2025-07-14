package quote

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"testing"
	"time"
)

func TestValidateQuoteText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "valid text",
			text:    "This is a valid quote",
			wantErr: false,
		},
		{
			name:    "empty text",
			text:    "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			text:    "   ",
			wantErr: true,
		},
		{
			name:    "very long text",
			text:    strings.Repeat("a", 1001),
			wantErr: true,
		},
		{
			name:    "max length text",
			text:    strings.Repeat("a", 1000),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuoteText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuoteText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAuthor(t *testing.T) {
	tests := []struct {
		name    string
		author  string
		wantErr bool
	}{
		{
			name:    "valid author",
			author:  "John Doe",
			wantErr: false,
		},
		{
			name:    "empty author",
			author:  "",
			wantErr: false,
		},
		{
			name:    "whitespace only",
			author:  "   ",
			wantErr: false,
		},
		{
			name:    "very long author",
			author:  strings.Repeat("a", 101),
			wantErr: true,
		},
		{
			name:    "max length author",
			author:  strings.Repeat("a", 100),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthor(tt.author)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAuthor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{
			name:     "valid category",
			category: "motivation",
			wantErr:  false,
		},
		{
			name:     "empty category",
			category: "",
			wantErr:  false,
		},
		{
			name:     "whitespace only",
			category: "   ",
			wantErr:  false,
		},
		{
			name:     "very long category",
			category: strings.Repeat("a", 51),
			wantErr:  true,
		},
		{
			name:     "max length category",
			category: strings.Repeat("a", 50),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCategory(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Mock command for testing print functions
type mockCommand struct {
	*cobra.Command
	output strings.Builder
}

func newMockCommand() *mockCommand {
	mockCmd := &mockCommand{
		Command: &cobra.Command{},
	}

	// Override the Printf method
	mockCmd.Command.SetOut(&mockCmd.output)

	return mockCmd
}

func (m *mockCommand) Printf(format string, args ...interface{}) {
	m.output.WriteString(fmt.Sprintf(format, args...))
}

func (m *mockCommand) Println(args ...interface{}) {
	m.output.WriteString(fmt.Sprintln(args...))
}

func TestPrintQuote(t *testing.T) {
	mockCmd := newMockCommand()

	quote := &Quote{
		ID:        1,
		Text:      "Test quote",
		Author:    "Test Author",
		Category:  "test",
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}

	PrintQuote(mockCmd.Command, quote)

	output := mockCmd.output.String()
	if !strings.Contains(output, "Test quote") {
		t.Error("PrintQuote() should contain quote text")
	}
	if !strings.Contains(output, "Test Author") {
		t.Error("PrintQuote() should contain author")
	}
}

func TestPrintQuoteWithoutAuthor(t *testing.T) {
	mockCmd := newMockCommand()

	quote := &Quote{
		ID:        1,
		Text:      "Test quote",
		Author:    "",
		Category:  "test",
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}

	PrintQuote(mockCmd.Command, quote)

	output := mockCmd.output.String()
	if !strings.Contains(output, "Test quote") {
		t.Error("PrintQuote() should contain quote text")
	}
	if strings.Contains(output, " - ") {
		t.Error("PrintQuote() should not contain author separator when no author")
	}
}

func TestPrintQuoteCreated(t *testing.T) {
	mockCmd := newMockCommand()

	quote := &Quote{
		ID:       42,
		Text:     "Test quote",
		Author:   "Test Author",
		Category: "test",
	}

	PrintQuoteCreated(mockCmd.Command, quote)

	output := mockCmd.output.String()
	if !strings.Contains(output, "✅ Added quote #42") {
		t.Error("PrintQuoteCreated() should contain success message with ID")
	}
	if !strings.Contains(output, "Test Author") {
		t.Error("PrintQuoteCreated() should contain author")
	}
}

func TestPrintQuoteCreatedWithoutAuthor(t *testing.T) {
	mockCmd := newMockCommand()

	quote := &Quote{
		ID:       42,
		Text:     "Test quote",
		Author:   "",
		Category: "test",
	}

	PrintQuoteCreated(mockCmd.Command, quote)

	output := mockCmd.output.String()
	if !strings.Contains(output, "✅ Added quote #42") {
		t.Error("PrintQuoteCreated() should contain success message with ID")
	}
	if strings.Contains(output, "Author:") {
		t.Error("PrintQuoteCreated() should not contain author line when no author")
	}
}

func TestPrintError(t *testing.T) {
	mockCmd := newMockCommand()

	err := fmt.Errorf("test error")
	PrintError(mockCmd.Command, err)

	output := mockCmd.output.String()
	if !strings.Contains(output, "❌ Error: test error") {
		t.Error("PrintError() should contain error message")
	}
}

func TestPrintSuccess(t *testing.T) {
	mockCmd := newMockCommand()

	message := "Operation successful"
	PrintSuccess(mockCmd.Command, message)

	output := mockCmd.output.String()
	if !strings.Contains(output, "✅ Operation successful") {
		t.Error("PrintSuccess() should contain success message")
	}
}
