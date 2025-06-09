package cmd

import (
	"fmt"
	"github.com/negadras/tada/cmd/db"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func newAddTadaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [description]",
		Short: "Add a todo task",
		Long: `Add a new todo task with the specified description and priority.
Priority levels:
  low, l     - Low priority (default)
  medium, m  - Medium priority
  high, h    - High priority`,
		Example: `  # Add a high priority task
  tada add "Fix the login bug" --priority high
  
  # Add with short priority flag
  tada add "Write documentation" -p medium
  
  # Add low priority task (default)
  tada add "Clean up code"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			todoDB, err := getTodoDB()
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer todoDB.Close()

			description := strings.TrimSpace(args[0])
			if description == "" {
				return fmt.Errorf("task description cannot be empty")
			}

			priorityFlag, _ := cmd.Flags().GetString("priority")
			priority, err := parsePriority(priorityFlag)
			if err != nil {
				return fmt.Errorf("invalid priority '%s':%w", priorityFlag, err)
			}

			todo, err := todoDB.CreateTodo(description, priority)
			if err != nil {
				return fmt.Errorf("failed to create todo: %w", err)
			}

			cmd.Printf("✅ Created todo #%d: %s\n", todo.ID, todo.Description)
			cmd.Printf("   Priority: %s\n", priority.String())
			cmd.Printf("   Status: %s\n", todo.Status.String())

			return nil
		},
	}

	cmd.Flags().StringP("priority", "p", "medium", "Priority level (low/l, medium/m, high/h)")
	return cmd
}

// getTodoDB initializes the database connection with proper path handling
func getTodoDB() (*db.TodoDB, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create .tada directory if it doesn't exist
	tadaDir := filepath.Join(homeDir, ".tada")
	if err := os.MkdirAll(tadaDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .tada directory: %w", err)
	}

	// Database file path
	dbPath := filepath.Join(tadaDir, "todos.db")

	return db.NewTodoDB(dbPath)
}

// parsePriority converts string priority to Priority enum
func parsePriority(priorityStr string) (db.Priority, error) {
	switch strings.ToLower(strings.TrimSpace(priorityStr)) {
	case "low", "l", "1":
		return db.Low, nil
	case "medium", "m", "2":
		return db.Medium, nil
	case "high", "h", "3":
		return db.High, nil
	default:
		return db.Medium, fmt.Errorf("must be one of: low, medium, high (or l, m, h)")
	}
}
