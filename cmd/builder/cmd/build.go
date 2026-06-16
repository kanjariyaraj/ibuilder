package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Manage iOS builds",
	Long:  "Trigger, monitor, and manage iOS builds via GitHub Actions.",
}

var buildRunCmd = &cobra.Command{
	Use:   "run [workflow-id]",
	Short: "Trigger an iOS build workflow",
	Args:  cobra.ExactArgs(1),
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
			return fmt.Errorf("no repository configured. Run 'builder repo connect' first")
		}

		workflowID := args[0]
		branch, _ := cmd.Flags().GetString("branch")
		if branch == "" {
			branch = cfg.Repo.Branch
		}

		scheme, _ := cmd.Flags().GetString("scheme")
		mode, _ := cmd.Flags().GetString("mode")

		inputs := make(map[string]string)
		if scheme != "" {
			inputs["scheme"] = scheme
		}
		if mode != "" {
			inputs["build_mode"] = mode
		}

		client := github.NewClient(tokenData.AccessToken)
		if err := github.DispatchWorkflow(client, cfg.Repo.Owner, cfg.Repo.Name, workflowID, branch, inputs); err != nil {
			return fmt.Errorf("failed to trigger build: %w", err)
		}

		fmt.Printf("Build triggered! Workflow: %s, Branch: %s\n", workflowID, branch)
		fmt.Println("Use 'builder build list' to check build status.")
		return nil
	},
}

var buildStatusCmd = &cobra.Command{
	Use:   "status [run-id]",
	Short: "Check build status",
	Args:  cobra.ExactArgs(1),
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

		runID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[0])
		}

		client := github.NewClient(tokenData.AccessToken)
		run, err := github.GetWorkflowRun(client, cfg.Repo.Owner, cfg.Repo.Name, runID)
		if err != nil {
			return fmt.Errorf("failed to get build status: %w", err)
		}

		fmt.Printf("Build #%d\n", run.RunNumber)
		fmt.Printf("Status:      %s\n", run.Status)
		fmt.Printf("Conclusion:  %s\n", run.Conclusion)
		fmt.Printf("Branch:      %s\n", run.HeadBranch)
		fmt.Printf("Created:     %s\n", run.CreatedAt.Format("Jan 2, 2006 15:04:05"))
		fmt.Printf("Updated:     %s\n", run.UpdatedAt.Format("Jan 2, 2006 15:04:05"))
		fmt.Printf("URL:         %s\n", run.HTMLURL)

		return nil
	},
}

var buildListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent builds",
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

		limit, _ := cmd.Flags().GetInt("limit")
		if limit <= 0 {
			limit = 10
		}

		client := github.NewClient(tokenData.AccessToken)
		runs, err := github.ListWorkflowRuns(client, cfg.Repo.Owner, cfg.Repo.Name, limit)
		if err != nil {
			return fmt.Errorf("failed to list builds: %w", err)
		}

		if len(runs) == 0 {
			fmt.Println("No builds found.")
			return nil
		}

		fmt.Printf("Recent builds (last %d):\n", limit)
		fmt.Println()
		for _, run := range runs {
			status := run.Status
			if run.Conclusion != "" && run.Conclusion != "null" {
				status = run.Conclusion
			}
			fmt.Printf("  #%-5d %-12s %s\n", run.RunNumber, status, run.CreatedAt.Format("Jan 2 15:04"))
		}

		return nil
	},
}

var buildLogCmd = &cobra.Command{
	Use:   "log [run-id]",
	Short: "View build logs",
	Args:  cobra.ExactArgs(1),
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

		runID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[0])
		}

		client := github.NewClient(tokenData.AccessToken)
		run, err := github.GetWorkflowRun(client, cfg.Repo.Owner, cfg.Repo.Name, runID)
		if err != nil {
			return fmt.Errorf("failed to get build: %w", err)
		}

		fmt.Printf("Build #%d - %s (%s)\n", run.RunNumber, run.Status, run.Conclusion)
		fmt.Printf("URL: %s\n", run.HTMLURL)
		fmt.Println()
		fmt.Println("Use the URL above to view detailed logs in GitHub.")
		return nil
	},
}

var buildArtifactsCmd = &cobra.Command{
	Use:   "artifacts [run-id]",
	Short: "List build artifacts",
	Args:  cobra.ExactArgs(1),
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

		runID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[0])
		}

		client := github.NewClient(tokenData.AccessToken)
		artifacts, err := github.ListArtifacts(client, cfg.Repo.Owner, cfg.Repo.Name, runID)
		if err != nil {
			return fmt.Errorf("failed to list artifacts: %w", err)
		}

		if len(artifacts) == 0 {
			fmt.Println("No artifacts found for this build.")
			return nil
		}

		fmt.Println("Build Artifacts:")
		for _, a := range artifacts {
			size := float64(a.Size) / 1024
			fmt.Printf("  %s (%.1f KB)\n", a.Name, size)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.AddCommand(buildRunCmd)
	buildCmd.AddCommand(buildStatusCmd)
	buildCmd.AddCommand(buildListCmd)
	buildCmd.AddCommand(buildLogCmd)
	buildCmd.AddCommand(buildArtifactsCmd)

	buildRunCmd.Flags().String("branch", "", "Branch to build from")
	buildRunCmd.Flags().String("scheme", "", "Xcode scheme (for Xcode builds)")
	buildRunCmd.Flags().String("mode", "release", "Build mode (debug/release)")
	buildListCmd.Flags().Int("limit", 10, "Number of builds to show")
}
