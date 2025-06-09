package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/negadras/tada/internal/todo"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tada",
		Short: "A command line todo app",
		Long: `tada is a CLI that will help you add todo list, list your 
todo list, edit, close ...`,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceUsage = true
		},
	}

	// Main commands
	addCmd := newAddTadaCommand()
	listCmd := newListTadaCommand()
	updateCmd := newUpdateCommand()
	deleteCmd := newDeleteCommand()

	cmd.AddCommand(newQuoteCommand())
	cmd.AddCommand(addCmd)
	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(listCmd)
	cmd.AddCommand(updateCmd)
	cmd.AddCommand(deleteCmd)

	// Aliases
	cmd.AddCommand(createAlias("ls", listCmd))
	cmd.AddCommand(createAlias("rm", deleteCmd))
	cmd.AddCommand(createAlias("del", deleteCmd))

	cmd.AddCommand(createDoneCommand())
	cmd.AddCommand(createOpenCommand())

	// Help for aliases
	cmd.AddCommand(createAliasesCommand())

	return cmd
}

// createAlias creates an alias for an existing command
func createAlias(alias string, originalCmd *cobra.Command) *cobra.Command {
	aliasCmd := &cobra.Command{
		Use:   alias,
		Short: fmt.Sprintf("Alias for '%s'", originalCmd.Name()),
		Long:  originalCmd.Long,
		Args:  originalCmd.Args,
		RunE:  originalCmd.RunE,
	}

	aliasCmd.Flags().AddFlagSet(originalCmd.Flags())

	return aliasCmd
}

// createDoneCommand creates a convenience command for marking todos as done
func createDoneCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "done [id]",
		Short: "Mark a todo as done (alias for 'update [id] --status done')",
		Example: `  # Mark todo #5 as done
  tada done 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create update command and set the status flag
			updateCmd := newUpdateCommand()
			updateCmd.SetArgs(append(args, "--status", "done"))
			return updateCmd.Execute()
		},
	}
}

// createOpenCommand creates a convenience command for marking todos as open
func createOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open [id]",
		Short: "Mark a todo as open (alias for 'update [id] --status open')",
		Example: `  # Mark todo #5 as open
  tada open 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create update command and set the status flag
			updateCmd := newUpdateCommand()
			updateCmd.SetArgs(append(args, "--status", "open"))
			return updateCmd.Execute()
		},
	}
}

// createAliasesCommand shows all available aliases
func createAliasesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "aliases",
		Short: "Show all available command aliases",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("📚 Available aliases:")
			cmd.Println()
			cmd.Println("  ls      → list         (List todos)")
			cmd.Println("  rm      → delete       (Remove a todo)")
			cmd.Println("  del     → delete       (Remove a todo)")
			cmd.Println("  done    → update --status done")
			cmd.Println("  open    → update --status open")
			cmd.Println()
			cmd.Println("Examples:")
			cmd.Println("  tada ls                    # List all todos")
			cmd.Println("  tada rm 5                  # Delete todo #5")
			cmd.Println("  tada done 3                # Mark todo #3 as done")
		},
	}
}

func newDeleteCommand() *cobra.Command {
	return &cobra.Command{
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
}

func Execute() {
	cmd := newRootCommand()

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", color.RedString("error:"), err)
		os.Exit(1)
	}
}
