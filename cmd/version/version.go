package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Version is set during build time
var Version = "2.0.0"

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "All software has versions. This is tada's",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("tada version %s\n", Version)
		},
	}
}
