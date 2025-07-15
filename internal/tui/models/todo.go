package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/negadras/tada/internal/todo"
	"github.com/negadras/tada/internal/tui/components"
	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

// Table column width constants
const (
	todoIDWidth          = 6
	todoPriorityWidth    = 12
	todoStatusWidth      = 8
	todoAgeWidth         = 10
	todoDescriptionWidth = 60
	todoTableHeight      = 15
	todoTableMargin      = 4
	todoTableReserved    = 10
)

// TodosLoadedMsg is sent when todos are loaded from the database
type TodosLoadedMsg struct {
	Todos []*todo.Todo
}

// TodoErrorMsg is sent when there's an error loading todos
type TodoErrorMsg struct {
	Error error
}

// TodoManager represents the todo management screen
type TodoManager struct {
	BaseModel
	styles            *styles.Styles
	keymap            utils.KeyMap
	table             table.Model
	todos             []*todo.Todo
	db                *todo.DB
	loading           bool
	errorMessage      string
	statusFilter      *todo.Status
	addForm           *components.Form
	showAddForm       bool
	editForm          *components.Form
	showEditForm      bool
	editingTodo       *todo.Todo
	showDeleteConfirm bool
	todoToDelete      *todo.Todo
}

// NewTodoManager creates a new todo manager model
func NewTodoManager(styles *styles.Styles, keymap utils.KeyMap) *TodoManager {
	// Create table columns
	columns := []table.Column{
		{Title: "ID", Width: todoIDWidth},
		{Title: "Priority", Width: todoPriorityWidth},
		{Title: "Status", Width: todoStatusWidth},
		{Title: "Age", Width: todoAgeWidth},
		{Title: "Description", Width: todoDescriptionWidth},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(todoTableHeight),
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
	addForm := components.NewForm(styles, keymap, "Add New Todo")
	addForm.AddField("Description", "Enter todo description", true)
	addForm.AddField("Priority", "low, medium, or high (default: medium)", false)

	// Create edit form
	editForm := components.NewForm(styles, keymap, "Edit Todo")
	editForm.AddField("Description", "Enter todo description", true)
	editForm.AddField("Priority", "low, medium, or high", false)

	return &TodoManager{
		styles:   styles,
		keymap:   keymap,
		table:    t,
		todos:    []*todo.Todo{},
		loading:  true,
		addForm:  addForm,
		editForm: editForm,
	}
}

// Init initializes the todo manager
func (t *TodoManager) Init() tea.Cmd {
	return t.loadTodos()
}

// Update handles messages for the todo manager
func (t *TodoManager) Update(msg tea.Msg) (*TodoManager, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case TodosLoadedMsg:
		t.todos = msg.Todos
		t.loading = false
		t.errorMessage = ""
		t.updateTable()
		return t, nil

	case TodoErrorMsg:
		t.loading = false
		t.errorMessage = msg.Error.Error()
		return t, nil

	case tea.KeyMsg:
		// Handle add form mode
		if t.showAddForm {
			t.addForm, cmd = t.addForm.Update(msg)

			if t.addForm.IsSubmitted() {
				description := t.addForm.GetValue(0)
				priorityStr := t.addForm.GetValue(1)

				// Parse priority
				priority := todo.Medium // default
				if priorityStr != "" {
					if p, err := todo.ParsePriority(priorityStr); err == nil {
						priority = p
					}
				}

				t.showAddForm = false
				t.addForm.Reset()
				return t, t.createTodo(description, priority)
			}

			if t.addForm.IsCancelled() {
				t.showAddForm = false
				t.addForm.Reset()
				return t, nil
			}

			return t, cmd
		}

		// Handle edit form mode
		if t.showEditForm {
			t.editForm, cmd = t.editForm.Update(msg)

			if t.editForm.IsSubmitted() {
				// Update todo
				description := t.editForm.GetValue(0)
				priorityStr := t.editForm.GetValue(1)

				// Parse priority
				priority := todo.Medium // default
				if priorityStr != "" {
					if p, err := todo.ParsePriority(priorityStr); err == nil {
						priority = p
					}
				}

				t.showEditForm = false
				t.editForm.Reset()
				todoID := t.editingTodo.ID
				t.editingTodo = nil
				return t, t.updateTodo(todoID, description, priority)
			}

			if t.editForm.IsCancelled() {
				t.showEditForm = false
				t.editForm.Reset()
				t.editingTodo = nil
				return t, nil
			}

			return t, cmd
		}

		// Handle delete confirmation mode
		if t.showDeleteConfirm {
			switch {
			case key.Matches(msg, t.keymap.Enter):
				// Confirm deletion
				t.showDeleteConfirm = false
				todoToDelete := t.todoToDelete
				t.todoToDelete = nil
				return t, t.confirmDeleteTodo(todoToDelete)
			case key.Matches(msg, t.keymap.Escape):
				// Cancel deletion
				t.showDeleteConfirm = false
				t.todoToDelete = nil
				return t, nil
			}
			return t, nil
		}

		// Handle table mode
		switch {
		case key.Matches(msg, t.keymap.Add):
			t.showAddForm = true
			t.addForm.SetSize(t.width, t.height)
			return t, t.addForm.Init()

		case key.Matches(msg, t.keymap.Edit):
			return t, t.openEditForm()

		case key.Matches(msg, t.keymap.Delete):
			return t, t.showDeleteConfirmation()

		case key.Matches(msg, t.keymap.Toggle):
			return t, t.toggleTodoStatus()

		case key.Matches(msg, t.keymap.Filter):
			return t, t.cycleStatusFilter()

		case key.Matches(msg, t.keymap.Enter):
			return t, t.toggleTodoStatus()

		default:
			// Handle table navigation
			t.table, cmd = t.table.Update(msg)
			return t, cmd
		}
	}

	return t, nil
}

