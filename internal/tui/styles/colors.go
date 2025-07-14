package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles contains all the styling definitions for the TUI
type Styles struct {
	Theme Theme

	// Base styles
	Base      lipgloss.Style
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	Header    lipgloss.Style
	Footer    lipgloss.Style
	StatusBar lipgloss.Style

	// Content styles
	Content        lipgloss.Style
	Panel          lipgloss.Style
	Card           lipgloss.Style
	List           lipgloss.Style
	ListItem       lipgloss.Style
	ListItemActive lipgloss.Style

	// Interactive styles
	Button           lipgloss.Style
	ButtonActive     lipgloss.Style
	Input            lipgloss.Style
	InputActive      lipgloss.Style
	InputPlaceholder lipgloss.Style

	// Table styles
	Table          lipgloss.Style
	TableHeader    lipgloss.Style
	TableRow       lipgloss.Style
	TableRowActive lipgloss.Style
	TableCell      lipgloss.Style

	// Status styles
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style
	Muted   lipgloss.Style

	// Special styles
	Border     lipgloss.Style
	Highlight  lipgloss.Style
	Help       lipgloss.Style
	KeyBinding lipgloss.Style
}

// NewStyles creates a new styles instance with the given theme
func NewStyles(theme Theme) *Styles {
	s := &Styles{Theme: theme}

	// Base styles
	s.Base = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Background(theme.Background)

	s.Title = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Padding(0, 1)

	s.Subtitle = lipgloss.NewStyle().
		Foreground(theme.Secondary).
		Italic(true).
		Padding(0, 1)

	s.Header = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Background(theme.Background).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)

	s.Footer = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Background(theme.Background).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Background(theme.Primary).
		Padding(0, 1)

	// Content styles
	s.Content = lipgloss.NewStyle().
		Padding(1, 2)

	s.Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2)

	s.Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2).
		Margin(0, 1)

	s.List = lipgloss.NewStyle().
		Padding(0, 1)

	s.ListItem = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 0)

	s.ListItemActive = lipgloss.NewStyle().
		Foreground(theme.Background).
		Background(theme.Highlight).
		Padding(0, 1).
		Bold(true)

	// Interactive styles
	s.Button = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Background(theme.Secondary).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)

	s.ButtonActive = lipgloss.NewStyle().
		Foreground(theme.Background).
		Background(theme.Primary).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		Bold(true)

	s.Input = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Background(theme.Background).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)

	s.InputActive = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Background(theme.Background).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary)

	s.InputPlaceholder = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Italic(true)

	// Table styles
	s.Table = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)

	s.TableHeader = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Background(theme.Background).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(theme.Border)

	s.TableRow = lipgloss.NewStyle().
		Padding(0, 1)

	s.TableRowActive = lipgloss.NewStyle().
		Foreground(theme.Background).
		Background(theme.Highlight).
		Padding(0, 1).
		Bold(true)

	s.TableCell = lipgloss.NewStyle().
		Padding(0, 1)

	// Status styles
	s.Success = lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	s.Warning = lipgloss.NewStyle().
		Foreground(theme.Warning).
		Bold(true)

	s.Error = lipgloss.NewStyle().
		Foreground(theme.Error).
		Bold(true)

	s.Info = lipgloss.NewStyle().
		Foreground(theme.Info).
		Bold(true)

	s.Muted = lipgloss.NewStyle().
		Foreground(theme.Muted)

	// Special styles
	s.Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)

	s.Highlight = lipgloss.NewStyle().
		Foreground(theme.Background).
		Background(theme.Highlight).
		Bold(true)

	s.Help = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Italic(true).
		Padding(0, 1)

	s.KeyBinding = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	return s
}

// DefaultStyles returns the default styles with the default theme
func DefaultStyles() *Styles {
	return NewStyles(DefaultTheme())
}
