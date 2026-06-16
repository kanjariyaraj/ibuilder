package cmd

import (
	"fmt"
	"os"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repository connection",
	Long:  "Connect, inspect, and validate GitHub repository settings.",
}

var repoConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Detect and save repository from git remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		remote, err := github.DetectGitRemote()
		if err != nil {
			return fmt.Errorf("no git remote 'origin' found: %w\nMake sure you are in a git repository with a remote 'origin' configured.", err)
		}

		path := cfgFile
		if path == "" {
			cwd, _ := os.Getwd()
			path = cwd + "/builder.json"
		}

		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		cfg.Repo.Owner = remote.Owner
		cfg.Repo.Name = remote.Name

		if err := config.Save(path, cfg); err != nil {
			return err
		}

		fmt.Printf("Repository connected: %s/%s\n", remote.Owner, remote.Name)
		fmt.Printf("  Saved to: %s\n", path)
		return nil
	},
}

var repoInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display repository metadata from GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated. Run 'builder auth github' first")
		}

		path := cfgFile
		if path == "" {
			cwd, _ := os.Getwd()
			path = cwd + "/builder.json"
		}

		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		if cfg.Repo.Owner == "" || cfg.Repo.Name == "" {
			fmt.Println("No repository configured. Run 'builder repo connect' first.")
			return nil
		}

		client := github.NewClient(tokenData.AccessToken)
		info, err := github.FetchRepoInfo(client, cfg.Repo.Owner, cfg.Repo.Name)
		if err != nil {
			return err
		}

		fmt.Printf("Owner:           %s\n", info.Owner)
		fmt.Printf("Repository:      %s\n", info.Name)
		fmt.Printf("Full Name:       %s\n", info.FullName)
		fmt.Printf("Default Branch:  %s\n", info.DefaultBranch)
		fmt.Printf("Visibility:      %s\n", info.Visibility)
		fmt.Printf("Admin Access:    %v\n", info.Permissions.Admin)
		fmt.Printf("Push Access:     %v\n", info.Permissions.Push)
		fmt.Printf("Pull Access:     %v\n", info.Permissions.Pull)
		fmt.Printf("Actions Enabled: %v\n", info.ActionsEnabled)
		return nil
	},
}

var repoValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate repository access and configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil {
			return fmt.Errorf("error loading token: %w", err)
		}

		path := cfgFile
		if path == "" {
			cwd, _ := os.Getwd()
			path = cwd + "/builder.json"
		}

		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		token := ""
		if tokenData != nil {
			token = tokenData.AccessToken
		}

		result := github.ValidateAll(token, cfg.Repo.Owner, cfg.Repo.Name)

		fmt.Println("Validation Results:")
		fmt.Printf("  Authenticated:   %v\n", result.Authenticated)
		fmt.Printf("  Repository:      %v\n", result.RepoExists)
		fmt.Printf("  Actions Enabled: %v\n", result.ActionsEnabled)
		fmt.Printf("  Can Push:        %v\n", result.CanPush)
		fmt.Printf("  Can Admin:       %v\n", result.CanAdmin)

		if len(result.Errors) > 0 {
			fmt.Println("\nIssues Found:")
			for _, e := range result.Errors {
				fmt.Printf("  - %s\n", e)
			}
		} else {
			fmt.Println("\nEverything looks good!")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoConnectCmd)
	repoCmd.AddCommand(repoInfoCmd)
	repoCmd.AddCommand(repoValidateCmd)
}
