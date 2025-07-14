package quote

import (
	"strings"
)

// MigrateHardcodedQuotes migrates the hardcoded quotes to the database
func (db *DB) MigrateHardcodedQuotes() error {
	quotes, err := db.List(nil, nil)
	if err != nil {
		return err
	}

	if len(quotes) > 0 {
		return nil
	}

	hardcodedQuotes := []string{
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

	for _, quoteText := range hardcodedQuotes {
		text, author := parseQuote(quoteText)
		_, err := db.Create(text, author, "motivational")
		if err != nil {
			return err
		}
	}

	return nil
}

// parseQuote extracts the text and author from a quote string
func parseQuote(quoteText string) (text, author string) {
	if strings.Contains(quoteText, " - ") {
		parts := strings.Split(quoteText, " - ")
		if len(parts) >= 2 {
			text = strings.TrimSpace(parts[0])
			author = strings.TrimSpace(parts[1])
			return text, author
		}
	}

	return quoteText, ""
}
