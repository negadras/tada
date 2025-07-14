package quote

import "testing"

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "quote" {
		t.Errorf("NewCommand() Use = %s, want %s", cmd.Use, "quote")
	}

	if cmd.Short != "Manage and display motivational quotes" {
		t.Errorf("NewCommand() Short = %v, want 'Manage and display motivational quotes'", cmd.Short)
	}

	expectedLong := "Manage your collection of motivational quotes with subcommands for adding, listing, updating, and deleting quotes. Running 'quote' without subcommands displays a random quote.\n\n💡 Tip: For interactive quote browsing and management, try 'tada --tui'"
	if cmd.Long != expectedLong {
		t.Errorf("NewCommand() Long = %v, want expected long description", cmd.Long)
	}
}
