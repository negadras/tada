package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/negadras/tada/internal/quote"
	"github.com/negadras/tada/internal/tui/components"
	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

// Table column width constants
const (
	quoteIDWidth       = 6
	quoteAuthorWidth   = 20
	quoteCategoryWidth = 15
	quoteAgeWidth      = 10
	quoteTextWidth     = 60
	quoteTableHeight   = 15
	quoteTableMargin   = 4
	quoteTableReserved = 10
)

// QuotesLoadedMsg is sent when quotes are loaded from the database
type QuotesLoadedMsg struct {
	Quotes []*quote.Quote
}

// QuoteErrorMsg is sent when there's an error loading quotes
type QuoteErrorMsg struct {
	Error error
}

// QuoteManager represents the quote management screen
type QuoteManager struct {
	BaseModel
	styles            *styles.Styles
	keymap            utils.KeyMap
	table             table.Model
	quotes            []*quote.Quote
	db                *quote.DB
	loading           bool
	errorMessage      string
	categoryFilter    *string
	currentQuote      *quote.Quote
	showQuote         bool
	addForm           *components.Form
	showAddForm       bool
	editForm          *components.Form
	showEditForm      bool
	editingQuote      *quote.Quote
	showDeleteConfirm bool
	quoteToDelete     *quote.Quote
}

// NewQuoteManager creates a new quote manager model
func NewQuoteManager(styles *styles.Styles, keymap utils.KeyMap) *QuoteManager {
	// Create table columns
	columns := []table.Column{
		{Title: "ID", Width: quoteIDWidth},
		{Title: "Author", Width: quoteAuthorWidth},
		{Title: "Category", Width: quoteCategoryWidth},
		{Title: "Age", Width: quoteAgeWidth},
		{Title: "Quote", Width: quoteTextWidth},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(quoteTableHeight),
	)

	// Apply table styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = styles.TableRowActive
	t.SetStyles(s)

	// Create add form
	addForm := components.NewForm(styles, keymap, "Add New Quote")
	addForm.AddField("Quote", "Enter the quote text", true)
	addForm.AddField("Author", "Author name (optional)", false)
	addForm.AddField("Category", "Category (optional)", false)

	// Create edit form
	editForm := components.NewForm(styles, keymap, "Edit Quote")
	editForm.AddField("Quote", "Enter the quote text", true)
	editForm.AddField("Author", "Author name (optional)", false)
	editForm.AddField("Category", "Category (optional)", false)

	return &QuoteManager{
		styles:   styles,
		keymap:   keymap,
		table:    t,
		quotes:   []*quote.Quote{},
		loading:  true,
		addForm:  addForm,
		editForm: editForm,
	}
}

// Init initializes the quote manager
func (q *QuoteManager) Init() tea.Cmd {
	return q.loadQuotes()
}

