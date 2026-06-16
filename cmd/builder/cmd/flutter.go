package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/flutter"
	"github.com/kanjariyaraj/Builder/internal/logger"
	"github.com/spf13/cobra"
)

func getFlutterSession() (*flutter.Session, *config.Config, error) {
	path := cfgFile
	if path == "" {
		cwd, _ := os.Getwd()
		path = cwd + "/builder.json"
	}

	cfg, err := config.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	log := logger.New(logger.LevelInfo)
	session := flutter.NewSession(&cfg.Flutter, log)

	cwd, _ := os.Getwd()
	session.DetectFlutterProject(cwd)

	return session, cfg, nil
}

var flutterCmd = &cobra.Command{
	Use:   "flutter",
	Short: "Flutter development commands",
	Long:  "Develop, debug, and deploy Flutter iOS apps on real devices via MobAI.",
}

var flutterDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start Flutter development mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, cfg, err := getFlutterSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		if _, err := session.DetectFlutterProject(cwd); err != nil {
			return fmt.Errorf("not a Flutter project: %w", err)
		}

		device, _ := cmd.Flags().GetString("device")
		if device != "" {
			session.SetDeviceID(device)
		}

		autoInstall, _ := cmd.Flags().GetBool("install")
		if autoInstall || cfg.Flutter.AutoInstall {
			session.BuildAndInstall()
		}

		result, err := session.DevMode()
		if err != nil {
			return fmt.Errorf("dev mode failed: %w", err)
		}

		fmt.Println("Flutter dev mode started!")
		fmt.Printf("  PID:    %d\n", result.PID)
		fmt.Printf("  Device: %s\n", result.Device)
		fmt.Printf("  Time:   %s\n", result.Started.Format(time.RFC3339))

		return nil
	},
}

var flutterAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a running Flutter app",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getFlutterSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		if _, err := session.DetectFlutterProject(cwd); err != nil {
			return fmt.Errorf("not a Flutter project: %w", err)
		}

		device, _ := cmd.Flags().GetString("device")

		result, err := session.Attach(device)
		if err != nil {
			return fmt.Errorf("attach failed: %w", err)
		}

		fmt.Println("Attached to Flutter app!")
		fmt.Printf("  PID:    %d\n", result.PID)
		fmt.Printf("  Device: %s\n", result.Device)
		fmt.Printf("  URI:    %s\n", result.URI)

		return nil
	},
}

var flutterWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for file changes and auto-reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getFlutterSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		if _, err := session.DetectFlutterProject(cwd); err != nil {
			return fmt.Errorf("not a Flutter project: %w", err)
		}

		watcher, err := session.WatchMode()
		if err != nil {
			return fmt.Errorf("watch mode failed: %w", err)
		}
		defer watcher.Stop()

		fmt.Println("Watching for file changes... (Ctrl+C to stop)")
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		eventCount := 0
		for {
			select {
			case <-sigChan:
				fmt.Println("\nWatch mode stopped.")
				return nil
			case event := <-watcher.Events():
				eventCount++
				fmt.Printf("\n[%s] %d change(s) detected:\n",
					time.Now().Format(time.Stamp), len(event.Changes))
				for _, change := range event.Changes {
					fmt.Printf("  %s: %s\n", change.Operation, change.Path)
				}
			}
		}
	},
}

var flutterReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Trigger Flutter hot reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getFlutterSession()
		if err != nil {
			return err
		}

		result, err := session.HotReload()
		if err != nil {
			return fmt.Errorf("hot reload failed: %w", err)
		}

		fmt.Println("Hot reload completed!")
		fmt.Printf("  Duration: %v\n", result.Duration)
		fmt.Printf("  Output:   %s\n", result.Output)
		return nil
	},
}

var flutterRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Trigger Flutter hot restart",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getFlutterSession()
		if err != nil {
			return err
		}

		result, err := session.Restart()
		if err != nil {
			return fmt.Errorf("hot restart failed: %w", err)
		}

		fmt.Println("Hot restart completed!")
		fmt.Printf("  Duration: %v\n", result.Duration)
		fmt.Printf("  Output:   %s\n", result.Output)
		return nil
	},
}

var flutterLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View Flutter device logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getFlutterSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		session.DetectFlutterProject(cwd)

		stream, _ := cmd.Flags().GetBool("stream")
		save, _ := cmd.Flags().GetString("save")
		level, _ := cmd.Flags().GetString("level")
		search, _ := cmd.Flags().GetString("search")
		since, _ := cmd.Flags().GetDuration("since")

		filter := &flutter.LogFilter{
			Level:  level,
			Search: search,
			Since:  since,
		}

		if stream {
			stopChan := make(chan struct{})
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			logChan, err := session.StreamLogs(stopChan)
			if err != nil {
				return fmt.Errorf("failed to stream logs: %w", err)
			}

			fmt.Println("Streaming Flutter logs... (Ctrl+C to stop)")
			go func() {
				<-sigChan
				close(stopChan)
			}()

			for entry := range logChan {
				fmt.Printf("[%s] [%s] %s\n",
					entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Message)
			}
			return nil
		}

		logs, err := session.FetchLogs(filter)
		if err != nil {
			return fmt.Errorf("failed to fetch logs: %w", err)
		}

		if len(logs) == 0 {
			fmt.Println("No logs found.")
			return nil
		}

		if save != "" {
			path, err := session.SaveLogs(logs, save)
			if err != nil {
				return fmt.Errorf("failed to save logs: %w", err)
			}
			fmt.Printf("Logs saved to: %s\n", path)
			return nil
		}

		for _, l := range logs {
			fmt.Printf("[%s] [%s] %s\n",
				l.Timestamp.Format(time.RFC3339), l.Level, l.Message)
		}

		return nil
	},
}

var flutterDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check Flutter development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getFlutterSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		session.DetectFlutterProject(cwd)

		report := session.Doctor()

		fmt.Println("Flutter Doctor Report")
		fmt.Println("=====================")
		for _, check := range report.Checks {
			symbol := "✓"
			switch check.Status {
			case flutter.StatusWarning:
				symbol = "⚠"
			case flutter.StatusFailure:
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

func init() {
	rootCmd.AddCommand(flutterCmd)
	flutterCmd.AddCommand(flutterDevCmd)
	flutterCmd.AddCommand(flutterAttachCmd)
	flutterCmd.AddCommand(flutterWatchCmd)
	flutterCmd.AddCommand(flutterReloadCmd)
	flutterCmd.AddCommand(flutterRestartCmd)
	flutterCmd.AddCommand(flutterLogsCmd)
	flutterCmd.AddCommand(flutterDoctorCmd)

	flutterDevCmd.Flags().String("device", "", "Target device ID")
	flutterDevCmd.Flags().Bool("install", false, "Install before running")

	flutterAttachCmd.Flags().String("device", "", "Target device ID")

	flutterLogsCmd.Flags().Bool("stream", false, "Stream live logs")
	flutterLogsCmd.Flags().String("save", "", "Save logs to directory")
	flutterLogsCmd.Flags().String("level", "", "Filter by log level")
	flutterLogsCmd.Flags().String("search", "", "Search logs")
	flutterLogsCmd.Flags().Duration("since", 0, "Show logs since duration")
}
