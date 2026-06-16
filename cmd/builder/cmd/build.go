package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/kanjariyaraj/Builder/internal/artifacts"
	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/github"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Manage builds",
	Long:  "View build history, inspect builds, fetch logs, and open build URLs.",
}

var buildHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show build history",
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

		branch, _ := cmd.Flags().GetString("branch")
		workflow, _ := cmd.Flags().GetString("workflow")
		status, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")
		page, _ := cmd.Flags().GetInt("page")
		asJSON, _ := cmd.Flags().GetBool("json")

		records, err := mgr.GetHistory(&artifacts.HistoryOptions{
			Branch:   branch,
			Workflow: workflow,
			Status:   status,
			Limit:    limit,
			Page:     page,
			JSON:     asJSON,
		})
		if err != nil {
			return fmt.Errorf("failed to get build history: %w", err)
		}

		if asJSON {
			return printJSON(records)
		}

		if len(records) == 0 {
			fmt.Println("No build history found.")
			return nil
		}

		fmt.Printf("%-6s %-10s %-12s %-20s %-10s\n", "#", "Status", "Conclusion", "Branch", "Duration")
		for _, r := range records {
			fmt.Printf("%-6d %-10s %-12s %-20s %-10s\n", r.RunNumber, r.Status, r.Conclusion, r.Branch, r.Duration)
		}
		return nil
	},
}

var buildInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a build",
	RunE: func(cmd *cobra.Command, args []string) error {
		runID, _ := cmd.Flags().GetInt64("run-id")
		if runID == 0 {
			return fmt.Errorf("--run-id is required")
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

		inspect, err := mgr.InspectBuild(runID)
		if err != nil {
			return fmt.Errorf("failed to inspect build: %w", err)
		}

		fmt.Printf("Build #%d\n", inspect.RunNumber)
		fmt.Printf("Status:     %s\n", inspect.Status)
		fmt.Printf("Conclusion: %s\n", inspect.Conclusion)
		fmt.Printf("Branch:     %s\n", inspect.Branch)
		fmt.Printf("Commit:     %s\n", inspect.CommitSHA)
		fmt.Printf("Author:     %s\n", inspect.Author)
		fmt.Printf("Workflow:   %s\n", inspect.Workflow)
		fmt.Printf("Duration:   %s\n", inspect.Duration)
		fmt.Printf("URL:        %s\n", inspect.URL)

		if len(inspect.Artifacts) > 0 {
			fmt.Println("\nArtifacts:")
			for _, a := range inspect.Artifacts {
				fmt.Printf("  - %s (%d bytes)\n", a.Name, a.Size)
			}
		}

		if len(inspect.Jobs) > 0 {
			fmt.Println("\nJobs:")
			for _, j := range inspect.Jobs {
				fmt.Printf("  - %s [%s/%s] (%d steps)\n", j.Name, j.Status, j.Conclusion, j.Steps)
			}
		}
		return nil
	},
}

var buildLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Fetch build logs",
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

		runID, _ := cmd.Flags().GetInt64("run-id")
		latest, _ := cmd.Flags().GetBool("latest")
		savePath, _ := cmd.Flags().GetString("save")

		logPath, err := mgr.FetchLogs(&artifacts.LogsOptions{
			RunID:    runID,
			Latest:   latest,
			SavePath: savePath,
		})
		if err != nil {
			return fmt.Errorf("failed to fetch logs: %w", err)
		}

		fmt.Printf("Logs saved to: %s\n", logPath)
		return nil
	},
}

var buildOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open build URL in browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		runID, _ := cmd.Flags().GetInt64("run-id")
		if runID == 0 {
			return fmt.Errorf("--run-id is required")
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

		url, err := mgr.OpenBuildURL(runID)
		if err != nil {
			return fmt.Errorf("failed to get build URL: %w", err)
		}

		fmt.Printf("Opening: %s\n", url)

		var cmd2 *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd2 = exec.Command("open", url)
		case "linux":
			cmd2 = exec.Command("xdg-open", url)
		case "windows":
			cmd2 = exec.Command("cmd", "/c", "start", url)
		default:
			return fmt.Errorf("unsupported OS for opening browser")
		}

		return cmd2.Start()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.AddCommand(buildHistoryCmd)
	buildCmd.AddCommand(buildInspectCmd)
	buildCmd.AddCommand(buildLogsCmd)
	buildCmd.AddCommand(buildOpenCmd)

	buildHistoryCmd.Flags().String("branch", "", "Filter by branch")
	buildHistoryCmd.Flags().String("workflow", "", "Filter by workflow")
	buildHistoryCmd.Flags().String("status", "", "Filter by status")
	buildHistoryCmd.Flags().Int("limit", 30, "Number of builds")
	buildHistoryCmd.Flags().Int("page", 1, "Page number")
	buildHistoryCmd.Flags().Bool("json", false, "JSON output")

	buildInspectCmd.Flags().Int64("run-id", 0, "Build run ID")
	buildLogsCmd.Flags().Int64("run-id", 0, "Build run ID")
	buildLogsCmd.Flags().Bool("latest", false, "Latest build logs")
	buildLogsCmd.Flags().String("save", "", "Save path for logs")
	buildOpenCmd.Flags().Int64("run-id", 0, "Build run ID")
}
