package delete

import "testing"

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "delete [id]" {
		t.Errorf("NewCommand.Use == %s, want %s", cmd.Use, "delete")
	}

	if cmd.Short != "Delete a todo" {
		t.Errorf("NewCommand() Short = %v, want 'Delete a todo'", cmd.Short)
	}
}

func TestDeleteCommand_Arguments(t *testing.T) {
	cmd := NewCommand()

	// Test that command requires exactly 1 argument
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("Expected error when no arguments provided")
	}

	cmd.SetArgs([]string{"arg1", "arg2"})
	if err := cmd.Execute(); err == nil {
		t.Error("Expected error when no arguments provided")
	}
}
