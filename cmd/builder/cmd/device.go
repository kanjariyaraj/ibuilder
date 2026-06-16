package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
	"github.com/kanjariyaraj/Builder/internal/mobai"
	"github.com/spf13/cobra"
)

func getDeviceClient() (*mobai.Client, *config.Config, error) {
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
	client := mobai.NewClient(&cfg.Mobai, log)
	return client, cfg, nil
}

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Manage connected iOS devices",
	Long:  "List, inspect, install, launch, and debug iOS devices.",
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List connected devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getDeviceClient()
		if err != nil {
			return err
		}

		if err := client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer client.Disconnect()

		devices, err := client.ListDevices()
		if err != nil {
			return fmt.Errorf("failed to list devices: %w", err)
		}

		if len(devices) == 0 {
			fmt.Println("No devices found.")
			return nil
		}

		asJSON, _ := cmd.Flags().GetBool("json")
		if asJSON {
			return printJSON(devices)
		}

		mobai.PrintDeviceTable(devices)
		return nil
	},
}

var deviceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show device information",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getDeviceClient()
		if err != nil {
			return err
		}

		if err := client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer client.Disconnect()

		udid, _ := cmd.Flags().GetString("udid")

		device, err := client.DeviceInfo(udid)
		if err != nil {
			return fmt.Errorf("failed to get device info: %w", err)
		}

		asJSON, _ := cmd.Flags().GetBool("json")
		if asJSON {
			return printJSON(device)
		}

		fmt.Print(mobai.FormatDeviceInfo(device))
		return nil
	},
}

var deviceLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View device logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getDeviceClient()
		if err != nil {
			return err
		}

		if err := client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer client.Disconnect()

		stream, _ := cmd.Flags().GetBool("stream")
		save, _ := cmd.Flags().GetString("save")
		process, _ := cmd.Flags().GetString("process")
		level, _ := cmd.Flags().GetString("level")
		search, _ := cmd.Flags().GetString("search")
		since, _ := cmd.Flags().GetDuration("since")

		filter := &mobai.LogFilter{
			Process: process,
			Level:   level,
			Search:  search,
			Since:   since,
		}

		if stream {
			stopChan := make(chan struct{})
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			logChan, err := client.StreamLogs(filter, stopChan)
			if err != nil {
				return fmt.Errorf("failed to stream logs: %w", err)
			}

			fmt.Println("Streaming logs... (Ctrl+C to stop)")
			go func() {
				<-sigChan
				close(stopChan)
			}()

			for entry := range logChan {
				fmt.Printf("[%s] [%s] [%s] %s\n",
					entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Process, entry.Message)
			}
			return nil
		}

		logs, err := client.FetchLogs(filter)
		if err != nil {
			return fmt.Errorf("failed to fetch logs: %w", err)
		}

		if len(logs) == 0 {
			fmt.Println("No logs found.")
			return nil
		}

		if save != "" {
			path, err := client.SaveLogs(logs, save)
			if err != nil {
				return fmt.Errorf("failed to save logs: %w", err)
			}
			fmt.Printf("Logs saved to: %s\n", path)
			return nil
		}

		saveDefault, _ := cmd.Flags().GetBool("save-default")
		if saveDefault {
			path, err := client.SaveLogs(logs, ".build/logs")
			if err != nil {
				return fmt.Errorf("failed to save logs: %w", err)
			}
			fmt.Printf("Logs saved to: %s\n", path)
			return nil
		}

		for _, l := range logs {
			fmt.Printf("[%s] [%s] [%s] %s\n",
				l.Timestamp.Format(time.RFC3339), l.Level, l.Process, l.Message)
		}

		return nil
	},
}

var deviceScreenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Capture device screenshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getDeviceClient()
		if err != nil {
			return err
		}

		if err := client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer client.Disconnect()

		outputDir, _ := cmd.Flags().GetString("output")
		if outputDir == "" {
			outputDir = "screenshots"
		}

		result, err := client.CaptureScreenshot(outputDir)
		if err != nil {
			return fmt.Errorf("screenshot failed: %w", err)
		}

		fmt.Printf("Screenshot saved to: %s\n", result.Path)
		fmt.Printf("Size: %d bytes\n", result.Size)
		fmt.Printf("Device: %s\n", result.Device)
		fmt.Printf("Dimensions: %dx%d\n", result.Width, result.Height)

		return nil
	},
}

var deviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install app on device",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getDeviceClient()
		if err != nil {
			return err
		}

		if err := client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer client.Disconnect()

		ipaPath, _ := cmd.Flags().GetString("ipa")
		artifactPath, _ := cmd.Flags().GetString("artifact")

		var result *mobai.InstallResult

		if ipaPath != "" {
			result, err = client.InstallIPA(ipaPath)
		} else if artifactPath != "" {
			result, err = client.InstallLatestArtifact(artifactPath)
		} else {
			return fmt.Errorf("specify --ipa <path> or --artifact <path>")
		}

		if err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		fmt.Println("Installation successful!")
		fmt.Printf("  App:      %s\n", result.AppName)
		fmt.Printf("  Bundle:   %s\n", result.BundleID)
		fmt.Printf("  Version:  %s\n", result.Version)
		fmt.Printf("  Duration: %v\n", result.Duration)

		return nil
	},
}

var deviceLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch app on device",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getDeviceClient()
		if err != nil {
			return err
		}

		if err := client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer client.Disconnect()

		bundleID, _ := cmd.Flags().GetString("bundle-id")
		if bundleID == "" {
			return fmt.Errorf("--bundle-id is required")
		}

		result, err := client.LaunchApp(bundleID)
		if err != nil {
			return fmt.Errorf("launch failed: %w", err)
		}

		fmt.Println("Launch successful!")
		fmt.Printf("  Bundle:   %s\n", result.BundleID)
		fmt.Printf("  PID:      %d\n", result.PID)
		fmt.Printf("  Duration: %v\n", result.Duration)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceInfoCmd)
	deviceCmd.AddCommand(deviceLogsCmd)
	deviceCmd.AddCommand(deviceScreenshotCmd)
	deviceCmd.AddCommand(deviceInstallCmd)
	deviceCmd.AddCommand(deviceLaunchCmd)

	deviceListCmd.Flags().Bool("json", false, "JSON output")

	deviceInfoCmd.Flags().String("udid", "", "Device UDID")
	deviceInfoCmd.Flags().Bool("json", false, "JSON output")

	deviceLogsCmd.Flags().Bool("stream", false, "Stream live logs")
	deviceLogsCmd.Flags().String("save", "", "Save logs to directory")
	deviceLogsCmd.Flags().Bool("save-default", false, "Save logs to default location (.build/logs/)")
	deviceLogsCmd.Flags().String("process", "", "Filter by process name")
	deviceLogsCmd.Flags().String("level", "", "Filter by log level (DEBUG, INFO, WARN, ERROR)")
	deviceLogsCmd.Flags().String("search", "", "Search logs")
	deviceLogsCmd.Flags().Duration("since", 0, "Show logs since duration (e.g. 5m, 1h)")

	deviceScreenshotCmd.Flags().String("output", "screenshots", "Output directory")

	deviceInstallCmd.Flags().String("ipa", "", "Path to .ipa file")
	deviceInstallCmd.Flags().String("artifact", "", "Path to artifact file")

	deviceLaunchCmd.Flags().String("bundle-id", "", "Bundle identifier")
}