// View renders the todo manager
func (t *TodoManager) View() string {
	if t.showAddForm {
		return t.addForm.View()
	}

	if t.showEditForm {
		return t.editForm.View()
	}

	// Show delete confirmation if active
	if t.showDeleteConfirm {
		return t.renderDeleteConfirmation()
	}

	var content strings.Builder

	title := t.styles.Title.Render("📝 Todo Management")
	stats := t.renderStats()
	header := lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", stats)
	content.WriteString(header)
	content.WriteString("\n\n")

	if t.statusFilter != nil {
		filterText := t.styles.Info.Render(fmt.Sprintf("Filter: %s", t.statusFilter.String()))
		content.WriteString(filterText)
		content.WriteString("\n")
	}

	if t.errorMessage != "" {
		errorText := t.styles.Error.Render(fmt.Sprintf("Error: %s", t.errorMessage))
		content.WriteString(errorText)
		content.WriteString("\n")
	}

	if t.loading {
		loadingText := t.styles.Info.Render("Loading todos...")
		content.WriteString(loadingText)
		content.WriteString("\n")
	}

	if !t.loading && t.errorMessage == "" {
		if len(t.todos) == 0 {
			emptyText := t.styles.Muted.Render("No todos found. Press 'a' to add one!")
			content.WriteString(emptyText)
		} else {
			content.WriteString(t.table.View())
		}
	}

	content.WriteString("\n")

	instructions := t.renderInstructions()
	content.WriteString(instructions)

	return t.styles.Content.Render(content.String())
}

// renderStats renders the todo statistics
func (t *TodoManager) renderStats() string {
	if len(t.todos) == 0 {
		return t.styles.Muted.Render("No todos")
	}

	var completed, total int
	for _, todoItem := range t.todos {
		total++
		if todoItem.Status == todo.Done {
			completed++
		}
	}

	completionRate := float64(completed) / float64(total) * 100
	statsText := fmt.Sprintf("%d total • %d completed • %.1f%% done", total, completed, completionRate)
	return t.styles.Success.Render(statsText)
}

// renderInstructions renders the keyboard instructions
func (t *TodoManager) renderInstructions() string {
	instructions := []string{
		"a: add todo",
		"e: edit",
		"d: delete",
		"t/enter: toggle status",
		"f: filter",
		"esc: back",
	}

	return t.styles.Help.Render(strings.Join(instructions, " • "))
}

