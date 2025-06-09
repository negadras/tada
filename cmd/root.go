package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
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

	cmd.AddCommand(newQuoteCommand())
	cmd.AddCommand(newAddTadaCommand())

	return cmd
}

func Execute() {
	cmd := newRootCommand()

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", color.RedString("error:"), err)
		os.Exit(1)
	}
}
