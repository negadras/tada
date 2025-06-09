package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Version is set during build time
var Version = "1.3.0"

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of tada",
		Long:  "All software has versions. This is tada's",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("tada version %s\n", Version)
		},
	}
}
