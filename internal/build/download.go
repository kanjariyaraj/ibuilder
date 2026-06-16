package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kanjariyaraj/Builder/internal/github"
)

func downloadArtifact(client *github.Client, opts *BuildOptions, runID int64) (string, error) {
	artifacts, err := github.ListArtifacts(client, opts.Owner, opts.Name, runID)
	if err != nil {
		return "", fmt.Errorf("failed to list artifacts: %w", err)
	}

	if len(artifacts) == 0 {
		return "", fmt.Errorf("no artifacts found for run %d", runID)
	}

	artifact := artifacts[0]
	fmt.Printf("Downloading artifact: %s (%.1f KB)\n", artifact.Name, float64(artifact.Size)/1024)

	distDir := filepath.Join(opts.Dir, "dist")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create dist directory: %w", err)
	}

	data, err := github.DownloadArtifact(client, opts.Owner, opts.Name, artifact.ID)
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(distDir, artifact.Name+".zip")
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save artifact: %w", err)
	}

	fmt.Printf("Saved to: %s\n", outputPath)
	return outputPath, nil
}

func downloadLatest(client *github.Client, opts *BuildOptions) (*BuildResult, error) {
	runs, err := github.ListWorkflowRuns(client, opts.Owner, opts.Name, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to find latest build: %w", err)
	}
	if len(runs) == 0 {
		return nil, fmt.Errorf("no builds found")
	}

	run := runs[0]
	if run.Conclusion != "success" {
		return nil, fmt.Errorf("latest build (#%d) did not succeed: %s", run.RunNumber, run.Conclusion)
	}

	artifact, err := downloadArtifact(client, opts, run.ID)
	if err != nil {
		return nil, err
	}

	return &BuildResult{
		RunID:      run.ID,
		RunNumber:  run.RunNumber,
		Status:     run.Status,
		Conclusion: run.Conclusion,
		Artifact:   artifact,
	}, nil
}
