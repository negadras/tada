package models

import tea "github.com/charmbracelet/bubbletea"

// NavigationMsg represents a navigation message between screens
type NavigationMsg struct {
	Screen string
	Data   interface{}
}

// ErrorMsg represents an error message
type ErrorMsg struct {
	Error error
}

// SuccessMsg represents a success message
type SuccessMsg struct {
	Message string
}

// Model represents a common interface for all screen models
type Model interface {
	tea.Model
	SetSize(width, height int)
}

// BaseModel provides common functionality for all screen models
type BaseModel struct {
	width  int
	height int
}

// SetSize sets the size of the model
func (m *BaseModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Width returns the width of the model
func (m *BaseModel) Width() int {
	return m.width
}

// Height returns the height of the model
func (m *BaseModel) Height() int {
	return m.height
}
