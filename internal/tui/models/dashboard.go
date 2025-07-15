package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/negadras/tada/internal/quote"
	"github.com/negadras/tada/internal/todo"
	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalTodos     int
	CompletedTodos int
	TotalQuotes    int
	TodayCompleted int
	CompletionRate float64
	Loading        bool
	Error          string
}

// Dashboard represents the main dashboard screen
type Dashboard struct {
	BaseModel
	styles       *styles.Styles
	keymap       utils.KeyMap
	selectedItem int
	menuItems    []MenuItem
	stats        DashboardStats
}

// MenuItem represents a menu item in the dashboard
type MenuItem struct {
	Icon        string
	Title       string
	Description string
	Action      string
}

// StatsLoadedMsg is sent when dashboard statistics are loaded
type StatsLoadedMsg struct {
	Stats DashboardStats
}

// StatsErrorMsg is sent when there's an error loading statistics
type StatsErrorMsg struct {
	Error error
}

// NewDashboard creates a new dashboard model
func NewDashboard(styles *styles.Styles, keymap utils.KeyMap) *Dashboard {
	menuItems := []MenuItem{
		{
			Icon:        "📝",
			Title:       "Todo Management",
			Description: "View, add, edit and manage your todos",
			Action:      "todos",
		},
		{
			Icon:        "💬",
			Title:       "Quote Collection",
			Description: "Browse and manage your motivational quotes",
			Action:      "quotes",
		},
		{
			Icon:        "📊",
			Title:       "Statistics",
			Description: "View your productivity statistics",
			Action:      "stats",
		},
		{
			Icon:        "⚙️",
			Title:       "Settings",
			Description: "Configure your preferences",
			Action:      "settings",
		},
	}

	return &Dashboard{
		styles:       styles,
		keymap:       keymap,
		selectedItem: 0,
		menuItems:    menuItems,
		stats: DashboardStats{
			Loading: true,
		},
	}
}

// Init initializes the dashboard and loads statistics
func (d *Dashboard) Init() tea.Cmd {
	return d.loadStats()
}

// Update handles messages for the dashboard
func (d *Dashboard) Update(msg tea.Msg) (*Dashboard, tea.Cmd) {
	switch msg := msg.(type) {
	case StatsLoadedMsg:
		d.stats = msg.Stats
		d.stats.Loading = false
		return d, nil

	case StatsErrorMsg:
		d.stats.Loading = false
		d.stats.Error = msg.Error.Error()
		return d, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keymap.Up):
			d.selectedItem = utils.Max(0, d.selectedItem-1)
			return d, nil

		case key.Matches(msg, d.keymap.Down):
			d.selectedItem = utils.Min(len(d.menuItems)-1, d.selectedItem+1)
			return d, nil

		case key.Matches(msg, d.keymap.Enter):
			action := d.menuItems[d.selectedItem].Action
			if action == "todos" || action == "quotes" {
				return d, func() tea.Msg {
					return NavigationMsg{Screen: action}
				}
			}
			return d, nil
		}
	}

	return d, nil
}

// View renders the dashboard
func (d *Dashboard) View() string {
	var content strings.Builder

	welcomeText := d.styles.Title.Render("Welcome to Tada!")
	subtitle := d.styles.Subtitle.Render("Your productivity companion")

	content.WriteString(lipgloss.JoinVertical(
		lipgloss.Center,
		welcomeText,
		subtitle,
		"",
	))

	// Quick stats (placeholder)
	stats := d.renderStats()
	content.WriteString(stats)
	content.WriteString("\n\n")

	// Navigation menu
	menu := d.renderMenu()
	content.WriteString(menu)

	// Instructions
	instructions := d.styles.Help.Render("↑/↓: Navigate • Enter: Select • Tab: Switch sections")
	content.WriteString("\n\n" + instructions)

	return d.styles.Content.Render(content.String())
}

// renderStats renders the quick statistics section with real database data
func (d *Dashboard) renderStats() string {
	var stats strings.Builder

	// Show loading state
	if d.stats.Loading {
		loadingCard := d.renderStatCard("⏳ Loading", "Fetching stats...", "Please wait")
		stats.WriteString(loadingCard)
		return stats.String()
	}

	// Show error state
	if d.stats.Error != "" {
		errorCard := d.renderStatCard("❌ Error", "Failed to load", d.stats.Error)
		stats.WriteString(errorCard)
		return stats.String()
	}

	// Create stat cards with real data
	todoCard := d.renderStatCard("📝 Todos", fmt.Sprintf("%d total", d.stats.TotalTodos), fmt.Sprintf("%.1f%% complete", d.stats.CompletionRate))
	quoteCard := d.renderStatCard("💬 Quotes", fmt.Sprintf("%d total", d.stats.TotalQuotes), "Collection")
	productivityCard := d.renderStatCard("📊 Today", fmt.Sprintf("%d completed", d.stats.TodayCompleted), "Great progress!")

	stats.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		todoCard,
		quoteCard,
		productivityCard,
	))

	return stats.String()
}

