package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

// Form validation constants
const (
	maxFieldLength = 255
	maxInputWidth  = 50
)

// Field represents a form field
type Field struct {
	Label       string
	Input       textinput.Model
	Required    bool
	Placeholder string
}

// Form represents a generic form component
type Form struct {
	styles    *styles.Styles
	keymap    utils.KeyMap
	title     string
	fields    []Field
	focused   int
	submitted bool
	cancelled bool
	width     int
	height    int
}

// NewForm creates a new form with the given styles, keymap, and title.
// The form starts with no fields - use AddField to add input fields.
func NewForm(styles *styles.Styles, keymap utils.KeyMap, title string) *Form {
	return &Form{
		styles: styles,
		keymap: keymap,
		title:  title,
		fields: []Field{},
	}
}

// AddField adds a new input field to the form with the given label, placeholder, and required status.
// The first field added will automatically receive focus.
func (f *Form) AddField(label, placeholder string, required bool) {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = maxFieldLength
	input.Width = maxInputWidth

	field := Field{
		Label:       label,
		Input:       input,
		Required:    required,
		Placeholder: placeholder,
	}

	f.fields = append(f.fields, field)

	if len(f.fields) == 1 {
		f.fields[0].Input.Focus()
	}
}

// SetSize sets the form dimensions and adjusts input field widths accordingly.
func (f *Form) SetSize(width, height int) {
	f.width = width
	f.height = height

	inputWidth := utils.Min(maxInputWidth, width-10)
	for i := range f.fields {
		f.fields[i].Input.Width = inputWidth
	}
}

func (f *Form) Init() tea.Cmd {
	if len(f.fields) > 0 {
		return f.fields[0].Input.Focus()
	}
	return nil
}

// Update handles form updates
func (f *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keymap.Enter):
			// Submit form
			if f.validateForm() {
				f.submitted = true
				return f, nil
			}
			return f, nil

		case key.Matches(msg, f.keymap.Escape):
			f.cancelled = true
			return f, nil

		case key.Matches(msg, f.keymap.Tab), key.Matches(msg, f.keymap.Down):
			if len(f.fields) > 0 {
				f.fields[f.focused].Input.Blur()
				f.focused = (f.focused + 1) % len(f.fields)
				cmds = append(cmds, f.fields[f.focused].Input.Focus())
			}

		case key.Matches(msg, f.keymap.ShiftTab), key.Matches(msg, f.keymap.Up):
			if len(f.fields) > 0 {
				f.fields[f.focused].Input.Blur()
				f.focused = (f.focused - 1 + len(f.fields)) % len(f.fields)
				cmds = append(cmds, f.fields[f.focused].Input.Focus())
			}

		default:
			if f.focused < len(f.fields) {
				var cmd tea.Cmd
				f.fields[f.focused].Input, cmd = f.fields[f.focused].Input.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return f, tea.Batch(cmds...)
}

func (f *Form) View() string {
	var content strings.Builder

	title := f.styles.Title.Render(f.title)
	content.WriteString(title)
	content.WriteString("\n\n")

	for i, field := range f.fields {
		label := field.Label
		if field.Required {
			label += " *"
		}

		var labelStyle lipgloss.Style
		if i == f.focused {
			labelStyle = f.styles.Success
		} else {
			labelStyle = f.styles.Muted
		}

		content.WriteString(labelStyle.Render(label))
		content.WriteString("\n")

		var inputStyle lipgloss.Style
		if i == f.focused {
			inputStyle = f.styles.InputActive
		} else {
			inputStyle = f.styles.Input
		}

		inputView := inputStyle.Render(field.Input.View())
		content.WriteString(inputView)
		content.WriteString("\n\n")
	}

	if !f.validateForm() {
		errorText := f.styles.Error.Render("Please fill in all required fields")
		content.WriteString(errorText)
		content.WriteString("\n")
	}

	instructions := f.styles.Help.Render("tab/↑↓: navigate • enter: submit • esc: cancel")
	content.WriteString(instructions)

	return f.styles.Panel.
		Width(f.width - 4).
		Height(f.height - 4).
		Render(content.String())
}

// validateForm validates all required fields.
// Returns true if all required fields have non-empty values.
func (f *Form) validateForm() bool {
	for _, field := range f.fields {
		if field.Required && !f.isFieldValid(field) {
			return false
		}
	}
	return true
}

// isFieldValid checks if a field contains a valid value.
// For required fields, this means non-empty after trimming whitespace.
func (f *Form) isFieldValid(field Field) bool {
	return strings.TrimSpace(field.Input.Value()) != ""
}

// IsSubmitted returns true if the form was submitted
func (f *Form) IsSubmitted() bool {
	return f.submitted
}

// IsCancelled returns true if the form was cancelled
func (f *Form) IsCancelled() bool {
	return f.cancelled
}

// GetValue returns the value of a field by index
func (f *Form) GetValue(index int) string {
	if index < len(f.fields) {
		return strings.TrimSpace(f.fields[index].Input.Value())
	}
	return ""
}

// SetValue sets the value of a field by index.
// Does nothing if the index is out of bounds.
func (f *Form) SetValue(index int, value string) {
	if index < len(f.fields) {
		f.fields[index].Input.SetValue(value)
	}
}

// Reset resets the form state
func (f *Form) Reset() {
	f.submitted = false
	f.cancelled = false
	f.focused = 0

	for i := range f.fields {
		f.fields[i].Input.SetValue("")
		f.fields[i].Input.Blur()
	}

	if len(f.fields) > 0 {
		f.fields[0].Input.Focus()
	}
}
