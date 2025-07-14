package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/negadras/tada/internal/tui/models"
	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

// Screen represents different screens in the TUI
type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenTodos
	ScreenQuotes
	ScreenHelp
)

// App represents the main TUI application
type App struct {
	currentScreen Screen
	width         int
	height        int
	keymap        utils.KeyMap
	styles        *styles.Styles
	help          help.Model
	showHelp      bool

	// Screen models
	dashboard *models.Dashboard
	todos     *models.TodoManager
	quotes    *models.QuoteManager

	// Navigation
	screens     []Screen
	screenIndex int
}

// NewApp creates a new TUI application
func NewApp() *App {
	keymap := utils.DefaultKeyMap()
	styles := styles.DefaultStyles()
	help := help.New()

	return &App{
		currentScreen: ScreenDashboard,
		keymap:        keymap,
		styles:        styles,
		help:          help,
		showHelp:      false,
		screens:       []Screen{ScreenDashboard, ScreenTodos, ScreenQuotes},
		screenIndex:   0,
	}
}

// Init initializes the TUI application
func (a *App) Init() tea.Cmd {
	// Initialize screen models
	a.dashboard = models.NewDashboard(a.styles, a.keymap)
	a.todos = models.NewTodoManager(a.styles, a.keymap)
	a.quotes = models.NewQuoteManager(a.styles, a.keymap)

	// Initialize all models
	var cmds []tea.Cmd
	cmds = append(cmds, a.dashboard.Init())
	cmds = append(cmds, a.todos.Init())
	cmds = append(cmds, a.quotes.Init())

	return tea.Batch(cmds...)
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.help.Width = msg.Width

		// Update all screen models with new size
		a.dashboard.SetSize(msg.Width, msg.Height)
		a.todos.SetSize(msg.Width, msg.Height)
		a.quotes.SetSize(msg.Width, msg.Height)

		return a, nil

	case tea.KeyMsg:
		// Global key bindings
		switch {
		case key.Matches(msg, a.keymap.Quit):
			return a, tea.Quit

		case key.Matches(msg, a.keymap.Help):
			a.showHelp = !a.showHelp
			return a, nil

		case key.Matches(msg, a.keymap.Tab):
			a.nextScreen()
			return a, nil

		case key.Matches(msg, a.keymap.ShiftTab):
			a.prevScreen()
			return a, nil

		case key.Matches(msg, a.keymap.Escape):
			if a.showHelp {
				a.showHelp = false
				return a, nil
			}
			// Handle escape in current screen
		}

		// Handle screen-specific key bindings
		if a.showHelp {
			return a, nil
		}

		return a.updateCurrentScreen(msg)
	}

	// Forward other messages to current screen
	return a.updateCurrentScreen(msg)
}

