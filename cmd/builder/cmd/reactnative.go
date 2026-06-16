package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
	"github.com/kanjariyaraj/Builder/internal/reactnative"
	"github.com/spf13/cobra"
)

func getReactNativeSession() (*reactnative.Session, *config.Config, error) {
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
	session := reactnative.NewSession(&cfg.ReactNative, log)

	cwd, _ := os.Getwd()
	session.DetectRNProject(cwd)

	return session, cfg, nil
}

var rnCmd = &cobra.Command{
	Use:   "rn",
	Short: "React Native development commands",
	Long:  "Develop, debug, and deploy React Native iOS apps on real devices via MobAI.",
}

var rnDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start React Native development mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, cfg, err := getReactNativeSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		if _, err := session.DetectRNProject(cwd); err != nil {
			return fmt.Errorf("not a React Native project: %w", err)
		}

		device, _ := cmd.Flags().GetString("device")
		if device != "" {
			session.SetDeviceID(device)
		}

		autoInstall, _ := cmd.Flags().GetBool("install")
		if autoInstall || cfg.ReactNative.AutoInstall {
			session.BuildAndInstall()
		}

		result, err := session.DevMode()
		if err != nil {
			return fmt.Errorf("dev mode failed: %w", err)
		}

		fmt.Println("React Native dev mode started!")
		fmt.Printf("  PID:       %d\n", result.PID)
		fmt.Printf("  Device:    %s\n", result.Device)
		fmt.Printf("  Port:      %d\n", result.MetroPort)
		fmt.Printf("  Time:      %s\n", result.Started.Format(time.RFC3339))

		return nil
	},
}

var rnAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a running React Native app",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		if _, err := session.DetectRNProject(cwd); err != nil {
			return fmt.Errorf("not a React Native project: %w", err)
		}

		device, _ := cmd.Flags().GetString("device")

		result, err := session.Attach(device)
		if err != nil {
			return fmt.Errorf("attach failed: %w", err)
		}

		fmt.Println("Attached to React Native app!")
		fmt.Printf("  PID:       %d\n", result.PID)
		fmt.Printf("  Device:    %s\n", result.Device)
		fmt.Printf("  Port:      %d\n", result.MetroPort)

		return nil
	},
}

var rnMetroCmd = &cobra.Command{
	Use:   "metro",
	Short: "Manage Metro bundler",
	Long:  "Start, stop, restart, or check status of the Metro bundler.",
}

var rnMetroStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Metro bundler",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		session.DetectRNProject(cwd)

		port, _ := cmd.Flags().GetInt("port")
		if port > 0 {
			session.UpdateConfig(&config.ReactNativeSettings{MetroPort: port})
		}

		result, err := session.StartMetro()
		if err != nil {
			return fmt.Errorf("metro start failed: %w", err)
		}

		fmt.Println("Metro bundler started!")
		fmt.Printf("  PID:  %d\n", result.PID)
		fmt.Printf("  Port: %d\n", result.Port)
		return nil
	},
}

var rnMetroStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Metro bundler",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		_, err = session.StopMetro()
		if err != nil {
			return fmt.Errorf("metro stop failed: %w", err)
		}

		fmt.Println("Metro bundler stopped.")
		return nil
	},
}

var rnMetroRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart Metro bundler",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		session.DetectRNProject(cwd)

		result, err := session.RestartMetro()
		if err != nil {
			return fmt.Errorf("metro restart failed: %w", err)
		}

		fmt.Println("Metro bundler restarted!")
		fmt.Printf("  PID:  %d\n", result.PID)
		fmt.Printf("  Port: %d\n", result.Port)
		return nil
	},
}

var rnMetroStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Metro bundler status",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		status := session.MetroStatus()
		if status.Success {
			fmt.Printf("Metro is RUNNING (PID: %d, Port: %d)\n", status.PID, status.Port)
		} else {
			fmt.Println("Metro is STOPPED")
		}
		return nil
	},
}

var rnReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Trigger React Native Fast Refresh or manual reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		refresh, _ := cmd.Flags().GetBool("fast-refresh")
		if refresh {
			result, err := session.FastRefresh()
			if err != nil {
				return fmt.Errorf("fast refresh failed: %w", err)
			}
			fmt.Println("Fast refresh triggered!")
			fmt.Printf("  Duration: %v\n", result.Duration)
			return nil
		}

		result, err := session.ManualReload()
		if err != nil {
			return fmt.Errorf("manual reload failed: %w", err)
		}

		fmt.Println("Manual reload triggered!")
		fmt.Printf("  Duration: %v\n", result.Duration)
		return nil
	},
}

var rnLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View React Native device logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		session.DetectRNProject(cwd)

		stream, _ := cmd.Flags().GetBool("stream")
		save, _ := cmd.Flags().GetString("save")
		level, _ := cmd.Flags().GetString("level")
		search, _ := cmd.Flags().GetString("search")
		since, _ := cmd.Flags().GetDuration("since")

		filter := &reactnative.LogFilter{
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

			fmt.Println("Streaming React Native logs... (Ctrl+C to stop)")
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

var rnDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check React Native development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		session.DetectRNProject(cwd)

		report := session.Doctor()

		fmt.Println("React Native Doctor Report")
		fmt.Println("==========================")
		for _, check := range report.Checks {
			symbol := "✓"
			switch check.Status {
			case reactnative.StatusWarning:
				symbol = "⚠"
			case reactnative.StatusFailure:
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

var rnInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install React Native app on device",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getReactNativeSession()
		if err != nil {
			return err
		}

		cwd, _ := os.Getwd()
		if _, err := session.DetectRNProject(cwd); err != nil {
			return fmt.Errorf("not a React Native project: %w", err)
		}

		device, _ := cmd.Flags().GetString("device")
		if device != "" {
			session.SetDeviceID(device)
		}

		artifact, _ := cmd.Flags().GetString("artifact")

		var result *reactnative.InstallResult
		if artifact != "" {
			result, err = session.InstallArtifact(artifact)
		} else {
			result, err = session.InstallLatest()
		}
		if err != nil {
			return fmt.Errorf("install failed: %w", err)
		}

		fmt.Println("React Native app installed!")
		fmt.Printf("  Artifact: %s\n", result.Artifact)
		fmt.Printf("  Device:   %s\n", result.Device)
		fmt.Printf("  Time:     %s\n", result.Installed.Format(time.RFC3339))

		verified, _ := session.VerifyInstall()
		if verified {
			fmt.Println("  Status:   VERIFIED")
		} else {
			fmt.Println("  Status:   COULD NOT VERIFY")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rnCmd)
	rnCmd.AddCommand(rnDevCmd)
	rnCmd.AddCommand(rnAttachCmd)
	rnCmd.AddCommand(rnMetroCmd)
	rnCmd.AddCommand(rnReloadCmd)
	rnCmd.AddCommand(rnLogsCmd)
	rnCmd.AddCommand(rnDoctorCmd)
	rnCmd.AddCommand(rnInstallCmd)

	rnMetroCmd.AddCommand(rnMetroStartCmd)
	rnMetroCmd.AddCommand(rnMetroStopCmd)
	rnMetroCmd.AddCommand(rnMetroRestartCmd)
	rnMetroCmd.AddCommand(rnMetroStatusCmd)

	rnDevCmd.Flags().String("device", "", "Target device ID")
	rnDevCmd.Flags().Bool("install", false, "Install before running")

	rnAttachCmd.Flags().String("device", "", "Target device ID")

	rnMetroStartCmd.Flags().Int("port", 8081, "Metro bundler port")

	rnReloadCmd.Flags().Bool("fast-refresh", false, "Trigger fast refresh instead of manual reload")

	rnLogsCmd.Flags().Bool("stream", false, "Stream live logs")
	rnLogsCmd.Flags().String("save", "", "Save logs to directory")
	rnLogsCmd.Flags().String("level", "", "Filter by log level")
	rnLogsCmd.Flags().String("search", "", "Search logs")
	rnLogsCmd.Flags().Duration("since", 0, "Show logs since duration")

	rnInstallCmd.Flags().String("device", "", "Target device ID")
	rnInstallCmd.Flags().String("artifact", "", "Path to custom build artifact")
}
