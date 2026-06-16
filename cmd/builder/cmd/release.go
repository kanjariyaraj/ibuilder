package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kanjariyaraj/Builder/internal/logger"
	"github.com/kanjariyaraj/Builder/internal/release"
	"github.com/kanjariyaraj/Builder/internal/releasepipeline"
	"github.com/spf13/cobra"
)

func getReleaseSession() *release.Session {
	log := logger.New(logger.LevelInfo)
	session := release.NewSession(log)
	cwd, _ := os.Getwd()
	if cfgFile != "" {
		session.SetProjectDir(filepath.Dir(cfgFile))
	} else {
		session.SetProjectDir(cwd)
	}
	return session
}

var testflightCmd = &cobra.Command{
	Use:   "testflight",
	Short: "TestFlight release management",
	Long:  "Upload builds, manage beta groups, and track TestFlight releases.",
}

var testflightUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload IPA to TestFlight",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getReleaseSession()

		artifact, _ := cmd.Flags().GetString("artifact")
		build, _ := cmd.Flags().GetString("build")

		var result *release.UploadResult
		var err error

		if artifact != "" {
			result, err = session.UploadArtifact(artifact)
		} else if build != "" {
			result, err = session.UploadBuild(build)
		} else {
			result, err = session.UploadLatest()
		}

		if err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}

		fmt.Println("TestFlight upload complete!")
		fmt.Printf("  IPA:      %s\n", result.IPAPath)
		fmt.Printf("  Build:    %s\n", result.BuildNumber)
		fmt.Printf("  Version:  %s\n", result.Version)
		fmt.Printf("  Status:   %s\n", result.Status)
		fmt.Printf("  Time:     %s\n", result.UploadedAt.Format("2006-01-02 15:04:05"))
		return nil
	},
}

var testflightStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check TestFlight upload and processing status",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getReleaseSession()

		buildNumber, _ := cmd.Flags().GetString("build")
		var result *release.StatusResult
		var err error

		if buildNumber != "" {
			result, err = session.CheckBuildStatus(buildNumber)
		} else {
			result, err = session.CheckStatus()
		}

		if err != nil {
			return fmt.Errorf("status check failed: %w", err)
		}

		fmt.Println("TestFlight Status")
		fmt.Println("=================")
		fmt.Printf("  Build:          %s\n", result.BuildNumber)
		fmt.Printf("  Version:        %s\n", result.Version)
		fmt.Printf("  Upload:         %s\n", result.UploadState)
		fmt.Printf("  Processing:     %s\n", result.ProcessingState)
		fmt.Printf("  Beta:           %s\n", result.BetaState)
		fmt.Printf("  Review:         %s\n", result.ReviewState)
		fmt.Printf("  Availability:   %s\n", result.Availability)
		return nil
	},
}

var testflightGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List TestFlight beta groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getReleaseSession()

		groupName, _ := cmd.Flags().GetString("inspect")

		if groupName != "" {
			group, err := session.InspectGroup(groupName)
			if err != nil {
				return fmt.Errorf("group not found: %w", err)
			}
			fmt.Printf("Group: %s\n", group.Name)
			fmt.Printf("  Testers: %d\n", group.TesterCount)
			fmt.Printf("  Builds:  %d\n", group.BuildCount)
			return nil
		}

		groups, err := session.ListGroups()
		if err != nil {
			return fmt.Errorf("failed to list groups: %w", err)
		}

		fmt.Println("TestFlight Groups")
		fmt.Println("=================")
		for _, g := range groups.Groups {
			fmt.Printf("  %s\n", g.Name)
			fmt.Printf("    Testers: %d\n", g.TesterCount)
			fmt.Printf("    Builds:  %d\n", g.BuildCount)
		}
		fmt.Printf("\nTotal: %d group(s)\n", groups.Total)
		return nil
	},
}

var testflightBuildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "List TestFlight builds",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getReleaseSession()

		buildNumber, _ := cmd.Flags().GetString("build")

		if buildNumber != "" {
			build, err := session.GetBuild(buildNumber)
			if err != nil {
				return fmt.Errorf("build not found: %w", err)
			}
			fmt.Printf("Build: %s\n", build.BuildNumber)
			fmt.Printf("  Version:     %s\n", build.Version)
			fmt.Printf("  Uploaded:    %s\n", build.UploadDate.Format("2006-01-02"))
			fmt.Printf("  Status:      %s\n", build.Status)
			return nil
		}

		builds, err := session.ListBuilds()
		if err != nil {
			return fmt.Errorf("failed to list builds: %w", err)
		}

		fmt.Println("TestFlight Builds")
		fmt.Println("=================")
		for _, b := range builds.Builds {
			fmt.Printf("  Build #%s (v%s)\n", b.BuildNumber, b.Version)
			fmt.Printf("    Uploaded: %s\n", b.UploadDate.Format("2006-01-02"))
			fmt.Printf("    Status:   %s\n", b.Status)
		}
		fmt.Printf("\nTotal: %d build(s)\n", builds.Total)
		return nil
	},
}

