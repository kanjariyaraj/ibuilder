package artifacts

import (
	"encoding/json"
	"fmt"
)

type BuildInspect struct {
	RunID       int64         `json:"run_id"`
	RunNumber   int           `json:"run_number"`
	Status      string        `json:"status"`
	Conclusion  string        `json:"conclusion"`
	Branch      string        `json:"branch"`
	CommitSHA   string        `json:"commit_sha"`
	CommitMsg   string        `json:"commit_msg"`
	Author      string        `json:"author"`
	Workflow    string        `json:"workflow"`
	URL         string        `json:"url"`
	Jobs        []JobInfo     `json:"jobs"`
	Artifacts   []Artifact    `json:"artifacts"`
	Duration    string        `json:"duration"`
}

type JobInfo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Conclusion string `json:"conclusion"`
	Steps     int    `json:"steps"`
}

func (m *ArtifactManager) InspectBuild(runID int64) (*BuildInspect, error) {
	data, err := m.client.Get(fmt.Sprintf("/repos/%s/%s/actions/runs/%d", m.owner, m.repo, runID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch build details: %w", err)
	}

	var run struct {
		ID          int64  `json:"id"`
		RunNumber   int    `json:"run_number"`
		Status      string `json:"status"`
		Conclusion  string `json:"conclusion"`
		HeadBranch  string `json:"head_branch"`
		HeadSHA     string `json:"head_sha"`
		Name        string `json:"name"`
		HTMLURL     string `json:"html_url"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		RunDuration int    `json:"run_duration_ms"`
		HeadCommit  struct {
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"head_commit"`
	}

	if err := json.Unmarshal(data, &run); err != nil {
		return nil, fmt.Errorf("failed to parse build details: %w", err)
	}

	artifacts, _ := m.ListWorkflowArtifacts(runID)

	jobs := m.fetchJobs(runID)

	duration := ""
	if run.RunDuration > 0 {
		duration = fmt.Sprintf("%ds", run.RunDuration/1000)
	}

	return &BuildInspect{
		RunID:      run.ID,
		RunNumber:  run.RunNumber,
		Status:     run.Status,
		Conclusion: run.Conclusion,
		Branch:     run.HeadBranch,
		CommitSHA:  run.HeadSHA,
		CommitMsg:  run.HeadCommit.Message,
		Author:     run.HeadCommit.Author.Name,
		Workflow:   run.Name,
		URL:        run.HTMLURL,
		Jobs:       jobs,
		Artifacts:  artifacts,
		Duration:   duration,
	}, nil
}

func (m *ArtifactManager) fetchJobs(runID int64) []JobInfo {
	data, err := m.client.Get(fmt.Sprintf("/repos/%s/%s/actions/runs/%d/jobs", m.owner, m.repo, runID))
	if err != nil {
		return nil
	}

	var resp struct {
		Jobs []struct {
			ID         int64  `json:"id"`
			Name       string `json:"name"`
			Status     string `json:"status"`
			Conclusion string `json:"conclusion"`
			Steps      []struct{} `json:"steps"`
		} `json:"jobs"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil
	}

	var jobs []JobInfo
	for _, j := range resp.Jobs {
		jobs = append(jobs, JobInfo{
			ID:         j.ID,
			Name:       j.Name,
			Status:     j.Status,
			Conclusion: j.Conclusion,
			Steps:      len(j.Steps),
		})
	}
	return jobs
}
