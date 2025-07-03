package add

import "testing"

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "add [description]" {
		t.Errorf("NewCommand() Use = %v, want 'add [description]'", cmd.Use)
	}

	if cmd.Short != "Add a todo task" {
		t.Errorf("NewCommand() Short = %v, want 'Add a todo task'", cmd.Short)
	}

	// Test that flags are set correctly
	priorityFlag := cmd.Flags().Lookup("priority")
	if priorityFlag == nil {
		t.Errorf("NewCommand() should have a priority flag")
	}

	if priorityFlag.DefValue != "medium" {
		t.Errorf("NewCommand() priority flag default = %v, want 'medium'", priorityFlag.DefValue)
	}
}

func TestAddCommand_Arguments(t *testing.T) {
	cmd := NewCommand()

	// test that command requires exactly 1 argument
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Errorf("Expected error when no arguments provided")
	}

	cmd.SetArgs([]string{"arg1", "arg2"})
	if err := cmd.Execute(); err == nil {
		t.Errorf("Expected error when too many arguments provided")
	}

	if err := cmd.Args(cmd, []string{"valid description"}); err != nil {
		t.Errorf("Expected no error with valid arguments, got %v", err)
	}
}
