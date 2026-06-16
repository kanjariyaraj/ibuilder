package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Builder configuration",
	Long:  `View, create, and validate Builder configuration.`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default builder.json configuration",
	Run: func(cmd *cobra.Command, args []string) {
		path := cfgFile
		if path == "" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			path = filepath.Join(cwd, "builder.json")
		}

		cfg := config.Default()
		if err := config.Save(path, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created default configuration: %s\n", path)
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		path := cfgFile
		if path == "" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			path = filepath.Join(cwd, "builder.json")
		}

		cfg, err := config.Load(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration: %s\n", path)
		fmt.Printf("Project:       %s\n", cfg.ProjectName)
		fmt.Printf("Repository:    %s\n", cfg.Repository)
		fmt.Printf("Repo Owner:    %s\n", cfg.Repo.Owner)
		fmt.Printf("Repo Name:     %s\n", cfg.Repo.Name)
		fmt.Printf("Repo Branch:   %s\n", cfg.Repo.Branch)
		fmt.Printf("GitHub Auth:   %v\n", cfg.GitHub.Authenticated)
		fmt.Printf("Build Workflow: %s\n", cfg.Build.WorkflowID)
		fmt.Printf("Build Branch:  %s\n", cfg.Build.Branch)
		fmt.Printf("Build Scheme:  %s\n", cfg.Build.Scheme)
		fmt.Printf("Build Mode:    %s\n", cfg.Build.BuildMode)
		fmt.Printf("iOS Min Ver:   %s\n", cfg.IOS.MinimumVersion)
		fmt.Printf("iOS Target:    %s\n", cfg.IOS.TargetVersion)
		fmt.Printf("Devices:       %v\n", cfg.IOS.Devices)
		fmt.Printf("MobAI:         %v\n", cfg.MobAI.Enabled)
		fmt.Printf("Flutter:       %v (channel: %s)\n", cfg.Flutter.Enabled, cfg.Flutter.Channel)
		fmt.Printf("React Native:  %v (entry: %s)\n", cfg.ReactNative.Enabled, cfg.ReactNative.Entry)
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		path := cfgFile
		if path == "" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			path = filepath.Join(cwd, "builder.json")
		}

		cfg, err := config.Load(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		errs := config.Validate(cfg)
		if len(errs) > 0 {
			fmt.Println("Configuration has errors:")
			for _, e := range errs {
				fmt.Printf("  - %v\n", e)
			}
		} else {
			fmt.Println("Configuration is valid.")
		}
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
}
