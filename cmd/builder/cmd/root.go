package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "builder",
	Short: "Build, Test, Sign and Release iOS Apps From Anywhere",
	Long: `Builder is an open-source CLI that allows developers to build iOS apps
from Windows, Linux and WSL using GitHub-hosted macOS runners and real iPhone devices.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Builder - iOS Build Toolchain")
		fmt.Println("Use 'builder --help' for available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to builder.json config file")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(configCmd)
}
