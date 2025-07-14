package quote

import "testing"

func TestNewDeleteCommand(t *testing.T) {
	cmd := newDeleteCommand()

	if cmd.Use != "delete [id]" {
		t.Errorf("newDeleteCommand() Use = %v, want 'delete [id]'", cmd.Use)
	}

	if cmd.Short != "Delete a quote" {
		t.Errorf("newDeleteCommand() Short = %v, want 'Delete a quote'", cmd.Short)
	}

	// Test that the command has no flags (delete command is simple)
	if cmd.Flags().NFlag() != 0 {
		t.Errorf("newDeleteCommand() should have no flags, got %d", cmd.Flags().NFlag())
	}
}

func TestDeleteCommand_Arguments(t *testing.T) {
	cmd := newDeleteCommand()

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
