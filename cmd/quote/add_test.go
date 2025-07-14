package quote

import "testing"

func TestNewAddCommand(t *testing.T) {
	cmd := newAddCommand()

	if cmd.Use != "add [text]" {
		t.Errorf("newAddCommand() Use = %v, want 'add [text]'", cmd.Use)
	}

	if cmd.Short != "Add a new quote" {
		t.Errorf("newAddCommand() Short = %v, want 'Add a new quote'", cmd.Short)
	}

	// Test that flags are set correctly
	authorFlag := cmd.Flags().Lookup("author")
	if authorFlag == nil {
		t.Error("newAddCommand() should have an author flag")
	}

	categoryFlag := cmd.Flags().Lookup("category")
	if categoryFlag == nil {
		t.Error("newAddCommand() should have a category flag")
	}

	// Test short flags
	if authorFlag.Shorthand != "a" {
		t.Errorf("newAddCommand() author flag shorthand = %v, want 'a'", authorFlag.Shorthand)
	}

	if categoryFlag.Shorthand != "c" {
		t.Errorf("newAddCommand() category flag shorthand = %v, want 'c'", categoryFlag.Shorthand)
	}

	// Test default values
	if authorFlag.DefValue != "" {
		t.Errorf("newAddCommand() author flag default = %v, want empty string", authorFlag.DefValue)
	}

	if categoryFlag.DefValue != "" {
		t.Errorf("newAddCommand() category flag default = %v, want empty string", categoryFlag.DefValue)
	}
}

func TestAddCommand_Arguments(t *testing.T) {
	cmd := newAddCommand()

	// Test that command requires exactly 1 argument
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("Expected error when no arguments provided")
	}

	cmd.SetArgs([]string{"arg1", "arg2"})
	if err := cmd.Execute(); err == nil {
		t.Error("Expected error when too many arguments provided")
	}

	if err := cmd.Args(cmd, []string{"valid quote text"}); err != nil {
		t.Errorf("Expected no error with valid arguments, got %v", err)
	}
}