// Update handles messages for the quote manager
func (q *QuoteManager) Update(msg tea.Msg) (*QuoteManager, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case QuotesLoadedMsg:
		q.quotes = msg.Quotes
		q.loading = false
		q.errorMessage = ""
		q.updateTable()
		return q, nil

	case QuoteErrorMsg:
		q.loading = false
		q.errorMessage = msg.Error.Error()
		return q, nil

	case tea.KeyMsg:
		// Handle add form mode
		if q.showAddForm {
			q.addForm, cmd = q.addForm.Update(msg)

			if q.addForm.IsSubmitted() {
				// Create new quote
				quoteText := q.addForm.GetValue(0)
				author := q.addForm.GetValue(1)
				category := q.addForm.GetValue(2)

				q.showAddForm = false
				q.addForm.Reset()
				return q, q.createQuote(quoteText, author, category)
			}

			if q.addForm.IsCancelled() {
				q.showAddForm = false
				q.addForm.Reset()
				return q, nil
			}

			return q, cmd
		}

		// Handle edit form mode
		if q.showEditForm {
			q.editForm, cmd = q.editForm.Update(msg)

			if q.editForm.IsSubmitted() {
				// Update quote
				quoteText := q.editForm.GetValue(0)
				author := q.editForm.GetValue(1)
				category := q.editForm.GetValue(2)

				q.showEditForm = false
				q.editForm.Reset()
				quoteID := q.editingQuote.ID
				q.editingQuote = nil
				return q, q.updateQuote(quoteID, quoteText, author, category)
			}

			if q.editForm.IsCancelled() {
				q.showEditForm = false
				q.editForm.Reset()
				q.editingQuote = nil
				return q, nil
			}

			return q, cmd
		}

		// Handle delete confirmation mode
		if q.showDeleteConfirm {
			switch {
			case key.Matches(msg, q.keymap.Enter):
				// Confirm deletion
				q.showDeleteConfirm = false
				quoteToDelete := q.quoteToDelete
				q.quoteToDelete = nil
				return q, q.confirmDeleteQuote(quoteToDelete)
			case key.Matches(msg, q.keymap.Escape):
				// Cancel deletion
				q.showDeleteConfirm = false
				q.quoteToDelete = nil
				return q, nil
			}
			return q, nil
		}

		// Handle quote view mode
		if q.showQuote {
			switch {
			case key.Matches(msg, q.keymap.Escape), key.Matches(msg, q.keymap.Back):
				q.showQuote = false
				return q, nil
			case key.Matches(msg, q.keymap.Enter):
				return q, q.getRandomQuote()
			}
			return q, nil
		}

		// Handle table mode
		switch {
		case key.Matches(msg, q.keymap.Add):
			q.showAddForm = true
			q.addForm.SetSize(q.width, q.height)
			return q, q.addForm.Init()

		case key.Matches(msg, q.keymap.Edit):
			return q, q.openEditForm()

		case key.Matches(msg, q.keymap.Delete):
			return q, q.showDeleteConfirmation()

		case key.Matches(msg, q.keymap.Filter):
			return q, q.cycleCategoryFilter()

		case key.Matches(msg, q.keymap.Enter):
			return q, q.viewQuote()

		case key.Matches(msg, q.keymap.Space):
			return q, q.getRandomQuote()

		default:
			// Handle table navigation
			q.table, cmd = q.table.Update(msg)
			return q, cmd
		}
	}

	return q, nil
}

// View renders the quote manager
func (q *QuoteManager) View() string {
	// Show add form if active
	if q.showAddForm {
		return q.addForm.View()
	}

	// Show edit form if active
	if q.showEditForm {
		return q.editForm.View()
	}

	// Show delete confirmation if active
	if q.showDeleteConfirm {
		return q.renderDeleteConfirmation()
	}

	if q.showQuote && q.currentQuote != nil {
		return q.renderQuoteView()
	}

	var content strings.Builder

	// Title and stats
	title := q.styles.Title.Render("💬 Quote Collection")
	stats := q.renderStats()
	header := lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", stats)
	content.WriteString(header)
	content.WriteString("\n\n")

	// Category filter info
	if q.categoryFilter != nil {
		filterText := q.styles.Info.Render(fmt.Sprintf("Filter: %s", *q.categoryFilter))
		content.WriteString(filterText)
		content.WriteString("\n")
	}

	// Error message
	if q.errorMessage != "" {
		errorText := q.styles.Error.Render(fmt.Sprintf("Error: %s", q.errorMessage))
		content.WriteString(errorText)
		content.WriteString("\n")
	}

	// Loading state
	if q.loading {
		loadingText := q.styles.Info.Render("Loading quotes...")
		content.WriteString(loadingText)
		content.WriteString("\n")
	}

	// Table
	if !q.loading && q.errorMessage == "" {
		if len(q.quotes) == 0 {
			emptyText := q.styles.Muted.Render("No quotes found. Press 'a' to add one!")
			content.WriteString(emptyText)
		} else {
			content.WriteString(q.table.View())
		}
	}

	content.WriteString("\n")

	// Instructions
	instructions := q.renderInstructions()
	content.WriteString(instructions)

	return q.styles.Content.Render(content.String())
}

