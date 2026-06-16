package cmd

import (
	"fmt"
	"os"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
	"github.com/kanjariyaraj/Builder/internal/signing"
	"github.com/spf13/cobra"
)

var signingCmd = &cobra.Command{
	Use:   "signing",
	Short: "Manage iOS code signing",
	Long:  "Configure, validate, and manage iOS code signing certificates and provisioning profiles.",
}

var signingSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup iOS code signing with certificate and provisioning profile",
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

		certPath, _ := cmd.Flags().GetString("cert")
		certPass, _ := cmd.Flags().GetString("password")
		profilePath, _ := cmd.Flags().GetString("profile")
		teamID, _ := cmd.Flags().GetString("team-id")
		bundleID, _ := cmd.Flags().GetString("bundle-id")
		force, _ := cmd.Flags().GetBool("force")

		if certPath == "" || profilePath == "" {
			return fmt.Errorf("--cert and --profile are required")
		}

		client := github.NewClient(tokenData.AccessToken)
		opts := &signing.SetupOptions{
			CertPath:     certPath,
			CertPassword: certPass,
			ProfilePath:  profilePath,
			TeamID:       teamID,
			BundleID:     bundleID,
			Force:        force,
		}

		result, err := signing.RunSetup(client, opts, cfg)
		if err != nil {
			return fmt.Errorf("signing setup failed: %w", err)
		}

		cfg.Signing.TeamID = teamID
		if err := config.Save(path, cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println("Signing setup complete!")
		fmt.Printf("  Certificate secret: %s\n", result.CertSecret)
		fmt.Printf("  Profile secret:     %s\n", result.ProfileSecret)
		fmt.Println()
		fmt.Println("Secrets uploaded to GitHub repository.")
		return nil
	},
}

var signingValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate signing configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated")
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

		cert, _ := cmd.Flags().GetString("cert")
		if cert != "" {
			if err := signing.ValidateCertFile(cert); err != nil {
				return err
			}
			fmt.Println("Certificate: valid")
		}

		profile, _ := cmd.Flags().GetString("profile")
		if profile != "" {
			if err := signing.ValidateProfileFile(profile); err != nil {
				return err
			}
			fmt.Println("Profile: valid")
		}

		client := github.NewClient(tokenData.AccessToken)
		sc := &signing.SigningConfig{
			TeamID:        cfg.Signing.TeamID,
			CertSecret:    cfg.Signing.Certificate,
			ProfileSecret: cfg.Signing.Provisioning,
		}
		result := signing.ValidateSigning(client, cfg.Repo.Owner, cfg.Repo.Name, sc)
		if result.Valid {
			fmt.Println("All checks passed!")
		} else {
			fmt.Println("Issues found:")
			for _, e := range result.Errors {
				fmt.Printf("  - %s\n", e)
			}
		}
		return nil
	},
}

var signingDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose signing configuration issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated")
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

		client := github.NewClient(tokenData.AccessToken)
		sc := &signing.SigningConfig{
			TeamID:        cfg.Signing.TeamID,
			CertSecret:    cfg.Signing.Certificate,
			ProfileSecret: cfg.Signing.Provisioning,
		}
		result := signing.RunDoctor(client, cfg.Repo.Owner, cfg.Repo.Name, sc)

		fmt.Println("Signing Doctor Results:")
		fmt.Printf("  Certificates:     %v\n", result.CertOK)
		fmt.Printf("  Profiles:         %v\n", result.ProfileOK)
		fmt.Printf("  Secrets:          %v\n", result.SecretsOK)
		fmt.Printf("  Team ID:          %v\n", result.TeamIDSet)
		fmt.Printf("  Overall Health:   %v\n", result.Healthy)

		if len(result.Issues) > 0 {
			fmt.Println("\nIssues:")
			for _, i := range result.Issues {
				fmt.Printf("  - %s\n", i)
			}
		}
		return nil
	},
}

var signingStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display signing configuration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated")
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

		client := github.NewClient(tokenData.AccessToken)
		status := signing.GetStatus(client, cfg.Repo.Owner, cfg.Repo.Name, cfg)

		fmt.Println("Signing Status:")
		fmt.Printf("  Team ID:      %s\n", status.TeamID)
		fmt.Printf("  Configured:   %v\n", status.Configured)
		fmt.Printf("  Cert Secret:  %s\n", status.CertSecret)
		fmt.Printf("  Profile Sec:  %s\n", status.ProfileSecret)
		fmt.Printf("  Uploaded:     %v\n", status.SecretsUploaded)
		return nil
	},
}

var signingRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove signing secrets from GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated")
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

		client := github.NewClient(tokenData.AccessToken)
		if err := signing.RemoveSigning(client, cfg.Repo.Owner, cfg.Repo.Name); err != nil {
			return err
		}

		fmt.Println("Signing secrets removed from GitHub repository.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(signingCmd)
	signingCmd.AddCommand(signingSetupCmd)
	signingCmd.AddCommand(signingValidateCmd)
	signingCmd.AddCommand(signingDoctorCmd)
	signingCmd.AddCommand(signingStatusCmd)
	signingCmd.AddCommand(signingRemoveCmd)

	signingSetupCmd.Flags().String("cert", "", "Path to .p12 certificate file")
	signingSetupCmd.Flags().String("password", "", "Certificate password")
	signingSetupCmd.Flags().String("profile", "", "Path to .mobileprovision file")
	signingSetupCmd.Flags().String("team-id", "", "Apple Team ID")
	signingSetupCmd.Flags().String("bundle-id", "", "Bundle identifier")
	signingSetupCmd.Flags().Bool("force", false, "Overwrite existing secrets")
	signingValidateCmd.Flags().String("cert", "", "Path to certificate file to validate")
	signingValidateCmd.Flags().String("profile", "", "Path to provisioning profile to validate")
}
