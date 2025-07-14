package quote

import (
	"strings"

	"github.com/negadras/tada/internal/quote"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [text]",
		Short: "Add a new quote",
		Long: `Add a new motivational quote to the database.
You can specify the author and category for better organization.`,
		Example: `  # Add a simple quote
  tada quote add "Success is not final, failure is not fatal."
  
  # Add with author
  tada quote add "The only way to do great work is to love what you do." --author "Steve Jobs"
  
  # Add with author and category
  tada quote add "Be yourself; everyone else is already taken." --author "Oscar Wilde" --category "inspiration"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, cleanup, err := quote.GetDB(cmd)
			if err != nil {
				return nil
			}
			defer cleanup()

			text := strings.TrimSpace(args[0])
			if err := quote.ValidateQuoteText(text); err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			author, _ := cmd.Flags().GetString("author")
			category, _ := cmd.Flags().GetString("category")

			if err := quote.ValidateAuthor(author); err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			if err := quote.ValidateCategory(category); err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			newQuote, err := db.Create(text, author, category)
			if err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			quote.PrintQuoteCreated(cmd, newQuote)
			return nil
		},
	}

	cmd.Flags().StringP("author", "a", "", "Author of the quote")
	cmd.Flags().StringP("category", "c", "", "Category of the quote")
	return cmd
}
