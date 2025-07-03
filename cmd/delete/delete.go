package delete

import (
	"github.com/negadras/tada/internal/todo"
	"github.com/spf13/cobra"
	"strconv"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a todo",
		Example: `  # Delete todo #5
  tada delete 5`,
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

			// Get todo before deletion
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
