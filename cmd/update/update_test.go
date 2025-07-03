package update

import "testing"

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "update [id]" {
		t.Errorf("NewCommand.Use = %s, want update", cmd.Use)
	}

	if cmd.Short != "Update a todo task" {
		t.Errorf("NewCommand.Short = %s, want update", cmd.Short)
	}

	// Test that flags are set correctly
	statusFlag := cmd.Flags().Lookup("status")
	if statusFlag == nil {
		t.Error("NewCommand() should have status flag")
	}

	priorityFlag := cmd.Flags().Lookup("priority")
	if priorityFlag == nil {
		t.Error("NewCommand() should have priority flag")
	}

	descFlag := cmd.Flags().Lookup("description")
	if descFlag == nil {
		t.Error("NewCommand() should have description flag")
	}
}

func TestUpdateCommand_Arguments(t *testing.T) {
	cmd := NewCommand()

	// Test that command requires exactly 1 argument
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}

	cmd.SetArgs([]string{"arg1", "arg2"})
	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error when too many arguments provided")
	}
}

func TestUpdateCommand_Flags(t *testing.T) {
	cmd := NewCommand()

	// Test short flags exist
	statusFlag := cmd.Flags().ShorthandLookup("s")
	if statusFlag == nil {
		t.Error("NewCommand() should have short status flag 's'")
	}

	priorityFlag := cmd.Flags().ShorthandLookup("p")
	if priorityFlag == nil {
		t.Error("NewCommand() should have short priority flag 'p'")
	}

	descFlag := cmd.Flags().ShorthandLookup("d")
	if descFlag == nil {
		t.Error("NewCommand() should have short description flag 'd'")
	}
}
