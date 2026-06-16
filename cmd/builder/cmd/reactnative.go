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

func getRNSession() (*reactnative.Session, *config.Config, error) {
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
		session, cfg, err := getRNSession()
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
		fmt.Printf("  Metro Port: %d\n", result.MetroPort)
		fmt.Printf("  Time:      %s\n", result.Started.Format(time.RFC3339))

		return nil
	},
}

var rnAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a running React Native app",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getRNSession()
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
		fmt.Printf("  Metro Port: %d\n", result.MetroPort)

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
		session, _, err := getRNSession()
		if err != nil {
			return err
		}

		if err := session.StartMetro(); err != nil {
			return fmt.Errorf("metro start failed: %w", err)
		}

		status := session.MetroStatus()
		fmt.Println("Metro bundler started!")
		fmt.Printf("  PID:  %d\n", status.PID)
		fmt.Printf("  Port: %d\n", status.Port)
		fmt.Printf("  Host: %s\n", status.Host)
		return nil
	},
}

var rnMetroStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Metro bundler",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getRNSession()
		if err != nil {
			return err
		}

		if err := session.StopMetro(); err != nil {
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
		session, _, err := getRNSession()
		if err != nil {
			return err
		}

		if err := session.RestartMetro(); err != nil {
			return fmt.Errorf("metro restart failed: %w", err)
		}

		status := session.MetroStatus()
		fmt.Println("Metro bundler restarted!")
		fmt.Printf("  PID:  %d\n", status.PID)
		fmt.Printf("  Port: %d\n", status.Port)
		return nil
	},
}

var rnMetroStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Metro bundler status",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getRNSession()
		if err != nil {
			return err
		}

		status := session.MetroStatus()
		if status.Running {
			fmt.Printf("Metro bundler is RUNNING\n")
			fmt.Printf("  PID:    %d\n", status.PID)
			fmt.Printf("  Port:   %d\n", status.Port)
			fmt.Printf("  Host:   %s\n", status.Host)
			if status.Uptime != "" {
				fmt.Printf("  Uptime: %s\n", status.Uptime)
			}
		} else {
			fmt.Println("Metro bundler is NOT running")
		}
		return nil
	},
}

var rnReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Trigger Fast Refresh or manual reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getRNSession()
		if err != nil {
			return err
		}

		mode, _ := cmd.Flags().GetString("mode")

		var result *reactnative.ReloadResult
		switch mode {
		case "manual":
			result, err = session.ManualReload()
		case "device":
			result, err = session.DeviceRefresh()
		default:
			result, err = session.FastRefresh()
		}

		if err != nil {
			return fmt.Errorf("reload failed: %w", err)
		}

		fmt.Println("Reload completed!")
		fmt.Printf("  Type:     %s\n", result.Type)
		fmt.Printf("  Duration: %v\n", result.Duration)
		fmt.Printf("  Output:   %s\n", result.Output)
		return nil
	},
}

var rnLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View React Native device logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, _, err := getRNSession()
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
		session, _, err := getRNSession()
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
		session, _, err := getRNSession()
		if err != nil {
			return err
		}

		artifact, _ := cmd.Flags().GetString("artifact")

		var result *reactnative.InstallResult
		if artifact != "" {
			result, err = session.InstallSpecificBuild(artifact)
		} else {
			result, err = session.InstallLatest()
		}

		if err != nil {
			return fmt.Errorf("install failed: %w", err)
		}

		fmt.Println("Install completed!")
		fmt.Printf("  Action: %s\n", result.Action)
		fmt.Printf("  Output: %s\n", result.Output)

		if err := session.VerifyInstallation(); err != nil {
			fmt.Printf("  Warning: %v\n", err)
		} else {
			fmt.Println("  Installation verified.")
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

	rnReloadCmd.Flags().String("mode", "fast", "Reload mode: fast, manual, device")

	rnLogsCmd.Flags().Bool("stream", false, "Stream live logs")
	rnLogsCmd.Flags().String("save", "", "Save logs to directory")
	rnLogsCmd.Flags().String("level", "", "Filter by log level")
	rnLogsCmd.Flags().String("search", "", "Search logs")
	rnLogsCmd.Flags().Duration("since", 0, "Show logs since duration")

	rnInstallCmd.Flags().String("artifact", "", "Path to specific build artifact")
}
