package github

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/kanjariyaraj/Builder/internal/errors"
)

type Workflow struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	State     string `json:"state"`
}

type DispatchRequest struct {
	Ref      string            `json:"ref"`
	Inputs   map[string]string `json:"inputs,omitempty"`
}

type WorkflowRun struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	Conclusion    string    `json:"conclusion"`
	HeadBranch    string    `json:"head_branch"`
	RunNumber     int       `json:"run_number"`
	Event         string    `json:"event"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	HTMLURL       string    `json:"html_url"`
	WorkflowID    int64     `json:"workflow_id"`
}

type workflowRunsResponse struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type Artifact struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size_in_bytes"`
	ArchiveURL  string `json:"archive_download_url"`
}

type artifactsResponse struct {
	Artifacts []Artifact `json:"artifacts"`
}

func ListWorkflows(client *Client, owner, name string) ([]Workflow, error) {
	data, err := client.Get(fmt.Sprintf("/repos/%s/%s/actions/workflows", owner, name))
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "failed to list workflows", err)
	}

	var resp struct {
		Workflows []Workflow `json:"workflows"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to parse workflows", err)
	}
	return resp.Workflows, nil
}

func DispatchWorkflow(client *Client, owner, name, workflowID, ref string, inputs map[string]string) error {
	req := DispatchRequest{
		Ref:    ref,
		Inputs: inputs,
	}

	dispatchURL := fmt.Sprintf("/repos/%s/%s/actions/workflows/%s/dispatches", owner, name, workflowID)

	body, _ := json.Marshal(req)

	url := APIRoot + dispatchURL
	httpReq, err := client.NewRequest("POST", url, string(body))
	if err != nil {
		return errors.Wrap(errors.KindInternal, "failed to create dispatch request", err)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}

	respBody, _ := io.ReadAll(resp.Body)
	return errors.Wrap(errors.KindNetwork, fmt.Sprintf("dispatch failed (status %d)", resp.StatusCode), fmt.Errorf("%s", string(respBody)))
}

func GetWorkflowRun(client *Client, owner, name string, runID int64) (*WorkflowRun, error) {
	data, err := client.Get(fmt.Sprintf("/repos/%s/%s/actions/runs/%d", owner, name, runID))
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "failed to get workflow run", err)
	}

	var run WorkflowRun
	if err := json.Unmarshal(data, &run); err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to parse workflow run", err)
	}
	return &run, nil
}

func ListWorkflowRuns(client *Client, owner, name string, limit int) ([]WorkflowRun, error) {
	url := fmt.Sprintf("/repos/%s/%s/actions/runs?per_page=%d&event=workflow_dispatch", owner, name, limit)
	data, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "failed to list workflow runs", err)
	}

	var resp workflowRunsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to parse workflow runs", err)
	}
	return resp.WorkflowRuns, nil
}

func ListArtifacts(client *Client, owner, name string, runID int64) ([]Artifact, error) {
	url := fmt.Sprintf("/repos/%s/%s/actions/runs/%d/artifacts", owner, name, runID)
	data, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "failed to list artifacts", err)
	}

	var resp artifactsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, errors.Wrap(errors.KindInternal, "failed to parse artifacts", err)
	}
	return resp.Artifacts, nil
}

func DownloadArtifact(client *Client, owner, name string, artifactID int64) ([]byte, error) {
	url := fmt.Sprintf("/repos/%s/%s/actions/artifacts/%d/zip", owner, name, artifactID)
	data, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrap(errors.KindNetwork, "failed to download artifact", err)
	}
	return data, nil
}
