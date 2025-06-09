package cmd

import (
	"fmt"
	"time"

	"github.com/negadras/tada/cmd/db"
	"github.com/negadras/tada/utils"
	"github.com/spf13/cobra"
)

func newListTadaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all todos",
		Long: `List all the todo tasks with optional filtering.

Status filtering: 
  open, o - tasks that are open (default)
  done, d - tasks that are completed
  all, a  - show all tasks regardless of status

Priority filtering:
  low, l     - Low priority tasks
  medium, m  - Medium priority tasks  
  high, h    - High priority tasks
  all, a     - Show all priorities (default)`,
		Example: `  # List all open tasks (default)
  tada list
  
  # List all tasks regardless of status
  tada list --status all
  
  # List only completed tasks
  tada list --status done
  
  # List high priority open tasks
  tada list --priority high
  
  # List all high priority tasks (open and done)
  tada list --status all --priority high`,
		RunE: func(cmd *cobra.Command, args []string) error {
			todoDB, err := db.GetTodoDB()
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer todoDB.Close()

			// Parse status filter
			statusFlag, _ := cmd.Flags().GetString("status")
			var statusFilter *db.Status

			if statusFlag != "all" && statusFlag != "a" {
				status, err := utils.ParseStatus(statusFlag)
				if err != nil {
					return fmt.Errorf("invalid status '%s': %w", statusFlag, err)
				}
				statusFilter = &status
			}

			// Parse priority filter
			priorityFlag, _ := cmd.Flags().GetString("priority")
			var priorityFilter *db.Priority

			if priorityFlag != "all" && priorityFlag != "a" {
				priority, err := utils.ParsePriority(priorityFlag)
				if err != nil {
					return fmt.Errorf("invalid priority '%s': %w", priorityFlag, err)
				}
				priorityFilter = &priority
			}

			tasks, err := todoDB.ListTodos(statusFilter, priorityFilter)
			if err != nil {
				return fmt.Errorf("failed to list todos: %w", err)
			}

			if len(tasks) == 0 {
				cmd.Println("📝 No todos found matching your criteria.")
				return nil
			}

			displayTodos(cmd, tasks)

			return nil
		},
	}

	cmd.Flags().StringP("status", "s", "open", "Status filter (open/o, done/d, all/a)")
	cmd.Flags().StringP("priority", "p", "all", "Priority filter (low/l, medium/m, high/h, all/a)")
	return cmd
}

// displayTodos formats and displays the tada list
func displayTodos(cmd *cobra.Command, todos []*db.Todo) {
	cmd.Printf("📋 Found %d todo(s):\n\n", len(todos))

	for _, todo := range todos {
		priorityIcon := "🟡"
		switch todo.Priority {
		case db.Low:
			priorityIcon = "🟢"
		case db.High:
			priorityIcon = "🔴"
		}

		age := formatAge(todo.Age())

		cmd.Printf("%s [#%d] %s\n", priorityIcon, todo.ID, todo.Description)
		cmd.Printf("   Priority: %-8s Status: %-6s Age: %s\n",
			todo.Priority.String(),
			todo.Status.String(),
			age)

		if todo.Status == db.Done && todo.CompletedAt != nil {
			completedAge := formatAge(*todo.CompletedAge())
			cmd.Printf("   Completed: %s ago\n", completedAge)
		}

		cmd.Println()
	}
}

// formatAge converts duration to human-readable format
func formatAge(duration time.Duration) string {
	if duration < time.Minute {
		seconds := int(duration.Seconds())
		return fmt.Sprintf("%dseconds", seconds)
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%dminutes", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dhours", hours)
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%ddays", days)
	}
}
