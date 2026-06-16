package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "0.1.0"
	Commit    = "none"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Builder",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Builder v%s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built:  %s\n", BuildDate)
	},
}
