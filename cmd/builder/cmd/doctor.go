package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type checkResult struct {
	name    string
	status  string
	message string
}

func runCheck(name, cmdName string, args ...string) checkResult {
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.Output()
	if err != nil {
		return checkResult{name: name, status: "FAILURE", message: fmt.Sprintf("not found: %v", err)}
	}
	return checkResult{name: name, status: "HEALTHY", message: strings.TrimSpace(string(output))}
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system dependencies and configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Builder Doctor")
		fmt.Println("==============")
		fmt.Println()

		results := []checkResult{}

		results = append(results, runCheck("Go", "go", "version"))
		results = append(results, runCheck("Git", "git", "version"))
		results = append(results, runCheck("GitHub CLI", "gh", "version"))
		results = append(results, runCheck("Node", "node", "--version"))
		results = append(results, runCheck("Java", "java", "-version"))

		runtimeOS := runtime.GOOS
		if runtimeOS == "linux" || runtimeOS == "darwin" {
			results = append(results, checkResult{name: "Platform", status: "HEALTHY", message: runtimeOS})
		} else {
			results = append(results, checkResult{name: "Platform", status: "WARNING", message: runtimeOS})
		}

		summary := map[string]int{"HEALTHY": 0, "WARNING": 0, "FAILURE": 0}
		for _, r := range results {
			summary[r.status]++
			statusSymbol := ""
			switch r.status {
			case "HEALTHY":
				statusSymbol = "✓"
			case "WARNING":
				statusSymbol = "⚠"
			case "FAILURE":
				statusSymbol = "✗"
			}
			fmt.Printf("  %s [%s] %s\n", statusSymbol, r.status, r.name)
			if r.message != "" {
				fmt.Printf("         %s\n", r.message)
			}
		}

		fmt.Println()
		fmt.Printf("Summary: %d healthy, %d warnings, %d failures\n",
			summary["HEALTHY"], summary["WARNING"], summary["FAILURE"])
	},
}