// renderQuoteView renders the full quote view
func (q *QuoteManager) renderQuoteView() string {
	var content strings.Builder

	// Title
	title := q.styles.Title.Render("💬 Quote View")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Quote content
	quoteCard := q.renderQuoteCard(q.currentQuote)
	content.WriteString(quoteCard)

	// Instructions
	instructions := q.styles.Help.Render("enter: random quote • esc: back to list")
	content.WriteString("\n\n" + instructions)

	return q.styles.Content.Render(content.String())
}

// renderQuoteCard renders a single quote as a card
func (q *QuoteManager) renderQuoteCard(quote *quote.Quote) string {
	var content strings.Builder

	// Quote text
	quoteText := q.styles.Title.Render(fmt.Sprintf("\"%s\"", quote.Text))
	content.WriteString(quoteText)
	content.WriteString("\n")

	// Author
	if quote.Author != "" {
		authorText := q.styles.Subtitle.Render(fmt.Sprintf("— %s", quote.Author))
		content.WriteString(authorText)
		content.WriteString("\n")
	}

	// Metadata
	var metadata []string
	if quote.Category != "" {
		metadata = append(metadata, fmt.Sprintf("Category: %s", quote.Category))
	}
	metadata = append(metadata, fmt.Sprintf("ID: #%d", quote.ID))
	metadata = append(metadata, fmt.Sprintf("Added: %s ago", utils.FormatDuration(quote.Age())))

	metadataText := q.styles.Muted.Render(strings.Join(metadata, " • "))
	content.WriteString(metadataText)

	return q.styles.Card.
		Width(80).
		Height(10).
		Render(content.String())
}

// renderStats renders the quote statistics
func (q *QuoteManager) renderStats() string {
	if len(q.quotes) == 0 {
		return q.styles.Muted.Render("No quotes")
	}

	// Count quotes by category
	categoryCount := make(map[string]int)
	for _, quote := range q.quotes {
		if quote.Category != "" {
			categoryCount[quote.Category]++
		}
	}

	var categories []string
	for category, count := range categoryCount {
		categories = append(categories, fmt.Sprintf("%s: %d", category, count))
	}

	statsText := fmt.Sprintf("%d total", len(q.quotes))
	if len(categories) > 0 {
		statsText += " • " + strings.Join(categories[:utils.Min(2, len(categories))], " • ")
	}

	return q.styles.Success.Render(statsText)
}

// renderInstructions renders the keyboard instructions
func (q *QuoteManager) renderInstructions() string {
	instructions := []string{
		"a: add quote",
		"e: edit",
		"d: delete",
		"enter: view",
		"space: random",
		"f: filter",
		"esc: back",
	}

	return q.styles.Help.Render(strings.Join(instructions, " • "))
}

// updateTable updates the table with current quotes.
// Formats each quote's data to fit within the designated column widths.
// Shows "Unknown" for empty authors and "None" for empty categories.
func (q *QuoteManager) updateTable() {
	rows := make([]table.Row, len(q.quotes))
	for i, quote := range q.quotes {
		// Handle empty author
		author := quote.Author
		if author == "" {
			author = "Unknown"
		}

		// Handle empty category
		category := quote.Category
		if category == "" {
			category = "None"
		}

		rows[i] = table.Row{
			fmt.Sprintf("#%d", quote.ID),
			author,
			category,
			utils.FormatDuration(quote.Age()),
			quote.Text,
		}
	}

	q.table.SetRows(rows)
}

