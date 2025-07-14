package quote

import (
	"github.com/negadras/tada/internal/quote"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote",
		Short: "Manage and display motivational quotes",
		Long:  "Manage your collection of motivational quotes with subcommands for adding, listing, updating, and deleting quotes. Running 'quote' without subcommands displays a random quote.",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, cleanup, err := quote.GetDB(cmd)
			if err != nil {
				return err
			}
			defer cleanup()

			// Migrate hardcoded quotes on first run
			if err := db.MigrateHardcodedQuotes(); err != nil {
				quote.PrintError(cmd, err)
				return err
			}

			// Get a random quote from the database
			randomQuote, err := db.GetRandom()
			if err != nil {
				quote.PrintError(cmd, err)
				return err
			}

			quote.PrintQuote(cmd, randomQuote)
			return nil
		},
	}

	// Add subcommands
	cmd.AddCommand(newAddCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newUpdateCommand())
	cmd.AddCommand(newDeleteCommand())

	return cmd
}
