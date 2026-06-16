package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/kanjariyaraj/Builder/internal/artifacts"
	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Manage build artifacts",
	Long:  "List, download, inspect, and clean build artifacts.",
}

var artifactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated")
		}

		cfgPath := cfgFile
		if cfgPath == "" {
			cwd, _ := os.Getwd()
			cfgPath = cwd + "/builder.json"
		}
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}

		client := github.NewClient(tokenData.AccessToken)
		mgr := artifacts.NewArtifactManager(client, cfg.Repo.Owner, cfg.Repo.Name)

		limit, _ := cmd.Flags().GetInt("limit")
		all, _ := cmd.Flags().GetBool("all")
		asJSON, _ := cmd.Flags().GetBool("json")

		artifactsList, err := mgr.ListAllArtifacts()
		if err != nil {
			return fmt.Errorf("failed to list artifacts: %w", err)
		}

		if !all && limit > 0 && len(artifactsList) > limit {
			artifactsList = artifactsList[:limit]
		}

		if asJSON {
			return printJSON(artifactsList)
		}

		if len(artifactsList) == 0 {
			fmt.Println("No artifacts found.")
			return nil
		}

		fmt.Printf("%-5s %-30s %-10s %-20s %-10s\n", "ID", "Name", "Size", "Created", "Status")
		for _, a := range artifactsList {
			fmt.Printf("%-5d %-30s %-10d %-20s %-10s\n", a.ID, a.Name, a.Size, a.CreatedAt.Format("2006-01-02 15:04"), a.Status)
		}
		return nil
	},
}

var artifactDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated")
		}

		cfgPath := cfgFile
		if cfgPath == "" {
			cwd, _ := os.Getwd()
			cfgPath = cwd + "/builder.json"
		}
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}

		client := github.NewClient(tokenData.AccessToken)
		mgr := artifacts.NewArtifactManager(client, cfg.Repo.Owner, cfg.Repo.Name)

		destDir, _ := cmd.Flags().GetString("dest")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		latest, _ := cmd.Flags().GetBool("latest")
		name, _ := cmd.Flags().GetString("name")
		build, _ := cmd.Flags().GetInt("build")

		opts := &artifacts.DownloadOptions{
			DestDir:   destDir,
			Overwrite: overwrite,
			Latest:    latest,
			Name:      name,
			Build:     build,
		}

		result, err := mgr.DownloadArtifact(opts)
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		fmt.Printf("Downloaded to: %s\n", result.Path)
		fmt.Printf("Size: %d bytes\n", result.Size)
		fmt.Printf("SHA256: %s\n", result.Checksum)
		return nil
	},
}

var artifactInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect artifact details",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		if id == 0 {
			return fmt.Errorf("--id is required")
		}

		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated")
		}

		cfgPath := cfgFile
		if cfgPath == "" {
			cwd, _ := os.Getwd()
			cfgPath = cwd + "/builder.json"
		}
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}

		client := github.NewClient(tokenData.AccessToken)
		mgr := artifacts.NewArtifactManager(client, cfg.Repo.Owner, cfg.Repo.Name)

		all, err := mgr.ListAllArtifacts()
		if err != nil {
			return err
		}

		for _, a := range all {
			if a.ID == id {
				fmt.Printf("ID:       %d\n", a.ID)
				fmt.Printf("Name:     %s\n", a.Name)
				fmt.Printf("Size:     %d bytes\n", a.Size)
				fmt.Printf("Build:    %d\n", a.BuildNumber)
				fmt.Printf("Created:  %s\n", a.CreatedAt.Format(time.RFC3339))
				fmt.Printf("Expires:  %s\n", a.ExpiresAt.Format(time.RFC3339))
				fmt.Printf("Status:   %s\n", a.Status)
				return nil
			}
		}
		return fmt.Errorf("artifact %d not found", id)
	},
}

var artifactLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Download latest artifact",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenData, err := github.LoadToken()
		if err != nil || tokenData == nil {
			return fmt.Errorf("not authenticated")
		}

		cfgPath := cfgFile
		if cfgPath == "" {
			cwd, _ := os.Getwd()
			cfgPath = cwd + "/builder.json"
		}
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}

		client := github.NewClient(tokenData.AccessToken)
		mgr := artifacts.NewArtifactManager(client, cfg.Repo.Owner, cfg.Repo.Name)

		result, err := mgr.DownloadLatest()
		if err != nil {
			return fmt.Errorf("failed to download latest: %w", err)
		}

		fmt.Printf("Latest artifact saved to: %s\n", result.Path)
		return nil
	},
}

var artifactCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up local artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		olderThan, _ := cmd.Flags().GetString("older-than")
		keep, _ := cmd.Flags().GetInt("keep")
		all, _ := cmd.Flags().GetBool("all")

		storage := artifacts.NewStorage()

		if all {
			if err := storage.CleanAll(); err != nil {
				return fmt.Errorf("cleanup failed: %w", err)
			}
			fmt.Println("All local artifacts removed.")
			return nil
		}

		if olderThan != "" {
			dur, err := parseDuration(olderThan)
			if err != nil {
				return fmt.Errorf("invalid duration: %w", err)
			}
			count, err := storage.CleanOldArtifacts(dur)
			if err != nil {
				return fmt.Errorf("cleanup failed: %w", err)
			}
			fmt.Printf("Removed %d old artifacts.\n", count)
		}

		if keep > 0 {
			fmt.Printf("Keeping %d artifacts.\n", keep)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	artifactCmd.AddCommand(artifactListCmd)
	artifactCmd.AddCommand(artifactDownloadCmd)
	artifactCmd.AddCommand(artifactInspectCmd)
	artifactCmd.AddCommand(artifactLatestCmd)
	artifactCmd.AddCommand(artifactCleanCmd)

	artifactListCmd.Flags().Int("limit", 30, "Limit number of artifacts")
	artifactListCmd.Flags().Bool("all", false, "Show all artifacts")
	artifactListCmd.Flags().Bool("json", false, "JSON output")

	artifactDownloadCmd.Flags().String("dest", "dist", "Destination directory")
	artifactDownloadCmd.Flags().Bool("overwrite", false, "Overwrite existing files")
	artifactDownloadCmd.Flags().Bool("latest", false, "Download latest artifact")
	artifactDownloadCmd.Flags().String("name", "", "Download artifact by name")
	artifactDownloadCmd.Flags().Int("build", 0, "Download artifact by build number")
	artifactInspectCmd.Flags().Int64("id", 0, "Artifact ID")

	artifactCleanCmd.Flags().String("older-than", "", "Remove artifacts older than duration (e.g. 72h, 7d)")
	artifactCleanCmd.Flags().Int("keep", 0, "Keep N most recent artifacts")
	artifactCleanCmd.Flags().Bool("all", false, "Remove all local artifacts")
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("empty duration")
	}
	return time.ParseDuration(s)
}
