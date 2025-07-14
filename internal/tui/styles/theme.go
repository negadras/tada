package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color scheme and styling for the TUI
type Theme struct {
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Info       lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
	Border     lipgloss.Color
	Highlight  lipgloss.Color
	Muted      lipgloss.Color
	Accent     lipgloss.Color
}

// DefaultTheme returns the default color theme
func DefaultTheme() Theme {
	return Theme{
		Primary:    lipgloss.Color("39"),  // Blue
		Secondary:  lipgloss.Color("240"), // Gray
		Success:    lipgloss.Color("42"),  // Green
		Warning:    lipgloss.Color("226"), // Yellow
		Error:      lipgloss.Color("196"), // Red
		Info:       lipgloss.Color("75"),  // Light Blue
		Background: lipgloss.Color("235"), // Dark Gray
		Foreground: lipgloss.Color("255"), // White
		Border:     lipgloss.Color("240"), // Medium Gray
		Highlight:  lipgloss.Color("57"),  // Purple
		Muted:      lipgloss.Color("245"), // Light Gray
		Accent:     lipgloss.Color("212"), // Pink
	}
}



