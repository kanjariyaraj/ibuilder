package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/kanjariyaraj/Builder/internal/config"
	"github.com/kanjariyaraj/Builder/internal/logger"
	"github.com/kanjariyaraj/Builder/internal/mobai"
	"github.com/spf13/cobra"
)

func getMobaiClient() (*mobai.Client, *config.Config, error) {
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

var mobaiCmd = &cobra.Command{
	Use:   "mobai",
	Short: "Manage MobAI device connections",
	Long:  "Connect, manage, and monitor real iPhone devices via MobAI.",
}

var mobaiConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to MobAI device",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, cfg, err := getMobaiClient()
		if err != nil {
			return err
		}

		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		device, _ := cmd.Flags().GetString("device")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if host != "" {
			cfg.Mobai.Host = host
		}
		if port > 0 {
			cfg.Mobai.Port = port
		}
		if device != "" {
			cfg.Mobai.Device = device
		}
		if timeout > 0 {
			cfg.Mobai.ConnectionTimeout = timeout
		}

		client.UpdateConfig(&cfg.Mobai)

		fmt.Println("Connecting to MobAI...")
		if err := client.Connect(); err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}

		path := cfgFile
		if path == "" {
			cwd, _ := os.Getwd()
			path = cwd + "/builder.json"
		}
		if err := config.Save(path, cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println("Connected successfully!")
		return nil
	},
}

var mobaiDisconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect from MobAI device",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getMobaiClient()
		if err != nil {
			return err
		}

		if err := client.Disconnect(); err != nil {
			return fmt.Errorf("disconnect failed: %w", err)
		}

		fmt.Println("Disconnected.")
		return nil
	},
}

var mobaiStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show MobAI connection status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getMobaiClient()
		if err != nil {
			return err
		}

		status := client.Status()

		fmt.Printf("State:       %s\n", status.State)
		if status.Device != nil {
			fmt.Printf("Device:      %s\n", status.Device.Name)
			fmt.Printf("Model:       %s\n", status.Device.Model)
			fmt.Printf("iOS:         %s\n", status.Device.OSVersion)
		}
		fmt.Printf("Latency:     %dms\n", status.Latency.Milliseconds())
		if !status.ConnectedAt.IsZero() {
			fmt.Printf("Connected:   %s\n", status.ConnectedAt.Format(time.RFC3339))
		}
		if status.Error != "" {
			fmt.Printf("Error:       %s\n", status.Error)
		}

		return nil
	},
}

var mobaiDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run MobAI health checks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getMobaiClient()
		if err != nil {
			return err
		}

		report := client.Doctor()

		fmt.Println("MobAI Doctor Report")
		fmt.Println("===================")
		for _, check := range report.Checks {
			symbol := "✓"
			switch check.Status {
			case mobai.StatusWarning:
				symbol = "⚠"
			case mobai.StatusFailure:
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

var mobaiPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping MobAI device connection",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getMobaiClient()
		if err != nil {
			return err
		}

		latency, err := client.Ping()
		if err != nil {
			return fmt.Errorf("ping failed: %w", err)
		}

		fmt.Printf("Pong! Latency: %dms\n", latency.Milliseconds())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mobaiCmd)
	mobaiCmd.AddCommand(mobaiConnectCmd)
	mobaiCmd.AddCommand(mobaiDisconnectCmd)
	mobaiCmd.AddCommand(mobaiStatusCmd)
	mobaiCmd.AddCommand(mobaiDoctorCmd)
	mobaiCmd.AddCommand(mobaiPingCmd)

	mobaiConnectCmd.Flags().String("host", "", "MobAI host address")
	mobaiConnectCmd.Flags().Int("port", 0, "MobAI port")
	mobaiConnectCmd.Flags().String("device", "", "Device name")
	mobaiConnectCmd.Flags().Int("timeout", 30, "Connection timeout in seconds")
}
