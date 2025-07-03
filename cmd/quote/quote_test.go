package quote

import "testing"

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "quote" {
		t.Errorf("NewCommand() Use = %s, want %s", cmd.Use, "quote")
	}

	if cmd.Short != "Show a motivational quote" {
		t.Errorf("NewCommand() Short = %v, want 'Show a motivational quote'", cmd.Short)
	}

	if cmd.Long != "Display a random motivational quote to inspire productivity" {
		t.Errorf("NewCommand() Long = %v, want expected long description", cmd.Long)
	}
}
