package artifacts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kanjariyaraj/Builder/internal/errors"
)

type LogsOptions struct {
	RunID    int64
	Latest   bool
	SavePath string
}

func (m *ArtifactManager) FetchLogs(opts *LogsOptions) (string, error) {
	var runID int64

	if opts.Latest {
		history, err := m.GetHistory(&HistoryOptions{Limit: 1})
		if err != nil {
			return "", fmt.Errorf("failed to get latest build: %w", err)
		}
		if len(history) == 0 {
			return "", errors.New(errors.KindNotFound, "no builds found")
		}
		runID = history[0].RunID
	} else {
		runID = opts.RunID
	}

	if runID == 0 {
		return "", fmt.Errorf("no build specified, use --run-id or --latest")
	}

	logPath := opts.SavePath
	if logPath == "" {
		logPath = filepath.Join(".build", "logs", fmt.Sprintf("run-%d.log", runID))
	}

	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return "", errors.Wrap(errors.KindInternal, "failed to create logs directory", err)
	}

	data, err := m.client.Get(fmt.Sprintf("/repos/%s/%s/actions/runs/%d/logs", m.owner, m.repo, runID))
	if err != nil {
		return "", fmt.Errorf("failed to fetch logs: %w", err)
	}

	if err := os.WriteFile(logPath, data, 0644); err != nil {
		return "", errors.Wrap(errors.KindInternal, "failed to save logs", err)
	}

	return logPath, nil
}

func (m *ArtifactManager) OpenBuildURL(runID int64) (string, error) {
	inspect, err := m.InspectBuild(runID)
	if err != nil {
		return "", err
	}
	return inspect.URL, nil
}
