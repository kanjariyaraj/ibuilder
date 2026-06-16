package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kanjariyaraj/Builder/internal/ai"
	"github.com/kanjariyaraj/Builder/internal/logger"
	"github.com/spf13/cobra"
)

func getAISession() *ai.Session {
	log := logger.New(logger.LevelInfo)
	return ai.NewSession(log)
}

func getProjectDir() string {
	cwd, _ := os.Getwd()
	if cfgFile != "" {
		return filepath.Dir(cfgFile)
	}
	return cwd
}

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI-powered diagnostics and troubleshooting",
	Long:  "Automatically analyze build failures, generate fix suggestions, and diagnose issues.",
}

var aiExplainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain the latest build or workflow failure",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getAISession()
		dir := getProjectDir()
		session.SetProjectDir(dir)

		result, err := session.Explain(dir)
		if err != nil {
			return fmt.Errorf("explain failed: %w", err)
		}

		fmt.Println("AI Diagnostic Explanation")
		fmt.Println("=========================")
		fmt.Printf("  Problem:    %s\n", result.Problem)
		fmt.Printf("  Category:   %s\n", result.Category)
		fmt.Printf("  Cause:      %s\n", result.Cause)
		fmt.Printf("  Impact:     %s\n", result.Impact)
		fmt.Printf("  Fix:        %s\n", result.Fix)
		fmt.Printf("  Confidence: %s (%.0f%%)\n", result.ConfidenceLabel, result.Confidence*100)
		fmt.Printf("  Time:       %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))

		return nil
	},
}

var aiAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze build and workflow logs for issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getAISession()
		dir := getProjectDir()
		session.SetProjectDir(dir)

		logs, err := session.Explain(dir)
		if err != nil {
			return fmt.Errorf("analyze failed: %w", err)
		}

		fmt.Println("AI Analysis Results")
		fmt.Println("===================")
		fmt.Printf("  Problem:       %s\n", logs.Problem)
		fmt.Printf("  Root Cause:    %s\n", logs.Cause)
		fmt.Printf("  Impact:        %s\n", logs.Impact)
		fmt.Printf("  Fix:           %s\n", logs.Fix)
		fmt.Printf("  Confidence:    %.0f%%\n", logs.Confidence*100)

		return nil
	},
}

var aiDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run full repository audit",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getAISession()
		dir := getProjectDir()
		session.SetProjectDir(dir)

		report := session.Doctor()

		fmt.Println("AI Doctor Report")
		fmt.Println("================")
		for _, check := range report.Checks {
			symbol := "✓"
			switch check.Status {
			case "WARNING":
				symbol = "⚠"
			case "FAILURE":
				symbol = "✗"
			}
			fmt.Printf("  %s [%s] %s\n", symbol, check.Status, check.Name)
			fmt.Printf("       %s\n", check.Message)
			if check.Suggest != "" {
				fmt.Printf("       Suggestion: %s\n", check.Suggest)
			}
		}

		fmt.Println()
		if report.Healthy {
			fmt.Println("Overall: HEALTHY")
		} else {
			fmt.Println("Overall: ISSUES FOUND")
		}

		return nil
	},
}

var aiFixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Generate fix suggestions for detected issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getAISession()
		dir := getProjectDir()
		session.SetProjectDir(dir)

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		apply, _ := cmd.Flags().GetBool("apply")

		if apply && !dryRun {
			fmt.Println("WARNING: Applying fixes automatically. Use --dry-run first to preview.")
		}

		result, err := session.AnalyzeAndSuggestFixes(dir)
		if err != nil {
			return fmt.Errorf("fix analysis failed: %w", err)
		}

		fmt.Println("AI Fix Suggestions")
		fmt.Println("=================")
		fmt.Printf("  Summary: %s\n", result.Summary)
		fmt.Println()

		for _, fix := range result.Fixes {
			fmt.Printf("  File:    %s\n", fix.File)
			fmt.Printf("  Current: %s\n", truncateStr(fix.Current, 50))
			fmt.Printf("  Suggest: %s\n", truncateStr(fix.Suggest, 50))
			fmt.Printf("  Reason:  %s\n", fix.Reason)
			fmt.Println()

			if apply && !dryRun {
				applyResult, err := session.ApplyFix(dir, fix, false)
				if err != nil {
					fmt.Printf("  ✗ Failed to apply: %v\n", err)
				} else {
					fmt.Printf("  ✓ %s\n", applyResult.Summary)
				}
			}
		}

		if dryRun && len(result.Fixes) > 0 {
			fmt.Println("  Run with --apply to apply these fixes.")
		}

		return nil
	},
}

var aiReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate AI diagnostic report",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getAISession()
		dir := getProjectDir()
		session.SetProjectDir(dir)

		format, _ := cmd.Flags().GetString("format")
		rf := ai.FormatMarkdown
		switch format {
		case "json":
			rf = ai.FormatJSON
		case "html":
			rf = ai.FormatHTML
		}

		path, err := session.GenerateReport(dir, rf)
		if err != nil {
			return fmt.Errorf("report generation failed: %w", err)
		}

		fmt.Printf("AI report generated: %s\n", path)
		return nil
	},
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func init() {
	rootCmd.AddCommand(aiCmd)
	aiCmd.AddCommand(aiExplainCmd)
	aiCmd.AddCommand(aiAnalyzeCmd)
	aiCmd.AddCommand(aiDoctorCmd)
	aiCmd.AddCommand(aiFixCmd)
	aiCmd.AddCommand(aiReportCmd)

	aiFixCmd.Flags().Bool("dry-run", true, "Preview fixes without applying")
	aiFixCmd.Flags().Bool("apply", false, "Apply suggested fixes (requires confirmation)")

	aiReportCmd.Flags().String("format", "markdown", "Report format: markdown, json, or html")
}
