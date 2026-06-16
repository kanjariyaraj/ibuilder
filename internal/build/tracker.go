package build

import (
	"fmt"
	"time"

	"github.com/kanjariyaraj/Builder/internal/github"
)

type statusEvent struct {
	Status     string
	Conclusion string
}

func waitForCompletion(client *github.Client, opts *BuildOptions, runID int64) (*BuildResult, error) {
	fmt.Println("Waiting for build to complete...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var lastStatus string

	for range ticker.C {
		run, err := github.GetWorkflowRun(client, opts.Owner, opts.Name, runID)
		if err != nil {
			return nil, fmt.Errorf("failed to get run status: %w", err)
		}

		if run.Status != lastStatus {
			statusSymbol := ""
			switch run.Status {
			case "queued":
				statusSymbol = "○"
			case "in_progress":
				statusSymbol = "◌"
			case "completed":
				statusSymbol = run.Conclusion
				switch run.Conclusion {
				case "success":
					statusSymbol = "✓"
				case "failure":
					statusSymbol = "✗"
				case "cancelled":
					statusSymbol = "⊘"
				default:
					statusSymbol = "?"
				}
			}
			fmt.Printf("\r  %s %s", statusSymbol, run.Status)
			lastStatus = run.Status
		}

		if run.Status == "completed" {
			fmt.Println()
			fmt.Printf("Conclusion: %s\n", run.Conclusion)

			result := &BuildResult{
				RunID:       run.ID,
				RunNumber:   run.RunNumber,
				Status:      run.Status,
				Conclusion:  run.Conclusion,
				WorkflowURL: run.HTMLURL,
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("timed out waiting for build completion")
}
