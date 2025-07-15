package update

import (
	"fmt"
	"strconv"

	"github.com/negadras/tada/internal/todo"
	"github.com/negadras/tada/internal/tui"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a todo task",
		Long: `Update various properties of a todo task including status, priority, and description.
		
At least one flag must be provided to specify what to update.

💡 Tip: Use --tui flag to launch interactive edit mode`,
		Example: `  # Mark todo #5 as done
  tada update 5 --status done
  
  # Change priority to high
  tada update 5 --priority high
  
  # Update description
  tada update 5 --description "New description"
  
  # Update multiple properties at once
  tada update 5 --status done --priority low
  
  # Using short flags
  tada update 5 -s done -p high -d "Updated task"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if TUI mode is requested
			tuiMode, _ := cmd.Flags().GetBool("tui")
			if tuiMode {
				return tui.RunWithScreen("todos")
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			// Check if at least one flag is provided
			if !cmd.Flags().Changed("status") &&
				!cmd.Flags().Changed("priority") &&
				!cmd.Flags().Changed("description") {
				todo.PrintError(cmd, fmt.Errorf("at least one flag (--status, --priority, or --description) must be provided"))
				return nil
			}

			// Get database connection
			db, cleanup, err := todo.GetDB(cmd)
			if err != nil {
				return nil
			}
			defer cleanup()

			if cmd.Flags().Changed("status") {
				statusFlag, _ := cmd.Flags().GetString("status")
				status, err := todo.ParseStatus(statusFlag)
				if err != nil {
					todo.PrintError(cmd, err)
					return nil
				}
				if err := db.UpdateStatus(id, status); err != nil {
					todo.PrintError(cmd, err)
					return nil
				}
			}

			if cmd.Flags().Changed("priority") {
				priorityFlag, _ := cmd.Flags().GetString("priority")
				priority, err := todo.ParsePriority(priorityFlag)
				if err != nil {
					todo.PrintError(cmd, err)
					return nil
				}
				if err := db.UpdatePriority(id, priority); err != nil {
					todo.PrintError(cmd, err)
					return nil
				}
			}

			if cmd.Flags().Changed("description") {
				description, _ := cmd.Flags().GetString("description")
				if err := todo.ValidateDescription(description); err != nil {
					todo.PrintError(cmd, err)
					return nil
				}
				if err := db.UpdateDescription(id, description); err != nil {
					todo.PrintError(cmd, err)
					return nil
				}
			}

			updatedTodo, err := db.Get(id)
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			todo.PrintSuccess(cmd, "Updated todo:")
			todo.PrintTodo(cmd, updatedTodo)
			return nil
		},
	}

	cmd.Flags().StringP("status", "s", "", "Update status (open/o, done/d)")
	cmd.Flags().StringP("priority", "p", "", "Update priority (low/l, medium/m, high/h)")
	cmd.Flags().StringP("description", "d", "", "Update description")
	cmd.Flags().BoolP("tui", "t", false, "Launch interactive TUI mode for editing")

	return cmd
}
