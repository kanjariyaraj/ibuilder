package artifacts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kanjariyaraj/Builder/internal/github"
)

type Artifact struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	BuildNumber  int       `json:"build_number"`
	WorkflowRun  string    `json:"workflow_run"`
	WorkflowName string    `json:"workflow_name"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	DownloadURL  string    `json:"download_url"`
	Status       string    `json:"status"`
}

type ArtifactFilter struct {
	Name   string
	Build  int
	Limit  int
	All    bool
	Latest bool
}

type ArtifactManager struct {
	client  *github.Client
	owner   string
	repo    string
	storage *Storage
}

func NewArtifactManager(client *github.Client, owner, repo string) *ArtifactManager {
	return &ArtifactManager{
		client:  client,
		owner:   owner,
		repo:    repo,
		storage: NewStorage(),
	}
}

func (m *ArtifactManager) ListWorkflowArtifacts(runID int64) ([]Artifact, error) {
	data, err := m.client.Get(fmt.Sprintf("/repos/%s/%s/actions/runs/%d/artifacts", m.owner, m.repo, runID))
	if err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}

	var resp struct {
		Artifacts []struct {
			ID           int64  `json:"id"`
			Name        string `json:"name"`
			Size        int64  `json:"size_in_bytes"`
			CreatedAt   string `json:"created_at"`
			ExpiresAt   string `json:"expires_at"`
			DownloadURL string `json:"archive_download_url"`
			WorkflowRun struct {
				ID     int64  `json:"id"`
				RunNumber int `json:"run_number"`
			} `json:"workflow_run"`
		} `json:"artifacts"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse artifacts: %w", err)
	}

	var artifacts []Artifact
	for _, a := range resp.Artifacts {
		createdAt, _ := time.Parse(time.RFC3339, a.CreatedAt)
		expiresAt, _ := time.Parse(time.RFC3339, a.ExpiresAt)
		artifacts = append(artifacts, Artifact{
			ID:          a.ID,
			Name:        a.Name,
			Size:        a.Size,
			BuildNumber: a.WorkflowRun.RunNumber,
			CreatedAt:   createdAt,
			ExpiresAt:   expiresAt,
			DownloadURL: a.DownloadURL,
			Status:      "available",
		})
	}

	return artifacts, nil
}

func (m *ArtifactManager) ListAllArtifacts() ([]Artifact, error) {
	data, err := m.client.Get(fmt.Sprintf("/repos/%s/%s/actions/artifacts?per_page=100", m.owner, m.repo))
	if err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}

	var resp struct {
		Artifacts []struct {
			ID           int64  `json:"id"`
			Name        string `json:"name"`
			Size        int64  `json:"size_in_bytes"`
			CreatedAt   string `json:"created_at"`
			ExpiresAt   string `json:"expires_at"`
			DownloadURL string `json:"archive_download_url"`
			WorkflowRun struct {
				ID         int64 `json:"id"`
				RunNumber  int   `json:"run_number"`
			} `json:"workflow_run"`
		} `json:"artifacts"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse artifacts: %w", err)
	}

	var artifacts []Artifact
	for _, a := range resp.Artifacts {
		createdAt, _ := time.Parse(time.RFC3339, a.CreatedAt)
		expiresAt, _ := time.Parse(time.RFC3339, a.ExpiresAt)
		artifacts = append(artifacts, Artifact{
			ID:          a.ID,
			Name:        a.Name,
			Size:        a.Size,
			BuildNumber: a.WorkflowRun.RunNumber,
			CreatedAt:   createdAt,
			ExpiresAt:   expiresAt,
			DownloadURL: a.DownloadURL,
			Status:      "available",
		})
	}

	return artifacts, nil
}