// renderStatCard renders a single statistics card
func (d *Dashboard) renderStatCard(title, value, description string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		d.styles.Title.Render(title),
		d.styles.Success.Render(value),
		d.styles.Muted.Render(description),
	)

	return d.styles.Card.
		Width(20).
		Height(6).
		Render(content)
}

// renderMenu renders the navigation menu
func (d *Dashboard) renderMenu() string {
	var menu strings.Builder

	menu.WriteString(d.styles.Subtitle.Render("Navigation"))
	menu.WriteString("\n\n")

	for i, item := range d.menuItems {
		var style lipgloss.Style
		if i == d.selectedItem {
			style = d.styles.ListItemActive
		} else {
			style = d.styles.ListItem
		}

		itemContent := fmt.Sprintf("%s %s", item.Icon, item.Title)
		if i == d.selectedItem {
			itemContent += fmt.Sprintf("\n  %s", d.styles.Muted.Render(item.Description))
		}

		menu.WriteString(style.Render(itemContent))
		menu.WriteString("\n")
	}

	return menu.String()
}

// loadStats loads dashboard statistics from the database.
// Returns a command that will send either StatsLoadedMsg or StatsErrorMsg.
func (d *Dashboard) loadStats() tea.Cmd {
	return func() tea.Msg {
		stats := DashboardStats{
			Loading: true,
		}

		// Load todo statistics
		todoStats, err := d.loadTodoStats()
		if err != nil {
			return StatsErrorMsg{Error: fmt.Errorf("failed to load todo stats: %w", err)}
		}
		stats.TotalTodos = todoStats.Total
		stats.CompletedTodos = todoStats.Completed
		stats.TodayCompleted = todoStats.TodayCompleted

		// Calculate completion rate
		if stats.TotalTodos > 0 {
			stats.CompletionRate = float64(stats.CompletedTodos) / float64(stats.TotalTodos) * 100
		}

		// Load quote statistics
		quoteStats, err := d.loadQuoteStats()
		if err != nil {
			return StatsErrorMsg{Error: fmt.Errorf("failed to load quote stats: %w", err)}
		}
		stats.TotalQuotes = quoteStats.Total

		return StatsLoadedMsg{Stats: stats}
	}
}

// TodoStats represents todo statistics
type TodoStats struct {
	Total          int
	Completed      int
	TodayCompleted int
}

// QuoteStats represents quote statistics
type QuoteStats struct {
	Total int
}

// loadTodoStats loads todo statistics from the database
func (d *Dashboard) loadTodoStats() (TodoStats, error) {
	// Get database path
	dbPath, err := todo.GetDatabasePath()
	if err != nil {
		return TodoStats{}, err
	}

	// Create database connection
	db, err := todo.NewDB(dbPath)
	if err != nil {
		return TodoStats{}, err
	}

	// Load all todos
	todos, err := db.List(nil, nil)
	if err != nil {
		return TodoStats{}, err
	}

	// Calculate statistics
	var completed, todayCompleted int
	today := time.Now().Truncate(24 * time.Hour)

	for _, todoItem := range todos {
		if todoItem.Status == todo.Done {
			completed++

			// Check if completed today
			if todoItem.CompletedAt != nil && todoItem.CompletedAt.Truncate(24*time.Hour).Equal(today) {
				todayCompleted++
			}
		}
	}

	return TodoStats{
		Total:          len(todos),
		Completed:      completed,
		TodayCompleted: todayCompleted,
	}, nil
}

// loadQuoteStats loads quote statistics from the database
func (d *Dashboard) loadQuoteStats() (QuoteStats, error) {
	// Get database path
	dbPath, err := quote.GetDatabasePath()
	if err != nil {
		return QuoteStats{}, err
	}

	// Create database connection
	db, err := quote.NewDB(dbPath)
	if err != nil {
		return QuoteStats{}, err
	}

	// Load all quotes
	quotes, err := db.List(nil, nil)
	if err != nil {
		return QuoteStats{}, err
	}

	return QuoteStats{
		Total: len(quotes),
	}, nil
}
