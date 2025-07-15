package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/negadras/tada/cmd/add"
	"github.com/negadras/tada/cmd/delete"
	"github.com/negadras/tada/cmd/list"
	"github.com/negadras/tada/cmd/quote"
	"github.com/negadras/tada/cmd/update"
	"github.com/negadras/tada/cmd/version"
	"github.com/negadras/tada/internal/tui"
	"github.com/spf13/cobra"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if TUI mode is requested
			tuiMode, _ := cmd.Flags().GetBool("tui")
			if tuiMode {
				screen, _ := cmd.Flags().GetString("screen")
				return tui.RunWithScreen(screen)
			}
			
			// Default behavior: show help
			return cmd.Help()
		},
	}

	// Main commands
	addCmd := add.NewCommand()
	listCmd := list.NewCommand()
	updateCmd := update.NewCommand()
	deleteCmd := delete.NewCommand()

	cmd.AddCommand(quote.NewCommand())
	cmd.AddCommand(addCmd)
	cmd.AddCommand(version.NewCommand())
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

	// Add TUI flags
	cmd.Flags().BoolP("tui", "t", false, "Launch interactive TUI mode")
	cmd.Flags().StringP("screen", "s", "", "Launch TUI at specific screen (dashboard, todos, quotes)")

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
			updateCmd := update.NewCommand()
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
			updateCmd := update.NewCommand()
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

func Execute() {
	cmd := newRootCommand()

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", color.RedString("error:"), err)
		os.Exit(1)
	}
}