// loadQuotes loads quotes from the database
func (q *QuoteManager) loadQuotes() tea.Cmd {
	return func() tea.Msg {
		// Get database path
		dbPath, err := quote.GetDatabasePath()
		if err != nil {
			return QuoteErrorMsg{Error: err}
		}

		// Create database connection
		db, err := quote.NewDB(dbPath)
		if err != nil {
			return QuoteErrorMsg{Error: err}
		}

		// Store database connection for future operations
		q.db = db

		// Migrate hardcoded quotes if needed
		if err := db.MigrateHardcodedQuotes(); err != nil {
			return QuoteErrorMsg{Error: err}
		}

		// Load quotes
		quotes, err := db.List(nil, q.categoryFilter)
		if err != nil {
			return QuoteErrorMsg{Error: err}
		}

		return QuotesLoadedMsg{Quotes: quotes}
	}
}

// viewQuote shows the selected quote in detail
func (q *QuoteManager) viewQuote() tea.Cmd {
	if len(q.quotes) == 0 {
		return nil
	}

	selectedIndex := q.table.Cursor()
	if selectedIndex >= len(q.quotes) {
		return nil
	}

	q.currentQuote = q.quotes[selectedIndex]
	q.showQuote = true

	return nil
}

// getRandomQuote gets a random quote
func (q *QuoteManager) getRandomQuote() tea.Cmd {
	if q.db == nil {
		return nil
	}

	return func() tea.Msg {
		randomQuote, err := q.db.GetRandom()
		if err != nil {
			return QuoteErrorMsg{Error: err}
		}

		q.currentQuote = randomQuote
		q.showQuote = true

		return nil
	}
}

// cycleCategoryFilter cycles through category filters
func (q *QuoteManager) cycleCategoryFilter() tea.Cmd {
	return func() tea.Msg {
		// Get unique categories
		categories := make(map[string]bool)
		for _, quote := range q.quotes {
			if quote.Category != "" {
				categories[quote.Category] = true
			}
		}

		// Convert to slice
		var categoryList []string
		for category := range categories {
			categoryList = append(categoryList, category)
		}

		// Cycle through filters
		if q.categoryFilter == nil && len(categoryList) > 0 {
			// Set first category filter
			q.categoryFilter = &categoryList[0]
		} else if q.categoryFilter != nil {
			// Find current category and move to next
			found := false
			for i, category := range categoryList {
				if category == *q.categoryFilter {
					if i < len(categoryList)-1 {
						q.categoryFilter = &categoryList[i+1]
					} else {
						q.categoryFilter = nil // Reset to all
					}
					found = true
					break
				}
			}
			if !found {
				q.categoryFilter = nil
			}
		}

		// Reload quotes with new filter
		quotes, err := q.db.List(nil, q.categoryFilter)
		if err != nil {
			return QuoteErrorMsg{Error: err}
		}

		return QuotesLoadedMsg{Quotes: quotes}
	}
}

// createQuote creates a new quote with the given text, author, and category.
// Returns a command that will send either QuotesLoadedMsg or QuoteErrorMsg.
func (q *QuoteManager) createQuote(text, author, category string) tea.Cmd {
	if q.db == nil {
		return nil
	}

	return func() tea.Msg {
		// Create new quote
		_, err := q.db.Create(text, author, category)
		if err != nil {
			return QuoteErrorMsg{Error: err}
		}

		// Reload quotes
		quotes, err := q.db.List(nil, q.categoryFilter)
		if err != nil {
			return QuoteErrorMsg{Error: err}
		}

		return QuotesLoadedMsg{Quotes: quotes}
	}
}

// openEditForm opens the edit form for the currently selected quote.
// Pre-populates the form with the quote's current text, author, and category.
// Returns nil if no quotes exist or selection is invalid.
func (q *QuoteManager) openEditForm() tea.Cmd {
	if len(q.quotes) == 0 {
		return nil
	}

	selectedIndex := q.table.Cursor()
	if selectedIndex >= len(q.quotes) {
		return nil
	}

	selectedQuote := q.quotes[selectedIndex]
	q.editingQuote = selectedQuote

	// Pre-populate form with current values
	q.editForm.Reset()
	q.editForm.SetValue(0, selectedQuote.Text)
	q.editForm.SetValue(1, selectedQuote.Author)
	q.editForm.SetValue(2, selectedQuote.Category)

	q.showEditForm = true
	q.editForm.SetSize(q.width, q.height)
	return q.editForm.Init()
}

