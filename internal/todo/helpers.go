package todo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

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

func PrintTodo(cmd *cobra.Command, todo *Todo) {
	priorityIcon := GetPriorityIcon(todo.Priority)
	age := FormatAge(todo.Age())

	cmd.Printf("%s [#%d] %s\n", priorityIcon, todo.ID, todo.Description)
	cmd.Printf("   Priority: %-8s Status: %-6s Age: %s\n",
		todo.Priority.String(),
		todo.Status.String(),
		age)

	if todo.Status == Done && todo.CompletedAt != nil {
		completedAge := FormatAge(*todo.CompletedAge())
		cmd.Printf("   Completed: %s ago\n", completedAge)
	}
}

func PrintCreated(cmd *cobra.Command, todo *Todo) {
	cmd.Printf("✅ Created todo #%d: %s\n", todo.ID, todo.Description)
	cmd.Printf("   Priority: %s\n", todo.Priority.String())
	cmd.Printf("   Status: %s\n", todo.Status.String())
}

func PrintError(cmd *cobra.Command, err error) {
	cmd.Printf("❌ Error: %v\n", err)
}

func PrintSuccess(cmd *cobra.Command, message string) {
	cmd.Printf("✅ %s\n", message)
}

func GetDB(cmd *cobra.Command) (*DB, func(), error) {
	dbPath, err := GetDatabasePath()
	if err != nil {
		PrintError(cmd, err)
		return nil, nil, err
	}

	db, err := NewDB(dbPath)
	if err != nil {
		PrintError(cmd, err)
		return nil, nil, err
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup, nil
}

func FormatAge(duration time.Duration) string {
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

func GetPriorityIcon(priority Priority) string {
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
