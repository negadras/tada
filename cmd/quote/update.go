package quote

import (
	"fmt"
	"strconv"

	"github.com/negadras/tada/internal/quote"
	"github.com/spf13/cobra"
)

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a quote",
		Long: `Update various properties of a quote including text, author, and category.
		
At least one flag must be provided to specify what to update.`,
		Example: `  # Update quote text
  tada quote update 5 --text "New quote text"
  
  # Update author
  tada quote update 5 --author "New Author"
  
  # Update category
  tada quote update 5 --category "inspiration"
  
  # Update multiple properties at once
  tada quote update 5 --text "Updated quote" --author "Author Name"
  
  # Using short flags
  tada quote update 5 -t "New text" -a "Author" -c "category"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			if !cmd.Flags().Changed("text") &&
				!cmd.Flags().Changed("author") &&
				!cmd.Flags().Changed("category") {
				quote.PrintError(cmd, fmt.Errorf("at least one flag (--text, --author, or --category) must be provided"))
				return nil
			}

			db, cleanup, err := quote.GetDB(cmd)
			if err != nil {
				return nil
			}
			defer cleanup()

			currentQuote, err := db.Get(id)
			if err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			text := currentQuote.Text
			author := currentQuote.Author
			category := currentQuote.Category

			if cmd.Flags().Changed("text") {
				text, _ = cmd.Flags().GetString("text")
				if err := quote.ValidateQuoteText(text); err != nil {
					quote.PrintError(cmd, err)
					return nil
				}
			}

			if cmd.Flags().Changed("author") {
				author, _ = cmd.Flags().GetString("author")
				if err := quote.ValidateAuthor(author); err != nil {
					quote.PrintError(cmd, err)
					return nil
				}
			}

			if cmd.Flags().Changed("category") {
				category, _ = cmd.Flags().GetString("category")
				if err := quote.ValidateCategory(category); err != nil {
					quote.PrintError(cmd, err)
					return nil
				}
			}

			if err := db.Update(id, text, author, category); err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			updatedQuote, err := db.Get(id)
			if err != nil {
				quote.PrintError(cmd, err)
				return nil
			}

			quote.PrintSuccess(cmd, "Updated quote:")
			quote.PrintQuote(cmd, updatedQuote)
			return nil
		},
	}

	cmd.Flags().StringP("text", "t", "", "Update quote text")
	cmd.Flags().StringP("author", "a", "", "Update author")
	cmd.Flags().StringP("category", "c", "", "Update category")

	return cmd
}
