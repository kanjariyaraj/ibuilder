package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kanjariyaraj/Builder/internal/build"
	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
	"github.com/spf13/cobra"
)

var iosCmd = &cobra.Command{
	Use:   "ios",
	Short: "iOS build commands",
	Long:  "Build, manage, and download iOS builds via GitHub Actions.",
}

var iosBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build iOS app remotely via GitHub Actions",
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

		workflowID, _ := cmd.Flags().GetString("workflow")
		if workflowID == "" {
			workflowID = cfg.Build.WorkflowID
		}
		if workflowID == "" {
			return fmt.Errorf("workflow ID required. Set in builder.json or use --workflow")
		}

		branch, _ := cmd.Flags().GetString("branch")
		if branch == "" {
			branch = cfg.Build.Branch
		}

		scheme, _ := cmd.Flags().GetString("scheme")
		if scheme == "" {
			scheme = cfg.Build.Scheme
		}

		mode, _ := cmd.Flags().GetString("mode")
		if mode == "" {
			mode = cfg.Build.BuildMode
		}

		wait, _ := cmd.Flags().GetBool("wait")
		logs, _ := cmd.Flags().GetBool("logs")
		jsonOut, _ := cmd.Flags().GetBool("json")
		downloadOnly, _ := cmd.Flags().GetBool("download-only")
		clean, _ := cmd.Flags().GetBool("clean")

		cwd, _ := os.Getwd()

		opts := &build.BuildOptions{
			WorkflowID:   workflowID,
			Branch:       branch,
			Scheme:       scheme,
			Mode:         mode,
			Wait:         wait,
			Logs:         logs,
			JSON:         jsonOut,
			DownloadOnly: downloadOnly,
			Clean:        clean,
			Token:        tokenData.AccessToken,
			Owner:        cfg.Repo.Owner,
			Name:         cfg.Repo.Name,
			Dir:          cwd,
		}

		if clean {
			os.RemoveAll(cwd + "/dist")
			fmt.Println("Cleaned dist/ directory")
			if downloadOnly {
				result, err := build.RunBuild(opts)
				if err != nil {
					return err
				}
				if jsonOut {
					data, _ := json.MarshalIndent(result, "", "  ")
					fmt.Println(string(data))
				}
				return nil
			}
			return nil
		}

		result, err := build.RunBuild(opts)
		if err != nil {
			return err
		}

		if jsonOut {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(iosCmd)
	iosCmd.AddCommand(iosBuildCmd)

	iosBuildCmd.Flags().String("workflow", "", "Workflow ID or filename")
	iosBuildCmd.Flags().String("branch", "", "Branch to build from")
	iosBuildCmd.Flags().String("scheme", "", "Xcode scheme")
	iosBuildCmd.Flags().String("mode", "release", "Build mode (debug/release)")
	iosBuildCmd.Flags().Bool("wait", false, "Wait for build completion")
	iosBuildCmd.Flags().Bool("logs", false, "Stream build logs")
	iosBuildCmd.Flags().Bool("json", false, "JSON output")
	iosBuildCmd.Flags().Bool("download-only", false, "Download latest artifact without building")
	iosBuildCmd.Flags().Bool("clean", false, "Clean dist/ directory")
}
