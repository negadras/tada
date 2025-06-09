package todo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Validation helpers

// ValidateDescription checks if a description is valid
func ValidateDescription(desc string) error {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return errors.New("description cannot be empty")
	}
	if len(desc) > 255 {
		return errors.New("description too long (max 255 characters)")
	}
	return nil
}

// ParsePriority converts string priority to Priority enum
func ParsePriority(priorityStr string) (Priority, error) {
	switch strings.ToLower(strings.TrimSpace(priorityStr)) {
	case "low", "l", "1":
		return Low, nil
	case "medium", "m", "2":
		return Medium, nil
	case "high", "h", "3":
		return High, nil
	default:
		return Medium, fmt.Errorf("must be one of: low, medium, high (or l, m, h)")
	}
}

// ParseStatus converts string status to Status enum
func ParseStatus(status string) (Status, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "open", "o", "1":
		return Open, nil
	case "done", "d", "2":
		return Done, nil
	default:
		return Open, fmt.Errorf("must be one of: open or done")
	}
}

// Formatting helpers

// PrintList formats and displays a list of todos
func PrintList(cmd *cobra.Command, todos []*Todo) {
	if len(todos) == 0 {
		cmd.Println("📝 No todos found matching your criteria.")
		return
	}

	cmd.Printf("📋 Found %d todo(s):\n\n", len(todos))

	for _, todo := range todos {
		PrintTodo(cmd, todo)
		cmd.Println()
	}
}

// PrintTodo formats and displays a single todo
func PrintTodo(cmd *cobra.Command, todo *Todo) {
	priorityIcon := getPriorityIcon(todo.Priority)
	age := formatAge(todo.Age())

	cmd.Printf("%s [#%d] %s\n", priorityIcon, todo.ID, todo.Description)
	cmd.Printf("   Priority: %-8s Status: %-6s Age: %s\n",
		todo.Priority.String(),
		todo.Status.String(),
		age)

	if todo.Status == Done && todo.CompletedAt != nil {
		completedAge := formatAge(*todo.CompletedAge())
		cmd.Printf("   Completed: %s ago\n", completedAge)
	}
}

// PrintCreated formats the output when a todo is created
func PrintCreated(cmd *cobra.Command, todo *Todo) {
	cmd.Printf("✅ Created todo #%d: %s\n", todo.ID, todo.Description)
	cmd.Printf("   Priority: %s\n", todo.Priority.String())
	cmd.Printf("   Status: %s\n", todo.Status.String())
}

// PrintError formats error messages consistently
func PrintError(cmd *cobra.Command, err error) {
	cmd.Printf("❌ Error: %v\n", err)
}

// PrintSuccess formats success messages consistently
func PrintSuccess(cmd *cobra.Command, message string) {
	cmd.Printf("✅ %s\n", message)
}

// formatAge converts duration to human-readable format
func formatAge(duration time.Duration) string {
	if duration < time.Minute {
		seconds := int(duration.Seconds())
		if seconds <= 1 {
			return "1 second"
		}
		return fmt.Sprintf("%d seconds", seconds)
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
}

// getPriorityIcon returns an emoji icon for the priority level
func getPriorityIcon(priority Priority) string {
	switch priority {
	case Low:
		return "🟢"
	case Medium:
		return "🟡"
	case High:
		return "🔴"
	default:
		return "⚪"
	}
}
