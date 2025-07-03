package quote

import (
	"github.com/spf13/cobra"
	"math/rand"
)

// Quote when the quote command is used it should return randomly a quote from the list
// should this quote be in a json so that we can manage them edit/update/delete
var Quote = []string{
	"🚿 Shower thoughts only work when you put in the work.",
	"⏰ Take the time it takes so it takes less time.",
	"🌱 Don't go out and try to find the quotes, they should come to you.",
	"🏆 Nobody cares how hard you worked, only the results.",
	"💡 If you want everything to be familiar, you will never learn anything new.",
	"✨ Do or do not. There is no try.",
	"⚡ Inspiration is perishable. Act on it immediately",
	"🚀 The way to get started is to quit talking and begin doing. - Walt Disney",
	"🌅 Don't let yesterday take up too much of today. - Will Rogers",
	"💪 You learn more from failure than from success. Don't let it stop you. Failure builds character. - Unknown",
	"🥊 It's not whether you get knocked down, it's whether you get up. - Vince Lombardi",
	"🎯 If you are working on something that you really care about, you don't have to be pushed. The vision pulls you. - Steve Jobs",
	"🎓 Knowing is not enough; we must apply. Wishing is not enough; we must do. - Johann Wolfgang von Goethe",
	"🧠 Whether you think you can or you think you can't, you're right. - Henry Ford",
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote",
		Short: "Show a motivational quote",
		Long:  "Display a random motivational quote to inspire productivity",
		RunE: func(cmd *cobra.Command, args []string) error {
			randomQuote := Quote[rand.Intn(len(Quote))]

			cmd.Printf("\n %s\n", randomQuote)
			return nil
		},
	}

	return cmd
}
