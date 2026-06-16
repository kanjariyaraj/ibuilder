package cmd

import (
	"fmt"
	"os"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage GitHub authentication",
	Long:  "Authenticate with GitHub using device flow, check status, or logout.",
}

var authGithubCmd = &cobra.Command{
	Use:   "github",
	Short: "Authenticate with GitHub via device flow",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := github.Authenticate("BuilderCLI")
		if err != nil {
			return err
		}

		path := cfgFile
		if path == "" {
			cwd, _ := os.Getwd()
			path = cwd + "/builder.json"
		}

		cfg, err := config.Load(path)
		if err != nil {
			return fmt.Errorf("warning: could not update config: %w", err)
		}

		cfg.GitHub.Authenticated = true

		if token.Scope != "" {
			fmt.Printf("Token scopes: %s\n", token.Scope)
		}

		if err := config.Save(path, cfg); err != nil {
			return fmt.Errorf("warning: could not save config: %w", err)
		}
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check GitHub authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		authenticated, token, err := github.AuthStatus()
		if err != nil {
			return fmt.Errorf("error checking auth status: %w", err)
		}

		if authenticated && token != nil {
			fmt.Println("Authenticated with GitHub")
			if token.Scope != "" {
				fmt.Printf("  Scopes: %s\n", token.Scope)
			}
		} else if token != nil {
			fmt.Println("Token found but invalid or expired")
			fmt.Println("  Run 'builder auth github' to re-authenticate.")
		} else {
			fmt.Println("Not authenticated with GitHub")
			fmt.Println("  Run 'builder auth github' to authenticate.")
		}
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove GitHub authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := github.Logout(); err != nil {
			return fmt.Errorf("error logging out: %w", err)
		}

		path := cfgFile
		if path == "" {
			cwd, _ := os.Getwd()
			path = cwd + "/builder.json"
		}

		cfg, err := config.Load(path)
		if err == nil {
			cfg.GitHub.Authenticated = false
			config.Save(path, cfg)
		}

		fmt.Println("Logged out successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authGithubCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}
