package list

import "testing"

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "list" {
		t.Errorf("NewCommand() Use = %v, want 'list'", cmd.Use)
	}

	if cmd.Short != "List todo tasks" {
		t.Errorf("NewCommand() Short = %v, want 'List todo tasks'", cmd.Short)
	}

	// test that flags are set correctly
	statusFlag := cmd.Flags().Lookup("status")
	if statusFlag == nil {
		t.Errorf("NewCommand() should have flag 'status'")
	}

	priorityFlag := cmd.Flags().Lookup("priority")
	if priorityFlag == nil {
		t.Errorf("NewCommand() should have flag 'priority'")
	}
}

func TestNewCommand_Flags(t *testing.T) {
	cmd := NewCommand()

	// test short flags exist
	statusFlags := cmd.Flags().ShorthandLookup("s")
	if statusFlags == nil {
		t.Errorf("NewCommand() should have short status flag 's'")
	}

	priorityFlags := cmd.Flags().ShorthandLookup("p")
	if priorityFlags == nil {
		t.Errorf("NewCommand() should have short priority flag 'p'")
	}

}
