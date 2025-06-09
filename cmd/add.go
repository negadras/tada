package cmd

import (
	"strings"

	"github.com/negadras/tada/internal/todo"
	"github.com/spf13/cobra"
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
			// Get database path
			dbPath, err := todo.GetDatabasePath()
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			// Open database
			db, err := todo.NewDB(dbPath)
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}
			defer db.Close()

			// Validate description
			description := strings.TrimSpace(args[0])
			if err := todo.ValidateDescription(description); err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			// Parse priority
			priorityFlag, _ := cmd.Flags().GetString("priority")
			priority, err := todo.ParsePriority(priorityFlag)
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			// Create todo
			newTodo, err := db.Create(description, priority)
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			todo.PrintCreated(cmd, newTodo)
			return nil
		},
	}

	cmd.Flags().StringP("priority", "p", "medium", "Priority level (low/l, medium/m, high/h)")
	return cmd
}
