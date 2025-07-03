package list

import (
	"github.com/negadras/tada/internal/todo"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List todo tasks",
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
			// Get database connection
			db, cleanup, err := todo.GetDB(cmd)
			if err != nil {
				return nil
			}
			defer cleanup()

			// Parse status filter
			statusFlag, _ := cmd.Flags().GetString("status")
			var statusFilter *todo.Status

			if statusFlag != "all" && statusFlag != "a" {
				status, err := todo.ParseStatus(statusFlag)
				if err != nil {
					todo.PrintError(cmd, err)
					return nil
				}
				statusFilter = &status
			}

			// Parse priority filter
			priorityFlag, _ := cmd.Flags().GetString("priority")
			var priorityFilter *todo.Priority

			if priorityFlag != "all" && priorityFlag != "a" {
				priority, err := todo.ParsePriority(priorityFlag)
				if err != nil {
					todo.PrintError(cmd, err)
					return nil
				}
				priorityFilter = &priority
			}

			// List todos
			tasks, err := db.List(statusFilter, priorityFilter)
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			todo.PrintList(cmd, tasks)
			return nil
		},
	}

	cmd.Flags().StringP("status", "s", "open", "Status filter (open/o, done/d, all/a)")
	cmd.Flags().StringP("priority", "p", "all", "Priority filter (low/l, medium/m, high/h, all/a)")
	return cmd
}
