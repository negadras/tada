package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/negadras/tada/internal/todo"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))
)

type TableModel struct {
	table table.Model
	todos []*todo.Todo
}

// NewTableModel creates a new table model with todos
func NewTableModel(todos []*todo.Todo) TableModel {
	// Define columns
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Priority", Width: 12},
		{Title: "Status", Width: 8},
		{Title: "Age", Width: 10},
		{Title: "Description", Width: 80},
	}

	rows := make([]table.Row, len(todos))
	for i, t := range todos {
		// helper func to get priority icon
		priorityIcon := todo.GetPriorityIcon(t.Priority)

		id := fmt.Sprintf("#%d", t.ID)
		priority := priorityIcon + " " + t.Priority.String()
		status := t.Status.String()
		age := todo.FormatAge(t.Age())
		description := t.Description

		rows[i] = table.Row{
			id,
			priority,
			status,
			age,
			description,
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = selectedStyle
	t.SetStyles(s)

	return TableModel{
		table: t,
		todos: todos,
	}
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if len(m.todos) > 0 && m.table.Cursor() < len(m.todos) {
				selectedTodo := m.todos[m.table.Cursor()]
				fmt.Printf("\nSelected todo #%d: %s\n", selectedTodo.ID, selectedTodo.Description)
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - 4)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	var b strings.Builder

	b.WriteString("\n📋 Todo List\n\n")
	b.WriteString(baseStyle.Render(m.table.View()))
	b.WriteString("\n")
	b.WriteString("  ↑/↓: Navigate • Enter: Select • q: Quit\n")

	return b.String()
}

func ShowTable(todos []*todo.Todo) error {
	if len(todos) == 0 {
		fmt.Println("📝 No todos found matching your criteria.")
		return nil
	}

	m := NewTableModel(todos)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