// updateTable updates the table with current todos.
// Formats each todo's data to fit within the designated column widths.
func (t *TodoManager) updateTable() {
	rows := make([]table.Row, len(t.todos))
	for i, todo := range t.todos {
		rows[i] = table.Row{
			fmt.Sprintf("#%d", todo.ID),
			strings.ToUpper(todo.Priority.String()),
			strings.ToUpper(todo.Status.String()),
			utils.FormatDuration(todo.Age()),
			todo.Description,
		}
	}

	t.table.SetRows(rows)
}

// loadTodos loads todos from the database and applies the current status filter.
// Returns a command that will send either TodosLoadedMsg or TodoErrorMsg.
func (t *TodoManager) loadTodos() tea.Cmd {
	return func() tea.Msg {
		dbPath, err := todo.GetDatabasePath()
		if err != nil {
			return TodoErrorMsg{Error: fmt.Errorf("failed to get database path: %w", err)}
		}

		db, err := todo.NewDB(dbPath)
		if err != nil {
			return TodoErrorMsg{Error: fmt.Errorf("failed to open database: %w", err)}
		}

		t.db = db

		todos, err := t.db.List(t.statusFilter, nil)
		if err != nil {
			return TodoErrorMsg{Error: fmt.Errorf("failed to load todos: %w", err)}
		}

		return TodosLoadedMsg{Todos: todos}
	}
}

// toggleTodoStatus toggles the status of the selected todo
func (t *TodoManager) toggleTodoStatus() tea.Cmd {
	if len(t.todos) == 0 || t.db == nil {
		return nil
	}

	selectedIndex := t.table.Cursor()
	if selectedIndex >= len(t.todos) {
		return nil
	}

	selectedTodo := t.todos[selectedIndex]

	return func() tea.Msg {
		var newStatus todo.Status
		if selectedTodo.Status == todo.Open {
			newStatus = todo.Done
		} else {
			newStatus = todo.Open
		}

		err := t.db.UpdateStatus(selectedTodo.ID, newStatus)
		if err != nil {
			return TodoErrorMsg{Error: err}
		}

		todos, err := t.db.List(t.statusFilter, nil)
		if err != nil {
			return TodoErrorMsg{Error: err}
		}

		return TodosLoadedMsg{Todos: todos}
	}
}

// cycleStatusFilter cycles through status filters
func (t *TodoManager) cycleStatusFilter() tea.Cmd {
	return func() tea.Msg {
		if t.statusFilter == nil {
			openStatus := todo.Open
			t.statusFilter = &openStatus
		} else if *t.statusFilter == todo.Open {
			doneStatus := todo.Done
			t.statusFilter = &doneStatus
		} else {
			t.statusFilter = nil
		}

		// Reload todos with new filter
		todos, err := t.db.List(t.statusFilter, nil)
		if err != nil {
			return TodoErrorMsg{Error: err}
		}

		return TodosLoadedMsg{Todos: todos}
	}
}

// createTodo creates a new todo with the given description and priority.
// Returns a command that will send either TodosLoadedMsg or TodoErrorMsg.
func (t *TodoManager) createTodo(description string, priority todo.Priority) tea.Cmd {
	if t.db == nil {
		return nil
	}

	return func() tea.Msg {
		_, err := t.db.Create(description, priority)
		if err != nil {
			return TodoErrorMsg{Error: err}
		}

		todos, err := t.db.List(t.statusFilter, nil)
		if err != nil {
			return TodoErrorMsg{Error: err}
		}

		return TodosLoadedMsg{Todos: todos}
	}
}

// openEditForm opens the edit form for the currently selected todo.
// Pre-populates the form with the todo's current description and priority.
// Returns nil if no todos exist or selection is invalid.
func (t *TodoManager) openEditForm() tea.Cmd {
	if len(t.todos) == 0 {
		return nil
	}

	selectedIndex := t.table.Cursor()
	if selectedIndex >= len(t.todos) {
		return nil
	}

	selectedTodo := t.todos[selectedIndex]
	t.editingTodo = selectedTodo

	t.editForm.Reset()
	t.editForm.SetValue(0, selectedTodo.Description)
	t.editForm.SetValue(1, selectedTodo.Priority.String())

	t.showEditForm = true
	t.editForm.SetSize(t.width, t.height)
	return t.editForm.Init()
}

