package cmd

import (
	"strconv"

	"github.com/negadras/tada/internal/todo"
	"github.com/spf13/cobra"
)

func newDoneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "done [id]",
		Short: "Mark a todo as done",
		Long:  "Mark a todo task as completed by its ID",
		Example: `  # Mark todo #5 as done
  tada done 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			dbPath, err := todo.GetDatabasePath()
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			db, err := todo.NewDB(dbPath)
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}
			defer db.Close()

			if err := db.UpdateStatus(id, todo.Done); err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			updateTodo, err := db.Get(id)
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			todo.PrintSuccess(cmd, "Marked as done:")
			todo.PrintTodo(cmd, updateTodo)
			return nil
		},
	}

	return cmd
}
