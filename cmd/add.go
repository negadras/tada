package cmd

import (
	"fmt"
	"github.com/negadras/tada/cmd/db"
	"github.com/negadras/tada/utils"
	"github.com/spf13/cobra"
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
			todoDB, err := db.GetTodoDB()
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer todoDB.Close()

			description := strings.TrimSpace(args[0])
			if description == "" {
				return fmt.Errorf("task description cannot be empty")
			}

			priorityFlag, _ := cmd.Flags().GetString("priority")
			priority, err := utils.ParsePriority(priorityFlag)
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
