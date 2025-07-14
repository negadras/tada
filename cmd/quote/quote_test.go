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

	if cmd.Long != "Manage your collection of motivational quotes with subcommands for adding, listing, updating, and deleting quotes. Running 'quote' without subcommands displays a random quote." {
		t.Errorf("NewCommand() Long = %v, want expected long description", cmd.Long)
	}
}