// updateQuote updates an existing quote's text, author, and category in the database.
// Returns a command that will send either QuotesLoadedMsg or QuoteErrorMsg.
func (q *QuoteManager) updateQuote(id int, text, author, category string) tea.Cmd {
	if q.db == nil {
		return nil
	}

	return func() tea.Msg {
		// Update quote in database
		err := q.db.Update(id, text, author, category)
		if err != nil {
			return QuoteErrorMsg{Error: fmt.Errorf("failed to update quote: %w", err)}
		}

		// Reload quotes
		quotes, err := q.db.List(nil, q.categoryFilter)
		if err != nil {
			return QuoteErrorMsg{Error: fmt.Errorf("failed to reload quotes: %w", err)}
		}

		return QuotesLoadedMsg{Quotes: quotes}
	}
}

// showDeleteConfirmation shows the delete confirmation dialog for the selected quote.
// Returns nil if no quotes exist or selection is invalid.
func (q *QuoteManager) showDeleteConfirmation() tea.Cmd {
	if len(q.quotes) == 0 {
		return nil
	}

	selectedIndex := q.table.Cursor()
	if selectedIndex >= len(q.quotes) {
		return nil
	}

	q.quoteToDelete = q.quotes[selectedIndex]
	q.showDeleteConfirm = true
	return nil
}

// confirmDeleteQuote actually deletes the quote after confirmation.
// Returns a command that will send either QuotesLoadedMsg or QuoteErrorMsg.
func (q *QuoteManager) confirmDeleteQuote(quoteToDelete *quote.Quote) tea.Cmd {
	if q.db == nil || quoteToDelete == nil {
		return nil
	}

	return func() tea.Msg {
		err := q.db.Delete(quoteToDelete.ID)
		if err != nil {
			return QuoteErrorMsg{Error: fmt.Errorf("failed to delete quote: %w", err)}
		}

		// Reload quotes
		quotes, err := q.db.List(nil, q.categoryFilter)
		if err != nil {
			return QuoteErrorMsg{Error: fmt.Errorf("failed to reload quotes: %w", err)}
		}

		return QuotesLoadedMsg{Quotes: quotes}
	}
}

// renderDeleteConfirmation renders the delete confirmation dialog.
func (q *QuoteManager) renderDeleteConfirmation() string {
	if q.quoteToDelete == nil {
		return ""
	}

	var content strings.Builder

	// Title
	title := q.styles.Title.Render("🗑️ Delete Quote")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Warning message
	warningMsg := q.styles.Error.Render("Are you sure you want to delete this quote?")
	content.WriteString(warningMsg)
	content.WriteString("\n\n")

	// Quote details
	author := q.quoteToDelete.Author
	if author == "" {
		author = "Unknown"
	}
	category := q.quoteToDelete.Category
	if category == "" {
		category = "None"
	}

	quoteDetails := fmt.Sprintf("ID: #%d\nAuthor: %s\nCategory: %s\n\nQuote: \"%s\"",
		q.quoteToDelete.ID,
		author,
		category,
		q.quoteToDelete.Text,
	)
	content.WriteString(q.styles.Muted.Render(quoteDetails))
	content.WriteString("\n\n")

	// Instructions
	instructions := q.styles.Help.Render("enter: confirm deletion • esc: cancel")
	content.WriteString(instructions)

	// Wrap in panel
	return q.styles.Panel.
		Width(q.width - 4).
		Height(q.height - 4).
		Render(content.String())
}

// SetSize sets the size of the quote manager and updates child components.
func (q *QuoteManager) SetSize(width, height int) {
	q.BaseModel.SetSize(width, height)
	q.table.SetWidth(width - quoteTableMargin)
	q.table.SetHeight(height - quoteTableReserved)
	if q.addForm != nil {
		q.addForm.SetSize(width, height)
	}
	if q.editForm != nil {
		q.editForm.SetSize(width, height)
	}
}
