package cmd

import (
	"github.com/negadras/tada/internal/todo"
	"github.com/spf13/cobra"
	"strconv"
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a todo",
		Long:  "Delete a todo task by its ID",
		Example: `  # Delete todo #5
  tada delete 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			// this thing looks like redudent accros the CRUD commands we have, what can we do about it
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

			todoItem, err := db.Get(id)
			if err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			if err := db.Delete(id); err != nil {
				todo.PrintError(cmd, err)
				return nil
			}

			cmd.Printf("🗑️  Deleted todo #%d: %s\n", id, todoItem.Description)
			return nil
		},
	}

	return cmd
}
