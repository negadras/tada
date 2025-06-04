package cmd

import (
	"github.com/spf13/cobra"
	"math/rand"
)

// Quote when the quote command is used it should return randomly a quote from the list
// should this quote be in a json so that we can manage them edit/update/delete
var Quote = []string{
	"Shower thoughts only work when you put in the work.",
	"Take the time it takes so it takes less time.",
	"Don't go out and try to find the quotes, they should come to you.",
	"Nobody cares how hard you worked, only the results.",
	"If you want everything to be familiar, you will never learn anything new.",
	"Do or do not. There is no try.",
}

func newQuoteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote",
		Short: "Tada random quote from our list",
		RunE: func(cmd *cobra.Command, args []string) error {
			randomQuote := Quote[rand.Intn(len(Quote))]

			cmd.Printf("\n 💬 %s\n", randomQuote)
			return nil
		},
	}

	return cmd
}
