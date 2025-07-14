package quote

import (
	"github.com/negadras/tada/internal/quote"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all quotes",
		Long:  `List all quotes with optional filtering by author or category.`,
		Example: `  # List all quotes
  tada quote list
  
  # List quotes by specific author
  tada quote list --author "Steve Jobs"
  
  # List quotes by category
  tada quote list --category "motivation"
  
  # List quotes by author and category
  tada quote list --author "Oscar Wilde" --category "inspiration"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, cleanup, err := quote.GetDB(cmd)
			if err != nil {
				return nil
			}
			defer cleanup()

			author, _ := cmd.Flags().GetString("author")
			category, _ := cmd.Flags().GetString("category")

			var authorFilter, categoryFilter *string
			if author != "" {
				authorFilter = &author
			}
			if category != "" {
				categoryFilter = &category
			}

			quotes, err := db.List(authorFilter, categoryFilter)
			if err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			if len(quotes) == 0 {
				cmd.Println("No quotes found.")
				return nil
			}

			for i, q := range quotes {
				if i > 0 {
					cmd.Println()
				}
				quote.PrintQuote(cmd, q)
			}

			return nil
		},
	}

	cmd.Flags().StringP("author", "a", "", "Filter by author")
	cmd.Flags().StringP("category", "c", "", "Filter by category")

	return cmd
}
