package artifacts

import (
	"fmt"
	"time"
)

type CleanupOptions struct {
	OlderThan time.Duration
	Keep      int
	All       bool
}

func (m *ArtifactManager) Cleanup(opts *CleanupOptions) (*CleanupResult, error) {
	result := &CleanupResult{}

	if opts.All {
		if err := m.storage.CleanAll(); err != nil {
			return nil, fmt.Errorf("cleanup failed: %w", err)
		}
		result.FilesRemoved = -1
		return result, nil
	}

	if opts.OlderThan > 0 {
		count, err := m.storage.CleanOldArtifacts(opts.OlderThan)
		if err != nil {
			return nil, fmt.Errorf("failed to clean old artifacts: %w", err)
		}
		result.FilesRemoved = count
	}

	if opts.Keep > 0 {
		result.FilesRemoved = 0
	}

	return result, nil
}

type CleanupResult struct {
	FilesRemoved int
}
