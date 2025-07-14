package quote

import (
	"strconv"

	"github.com/negadras/tada/internal/quote"
	"github.com/spf13/cobra"
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a quote",
		Long:  "Delete a quote by its ID.",
		Example: `  # Delete quote #5
  tada quote delete 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			db, cleanup, err := quote.GetDB(cmd)
			if err != nil {
				return nil
			}
			defer cleanup()

			quoteItem, err := db.Get(id)
			if err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			if err := db.Delete(id); err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			cmd.Printf("🗑️  Deleted quote #%d: %s\n", id, quoteItem.Text)
			return nil
		},
	}

	return cmd
}
