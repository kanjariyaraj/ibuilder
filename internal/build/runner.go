package build

import (
	"fmt"
	"os"

	"github.com/kanjariyaraj/Builder/internal/github"
)

type BuildOptions struct {
	WorkflowID   string
	Branch       string
	Scheme       string
	Mode         string
	Wait         bool
	Logs         bool
	JSON         bool
	DownloadOnly bool
	Clean        bool
	Token        string
	Owner        string
	Name         string
	Dir          string
}

type BuildResult struct {
	RunID      int64  `json:"run_id"`
	RunNumber  int    `json:"run_number"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	WorkflowURL string `json:"workflow_url"`
	Artifact   string `json:"artifact,omitempty"`
	ReportPath string `json:"report_path,omitempty"`
}

type BuildReport struct {
	RunID       int64  `json:"run_id"`
	RunNumber   int    `json:"run_number"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	Duration    string `json:"duration"`
	WorkflowURL string `json:"workflow_url"`
	Artifact    string `json:"artifact"`
	Error       string `json:"error,omitempty"`
}

func RunBuild(opts *BuildOptions) (*BuildResult, error) {
	client := github.NewClient(opts.Token)

	if opts.DownloadOnly {
		return downloadLatest(client, opts)
	}

	if err := validatePrereqs(opts); err != nil {
		return nil, err
	}

	inputs := make(map[string]string)
	if opts.Scheme != "" {
		inputs["scheme"] = opts.Scheme
	}
	if opts.Mode != "" {
		inputs["build_mode"] = opts.Mode
	}

	if err := github.DispatchWorkflow(client, opts.Owner, opts.Name, opts.WorkflowID, opts.Branch, inputs); err != nil {
		return nil, fmt.Errorf("failed to dispatch workflow: %w", err)
	}

	runs, err := github.ListWorkflowRuns(client, opts.Owner, opts.Name, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow run: %w", err)
	}
	if len(runs) == 0 {
		return nil, fmt.Errorf("no workflow run found after dispatch")
	}

	run := runs[0]
	result := &BuildResult{
		RunID:       run.ID,
		RunNumber:   run.RunNumber,
		Status:      run.Status,
		WorkflowURL: run.HTMLURL,
	}

	fmt.Printf("Build #%d triggered\n", run.RunNumber)
	fmt.Printf("Status: %s\n", run.Status)
	fmt.Printf("URL:    %s\n", run.HTMLURL)

	if opts.Wait {
		result, err = waitForCompletion(client, opts, run.ID)
		if err != nil {
			return result, err
		}

		if result.Conclusion == "success" {
			artifact, err := downloadArtifact(client, opts, run.ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: artifact download failed: %v\n", err)
			} else {
				result.Artifact = artifact
			}
		}

		reportPath := generateReport(result)
		result.ReportPath = reportPath
	}

	return result, nil
}

func validatePrereqs(opts *BuildOptions) error {
	if opts.Token == "" {
		return fmt.Errorf("not authenticated - run 'builder auth github' first")
	}
	if opts.Owner == "" || opts.Name == "" {
		return fmt.Errorf("no repository configured - run 'builder repo connect' first")
	}
	if opts.WorkflowID == "" {
		return fmt.Errorf("workflow ID is required")
	}
	return nil
}
