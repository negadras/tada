package quote

import "testing"

func TestNewUpdateCommand(t *testing.T) {
	cmd := newUpdateCommand()

	if cmd.Use != "update [id]" {
		t.Errorf("newUpdateCommand() Use = %v, want 'update [id]'", cmd.Use)
	}

	if cmd.Short != "Update a quote" {
		t.Errorf("newUpdateCommand() Short = %v, want 'Update a quote'", cmd.Short)
	}

	// Test that flags are set correctly
	textFlag := cmd.Flags().Lookup("text")
	if textFlag == nil {
		t.Error("newUpdateCommand() should have a text flag")
	}

	authorFlag := cmd.Flags().Lookup("author")
	if authorFlag == nil {
		t.Error("newUpdateCommand() should have an author flag")
	}

	categoryFlag := cmd.Flags().Lookup("category")
	if categoryFlag == nil {
		t.Error("newUpdateCommand() should have a category flag")
	}

	// Test short flags
	if textFlag.Shorthand != "t" {
		t.Errorf("newUpdateCommand() text flag shorthand = %v, want 't'", textFlag.Shorthand)
	}

	if authorFlag.Shorthand != "a" {
		t.Errorf("newUpdateCommand() author flag shorthand = %v, want 'a'", authorFlag.Shorthand)
	}

	if categoryFlag.Shorthand != "c" {
		t.Errorf("newUpdateCommand() category flag shorthand = %v, want 'c'", categoryFlag.Shorthand)
	}

	// Test default values
	if textFlag.DefValue != "" {
		t.Errorf("newUpdateCommand() text flag default = %v, want empty string", textFlag.DefValue)
	}

	if authorFlag.DefValue != "" {
		t.Errorf("newUpdateCommand() author flag default = %v, want empty string", authorFlag.DefValue)
	}

	if categoryFlag.DefValue != "" {
		t.Errorf("newUpdateCommand() category flag default = %v, want empty string", categoryFlag.DefValue)
	}
}

func TestUpdateCommand_Arguments(t *testing.T) {
	cmd := newUpdateCommand()

	// Test that command requires exactly 1 argument
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("Expected error when no arguments provided")
	}

	cmd.SetArgs([]string{"1", "2"})
	if err := cmd.Execute(); err == nil {
		t.Error("Expected error when too many arguments provided")
	}

	if err := cmd.Args(cmd, []string{"1"}); err != nil {
		t.Errorf("Expected no error with valid arguments, got %v", err)
	}
}
