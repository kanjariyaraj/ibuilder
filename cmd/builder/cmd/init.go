package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	initpkg "github.com/kanjariyaraj/Builder/internal/init"
	"github.com/kanjariyaraj/Builder/internal/logger"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project for Builder",
	Long:  "Detect project type, generate builder.json configuration, and create GitHub Actions workflow.",
}

var initRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run project initialization",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		yes, _ := cmd.Flags().GetBool("yes")
		jsonOut, _ := cmd.Flags().GetBool("json")

		dir, _ := os.Getwd()
		if len(args) > 0 {
			dir = args[0]
		}

		log := logger.New(logger.LevelInfo)

		opts := &initpkg.InitOptions{
			Force:  force,
			DryRun: dryRun,
			Yes:    yes,
			JSON:   jsonOut,
			Dir:    dir,
		}

		result, err := initpkg.Run(log, opts)
		if err != nil {
			return err
		}

		if jsonOut {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
		} else if !dryRun {
			fmt.Println()
			fmt.Println("Initialization complete!")
			fmt.Printf("  Config:  %s\n", result.ConfigPath)
			if result.WorkflowPath != "" {
				fmt.Printf("  Workflow: %s\n", result.WorkflowPath)
			}
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Println("  1. Review builder.json")
			fmt.Println("  2. Run 'builder auth github' to authenticate")
			fmt.Println("  3. Run 'builder repo connect' to connect repository")
			fmt.Println("  4. Run 'builder build run <workflow-id>' to build")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initRunCmd)
	initRunCmd.Flags().Bool("force", false, "Overwrite existing files without prompting")
	initRunCmd.Flags().Bool("dry-run", false, "Show what would be done without making changes")
	initRunCmd.Flags().Bool("yes", false, "Answer yes to all prompts (for CI)")
	initRunCmd.Flags().Bool("json", false, "Output in JSON format")
}