// updateTodo updates an existing todo's description and priority in the database.
// Calls both UpdateDescription and UpdatePriority methods, then reloads the todo list.
// Returns a command that will send either TodosLoadedMsg or TodoErrorMsg.
func (t *TodoManager) updateTodo(id int, description string, priority todo.Priority) tea.Cmd {
	if t.db == nil {
		return nil
	}

	return func() tea.Msg {

		err := t.db.UpdateDescription(id, description)
		if err != nil {
			return TodoErrorMsg{Error: err}
		}

		err = t.db.UpdatePriority(id, priority)
		if err != nil {
			return TodoErrorMsg{Error: err}
		}

		todos, err := t.db.List(t.statusFilter, nil)
		if err != nil {
			return TodoErrorMsg{Error: err}
		}

		return TodosLoadedMsg{Todos: todos}
	}
}

// showDeleteConfirmation shows the delete confirmation dialog for the selected todo.
// Returns nil if no todos exist or selection is invalid.
func (t *TodoManager) showDeleteConfirmation() tea.Cmd {
	if len(t.todos) == 0 {
		return nil
	}

	selectedIndex := t.table.Cursor()
	if selectedIndex >= len(t.todos) {
		return nil
	}

	t.todoToDelete = t.todos[selectedIndex]
	t.showDeleteConfirm = true
	return nil
}

// confirmDeleteTodo actually deletes the todo after confirmation.
// Returns a command that will send either TodosLoadedMsg or TodoErrorMsg.
func (t *TodoManager) confirmDeleteTodo(todoToDelete *todo.Todo) tea.Cmd {
	if t.db == nil || todoToDelete == nil {
		return nil
	}

	return func() tea.Msg {
		err := t.db.Delete(todoToDelete.ID)
		if err != nil {
			return TodoErrorMsg{Error: fmt.Errorf("failed to delete todo: %w", err)}
		}

		// Reload todos
		todos, err := t.db.List(t.statusFilter, nil)
		if err != nil {
			return TodoErrorMsg{Error: fmt.Errorf("failed to reload todos: %w", err)}
		}

		return TodosLoadedMsg{Todos: todos}
	}
}

// renderDeleteConfirmation renders the delete confirmation dialog.
func (t *TodoManager) renderDeleteConfirmation() string {
	if t.todoToDelete == nil {
		return ""
	}

	var content strings.Builder

	// Title
	title := t.styles.Title.Render("🗑️ Delete Todo")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Warning message
	warningMsg := t.styles.Error.Render("Are you sure you want to delete this todo?")
	content.WriteString(warningMsg)
	content.WriteString("\n\n")

	// Todo details
	todoDetails := fmt.Sprintf("ID: #%d\nDescription: %s\nPriority: %s\nStatus: %s",
		t.todoToDelete.ID,
		t.todoToDelete.Description,
		strings.ToUpper(t.todoToDelete.Priority.String()),
		strings.ToUpper(t.todoToDelete.Status.String()),
	)
	content.WriteString(t.styles.Muted.Render(todoDetails))
	content.WriteString("\n\n")

	// Instructions
	instructions := t.styles.Help.Render("enter: confirm deletion • esc: cancel")
	content.WriteString(instructions)

	// Wrap in panel
	return t.styles.Panel.
		Width(t.width - 4).
		Height(t.height - 4).
		Render(content.String())
}

// SetSize sets the size of the todo manager and updates child components.
func (t *TodoManager) SetSize(width, height int) {
	t.BaseModel.SetSize(width, height)
	t.table.SetWidth(width - todoTableMargin)
	t.table.SetHeight(height - todoTableReserved)

	if t.addForm != nil {
		t.addForm.SetSize(width, height)
	}
	if t.editForm != nil {
		t.editForm.SetSize(width, height)
	}
}
