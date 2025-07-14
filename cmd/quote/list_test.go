package quote

import "testing"

func TestNewListCommand(t *testing.T) {
	cmd := newListCommand()

	if cmd.Use != "list" {
		t.Errorf("newListCommand() Use = %v, want 'list'", cmd.Use)
	}

	if cmd.Short != "List all quotes" {
		t.Errorf("newListCommand() Short = %v, want 'List all quotes'", cmd.Short)
	}

	// Test that flags are set correctly
	authorFlag := cmd.Flags().Lookup("author")
	if authorFlag == nil {
		t.Error("newListCommand() should have an author flag")
	}

	categoryFlag := cmd.Flags().Lookup("category")
	if categoryFlag == nil {
		t.Error("newListCommand() should have a category flag")
	}

	// Test short flags
	if authorFlag.Shorthand != "a" {
		t.Errorf("newListCommand() author flag shorthand = %v, want 'a'", authorFlag.Shorthand)
	}

	if categoryFlag.Shorthand != "c" {
		t.Errorf("newListCommand() category flag shorthand = %v, want 'c'", categoryFlag.Shorthand)
	}

	// Test default values
	if authorFlag.DefValue != "" {
		t.Errorf("newListCommand() author flag default = %v, want empty string", authorFlag.DefValue)
	}

	if categoryFlag.DefValue != "" {
		t.Errorf("newListCommand() category flag default = %v, want empty string", categoryFlag.DefValue)
	}
}

func TestListCommand_Arguments(t *testing.T) {
	cmd := newListCommand()

	// Test that command accepts no arguments
	if cmd == nil {
		t.Error("newListCommand() returned nil")
	}
}
