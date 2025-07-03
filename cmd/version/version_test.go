package version

import "testing"

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "version" {
		t.Errorf("NewCommand.Use = %s, want version", cmd.Use)
	}

	if cmd.Short != "Show version information" {
		t.Errorf("NewCommand.Short = %v, want 'Show version information'", cmd.Short)
	}
}