var testflightTestersCmd = &cobra.Command{
	Use:   "testers",
	Short: "List TestFlight testers",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getReleaseSession()

		testers, err := session.ListTesters()
		if err != nil {
			return fmt.Errorf("failed to list testers: %w", err)
		}

		fmt.Println("TestFlight Testers")
		fmt.Println("==================")
		fmt.Printf("Internal Testers (%d):\n", len(testers.Internal))
		for _, t := range testers.Internal {
			fmt.Printf("  - %s (%s)\n", t.Name, t.Email)
		}
		fmt.Printf("\nExternal Testers (%d):\n", len(testers.External))
		for _, t := range testers.External {
			fmt.Printf("  - %s (%s)\n", t.Name, t.Email)
		}
		fmt.Printf("\nTotal: %d tester(s)\n", testers.Total)
		return nil
	},
}

func getPipelineSession() *releasepipeline.Pipeline {
	log := logger.New(logger.LevelInfo)
	p := releasepipeline.NewPipeline(log)
	cwd, _ := os.Getwd()
	if cfgFile != "" {
		p.SetProjectDir(filepath.Dir(cfgFile))
	} else {
		p.SetProjectDir(cwd)
	}
	return p
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "One-command release pipeline",
	Long: `Execute the complete release pipeline:
validate → build → sign → notes → upload → github release → report`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := getPipelineSession()
		started := time.Now()

		beta, _ := cmd.Flags().GetBool("beta")
		production, _ := cmd.Flags().GetBool("production")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		notesOnly, _ := cmd.Flags().GetBool("notes")
		jsonOut, _ := cmd.Flags().GetBool("json")

		mode := releasepipeline.ModeBeta
		if production {
			mode = releasepipeline.ModeProduction
		}
		if !beta && !production {
			mode = releasepipeline.ModeBeta
		}

		p.SetMode(mode)
		p.SetDryRun(dryRun)
		p.StartStatus()
		defer p.FinishStatus()

		fmt.Printf("Starting %s release pipeline...\n", mode)
		if dryRun {
			fmt.Println("DRY RUN — no changes will be made")
		}
		fmt.Println()

		stages := []struct {
			name string
			fn   func() *releasepipeline.StageResult
		}{
			{"validate", p.Validate},
			{"build", p.Build},
			{"sign", p.Sign},
			{"notes", p.GenerateNotes},
			{"upload", p.Upload},
			{"github release", p.CreateGitHubRelease},
			{"report", func() *releasepipeline.StageResult { return p.GenerateReport(started) }},
		}

		if notesOnly {
			stages = []struct {
				name string
				fn   func() *releasepipeline.StageResult
			}{{name: "notes", fn: p.GenerateNotes}}
		}

		for i, s := range stages {
			p.UpdateStatus(releasepipeline.PipelineStage(s.name), fmt.Sprintf("stage %d/%d", i+1, len(stages)))
			fmt.Printf("[%d/%d] %s... ", i+1, len(stages), s.name)
			result := s.fn()
			if result.Success {
				fmt.Println("✓")
			} else {
				fmt.Println("✗")
				fmt.Printf("  Error: %s\n", result.Message)
				if result.Error != "" {
					fmt.Printf("  Details: %s\n", result.Error)
				}
			}
		}

		fmt.Println()

		if notesOnly {
			p.GenerateReport(started)
		}

		fmt.Println("Pipeline stages:")
		for _, r := range p.Results() {
			icon := "✓"
			if !r.Success {
				icon = "✗"
			}
			fmt.Printf("  %s %s — %s\n", icon, r.Stage, r.Message)
		}

		fmt.Println()
		if !dryRun {
			fmt.Printf("Release mode: %s\n", mode)
			if jsonOut {
				fmt.Println("JSON output requested (TODO: implement structured output)")
			}
		}

		return nil
	},
}

var releaseStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current release pipeline status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(releasepipeline.StatusSummary())
		return nil
	},
}

var releaseNotesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Generate release notes from git history",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getReleaseSession()

		format, _ := cmd.Flags().GetString("format")
		notes, err := session.GenerateNotes()
		if err != nil {
			return fmt.Errorf("release notes generation failed: %w", err)
		}

		var rf release.NotesFormat
		switch format {
		case "json":
			rf = release.NotesJSON
		case "html":
			rf = release.NotesHTML
		default:
			rf = release.NotesMarkdown
		}

		savePath, err := session.SaveNotes(notes, rf)
		if err != nil {
			return fmt.Errorf("failed to save release notes: %w", err)
		}

		fmt.Println("Release notes generated!")
		fmt.Printf("  Path: %s\n", savePath)
		return nil
	},
}

var releaseHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View release history",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getReleaseSession()

		version, _ := cmd.Flags().GetString("version")

		if version != "" {
			entry, err := session.GetHistoryEntry(version)
			if err != nil {
				return fmt.Errorf("version not found: %w", err)
			}
			fmt.Printf("Version: %s (build %s)\n", entry.Version, entry.BuildNumber)
			fmt.Printf("  Date:   %s\n", entry.Date.Format("2006-01-02"))
			fmt.Printf("  Status: %s\n", entry.Status)
			fmt.Printf("  Notes:  %s\n", entry.Notes)
			return nil
		}

		history, err := session.GetHistory()
		if err != nil {
			return fmt.Errorf("failed to get history: %w", err)
		}

		fmt.Println("Release History")
		fmt.Println("===============")
		for _, e := range history.Entries {
			fmt.Printf("  v%s (build %s)\n", e.Version, e.BuildNumber)
			fmt.Printf("    Date:   %s\n", e.Date.Format("2006-01-02"))
			fmt.Printf("    Status: %s\n", e.Status)
		}
		fmt.Printf("\nTotal: %d release(s)\n", history.Total)
		return nil
	},
}

var releasePrepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Validate everything is ready for a release",
	RunE: func(cmd *cobra.Command, args []string) error {
		session := getReleaseSession()

		result, err := session.Prepare()
		if err != nil {
			return fmt.Errorf("release preparation failed: %w", err)
		}

		fmt.Println("Release Preparation Report")
		fmt.Println("==========================")
		fmt.Printf("  Signing:       %s\n", passFail(result.SigningOK))
		fmt.Printf("  Build:         %s\n", passFail(result.BuildOK))
		fmt.Printf("  IPA:           %s\n", passFail(result.IPAOK))
		fmt.Printf("  Metadata:      %s\n", passFail(result.MetadataOK))
		fmt.Printf("  Release Notes: %s\n", passFail(result.NotesOK))
		fmt.Printf("  Git State:     %s\n", passFail(result.GitStateOK))
		fmt.Println()

		if len(result.Issues) > 0 {
			fmt.Println("Issues:")
			for _, issue := range result.Issues {
				fmt.Printf("  ✗ %s\n", issue)
			}
		}
		if len(result.Warnings) > 0 {
			fmt.Println("Warnings:")
			for _, warn := range result.Warnings {
				fmt.Printf("  ⚠ %s\n", warn)
			}
		}

		fmt.Println()
		if result.Success {
			fmt.Println("Ready for release!")
		} else {
			fmt.Println("Fix the issues above before releasing.")
		}
		return nil
	},
}

func passFail(ok bool) string {
	if ok {
		return "✓ Pass"
	}
	return "✗ Fail"
}

func init() {
	rootCmd.AddCommand(testflightCmd)
	testflightCmd.AddCommand(testflightUploadCmd)
	testflightCmd.AddCommand(testflightStatusCmd)
	testflightCmd.AddCommand(testflightGroupsCmd)
	testflightCmd.AddCommand(testflightBuildsCmd)
	testflightCmd.AddCommand(testflightTestersCmd)

	rootCmd.AddCommand(releaseCmd)
	releaseCmd.AddCommand(releaseStatusCmd)
	releaseCmd.AddCommand(releaseNotesCmd)
	releaseCmd.AddCommand(releaseHistoryCmd)
	releaseCmd.AddCommand(releasePrepareCmd)

	releaseCmd.Flags().Bool("beta", false, "Beta release mode")
	releaseCmd.Flags().Bool("production", false, "Production release mode")
	releaseCmd.Flags().Bool("dry-run", false, "Preview release without changes")
	releaseCmd.Flags().Bool("notes", false, "Generate release notes only")
	releaseCmd.Flags().Bool("json", false, "Output in JSON format")

	testflightUploadCmd.Flags().String("artifact", "", "Path to IPA artifact")
	testflightUploadCmd.Flags().String("build", "", "Build number to upload")

	testflightStatusCmd.Flags().String("build", "", "Check status of specific build")

	testflightGroupsCmd.Flags().String("inspect", "", "Inspect a specific group")

	testflightBuildsCmd.Flags().String("build", "", "Get details of specific build")

	releaseNotesCmd.Flags().String("format", "markdown", "Output format: markdown, json, or html")

	releaseHistoryCmd.Flags().String("version", "", "Get details of specific version")
}