// updateCurrentScreen updates the current screen with the given message
func (a *App) updateCurrentScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Always forward loading messages to their respective managers
	if _, ok := msg.(models.TodosLoadedMsg); ok {
		if a.todos != nil {
			a.todos, cmd = a.todos.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}
	if _, ok := msg.(models.TodoErrorMsg); ok {
		if a.todos != nil {
			a.todos, cmd = a.todos.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}
	if _, ok := msg.(models.QuotesLoadedMsg); ok {
		if a.quotes != nil {
			a.quotes, cmd = a.quotes.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}
	if _, ok := msg.(models.QuoteErrorMsg); ok {
		if a.quotes != nil {
			a.quotes, cmd = a.quotes.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}
	if _, ok := msg.(models.StatsLoadedMsg); ok {
		if a.dashboard != nil {
			a.dashboard, cmd = a.dashboard.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}
	if _, ok := msg.(models.StatsErrorMsg); ok {
		if a.dashboard != nil {
			a.dashboard, cmd = a.dashboard.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	// Handle current screen updates
	switch a.currentScreen {
	case ScreenDashboard:
		a.dashboard, cmd = a.dashboard.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Handle navigation from dashboard
		if navMsg, ok := msg.(models.NavigationMsg); ok {
			switch navMsg.Screen {
			case "todos":
				a.currentScreen = ScreenTodos
				a.screenIndex = 1
			case "quotes":
				a.currentScreen = ScreenQuotes
				a.screenIndex = 2
			}
		}

	case ScreenTodos:
		a.todos, cmd = a.todos.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case ScreenQuotes:
		a.quotes, cmd = a.quotes.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return a, tea.Batch(cmds...)
}

// View renders the TUI application
func (a *App) View() string {
	if a.showHelp {
		return a.renderHelp()
	}

	var content string

	switch a.currentScreen {
	case ScreenDashboard:
		content = a.dashboard.View()
	case ScreenTodos:
		content = a.todos.View()
	case ScreenQuotes:
		content = a.quotes.View()
	}

	return a.renderWithChrome(content)
}

// renderWithChrome renders the content with navigation and status bar
func (a *App) renderWithChrome(content string) string {
	// Create status bar
	statusBar := a.renderStatusBar()

	// Create navigation breadcrumb
	breadcrumb := a.renderBreadcrumb()

	// Calculate available height for content
	contentHeight := a.height - 4 // Reserve space for status bar, breadcrumb, and margins

	// Render main content
	mainContent := lipgloss.NewStyle().
		Width(a.width).
		Height(contentHeight).
		Render(content)

	// Combine all elements
	return lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		breadcrumb,
		mainContent,
		a.renderFooter(),
	)
}

// renderStatusBar renders the top status bar
func (a *App) renderStatusBar() string {
	title := a.styles.Title.Render("📋 Tada")

	var screenName string
	switch a.currentScreen {
	case ScreenDashboard:
		screenName = "Dashboard"
	case ScreenTodos:
		screenName = "Todos"
	case ScreenQuotes:
		screenName = "Quotes"
	}

	subtitle := a.styles.Subtitle.Render(screenName)

	left := lipgloss.JoinHorizontal(lipgloss.Left, title, subtitle)
	right := a.styles.Muted.Render("? for help")

	return a.styles.StatusBar.
		Width(a.width).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			left,
			strings.Repeat(" ", utils.Max(0, a.width-lipgloss.Width(left)-lipgloss.Width(right))),
			right,
		))
}

// renderBreadcrumb renders the navigation breadcrumb
func (a *App) renderBreadcrumb() string {
	var items []string

	for i, screen := range a.screens {
		var name string
		switch screen {
		case ScreenDashboard:
			name = "Dashboard"
		case ScreenTodos:
			name = "Todos"
		case ScreenQuotes:
			name = "Quotes"
		}

		if i == a.screenIndex {
			items = append(items, a.styles.Highlight.Render(name))
		} else {
			items = append(items, a.styles.Muted.Render(name))
		}
	}

	breadcrumb := lipgloss.JoinHorizontal(lipgloss.Left, strings.Join(items, " • "))
	return a.styles.Content.Render(breadcrumb)
}

// renderFooter renders the bottom footer with key bindings
func (a *App) renderFooter() string {
	helpText := a.styles.Help.Render("tab: next • shift+tab: prev • q: quit • ?: help")
	return a.styles.Footer.
		Width(a.width).
		Render(helpText)
}

// renderHelp renders the help screen
func (a *App) renderHelp() string {
	helpContent := lipgloss.NewStyle().
		Width(a.width-4).
		Height(a.height-4).
		Padding(1, 2).
		Render(a.help.View(a.keymap))

	return a.styles.Panel.
		Width(a.width).
		Height(a.height).
		Render(helpContent)
}

// nextScreen navigates to the next screen
func (a *App) nextScreen() {
	a.screenIndex = (a.screenIndex + 1) % len(a.screens)
	a.currentScreen = a.screens[a.screenIndex]
}

// prevScreen navigates to the previous screen
func (a *App) prevScreen() {
	a.screenIndex = (a.screenIndex - 1 + len(a.screens)) % len(a.screens)
	a.currentScreen = a.screens[a.screenIndex]
}

// Run starts the TUI application
func Run() error {
	return RunWithScreen("")
}

// RunWithScreen starts the TUI application at a specific screen
func RunWithScreen(screenName string) error {
	app := NewApp()

	// Set initial screen based on parameter
	switch screenName {
	case "todos":
		app.currentScreen = ScreenTodos
		app.screenIndex = 1
	case "quotes":
		app.currentScreen = ScreenQuotes
		app.screenIndex = 2
	case "dashboard", "":
		// Default to dashboard
		app.currentScreen = ScreenDashboard
		app.screenIndex = 0
	default:
		// Invalid screen name, default to dashboard
		app.currentScreen = ScreenDashboard
		app.screenIndex = 0
	}

	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
