package artifacts

import (
	"encoding/json"
	"fmt"
)

type BuildRecord struct {
	RunID       int64  `json:"run_id"`
	RunNumber   int    `json:"run_number"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	Branch      string `json:"branch"`
	CommitSHA   string `json:"commit_sha"`
	WorkflowID  int64  `json:"workflow_id"`
	Workflow    string `json:"workflow"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Duration    string `json:"duration"`
	URL         string `json:"url"`
}

type HistoryOptions struct {
	Branch   string
	Workflow string
	Status   string
	Limit    int
	Page     int
	JSON     bool
}

func (m *ArtifactManager) GetHistory(opts *HistoryOptions) ([]BuildRecord, error) {
	if opts.Limit == 0 {
		opts.Limit = 30
	}

	path := fmt.Sprintf("/repos/%s/%s/actions/runs?per_page=%d&page=%d", m.owner, m.repo, opts.Limit, opts.Page)
	if opts.Branch != "" {
		path += "&branch=" + opts.Branch
	}
	if opts.Status != "" {
		path += "&status=" + opts.Status
	}

	data, err := m.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch build history: %w", err)
	}

	var resp struct {
		WorkflowRuns []struct {
			ID          int64  `json:"id"`
			RunNumber   int    `json:"run_number"`
			Status      string `json:"status"`
			Conclusion  string `json:"conclusion"`
			HeadBranch  string `json:"head_branch"`
			HeadSHA     string `json:"head_sha"`
			WorkflowID  int64  `json:"workflow_id"`
			Name        string `json:"name"`
			CreatedAt   string `json:"created_at"`
			UpdatedAt   string `json:"updated_at"`
			RunDuration int    `json:"run_duration_ms"`
			HTMLURL     string `json:"html_url"`
		} `json:"workflow_runs"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse build history: %w", err)
	}

	var records []BuildRecord
	for _, r := range resp.WorkflowRuns {
		records = append(records, BuildRecord{
			RunID:      r.ID,
			RunNumber:  r.RunNumber,
			Status:     r.Status,
			Conclusion: r.Conclusion,
			Branch:     r.HeadBranch,
			CommitSHA:  r.HeadSHA,
			WorkflowID: r.WorkflowID,
			Workflow:   r.Name,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
			Duration:   fmt.Sprintf("%ds", r.RunDuration/1000),
			URL:        r.HTMLURL,
		})
	}

	return records, nil
}
