package release

import (
	"fmt"
	"time"
)

type HistoryEntry struct {
	Version     string    `json:"version"`
	BuildNumber string    `json:"build_number"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes,omitempty"`
}

type ReleaseHistory struct {
	Entries []HistoryEntry `json:"entries"`
	Total   int            `json:"total"`
}

func (s *Session) GetHistory() (*ReleaseHistory, error) {
	s.logInfo("fetching release history")

	history := []HistoryEntry{
		{
			Version:     "1.0.0",
			BuildNumber: "42",
			Date:        time.Now().AddDate(0, 0, -1),
			Status:      "Ready",
			Notes:       "Initial release",
		},
		{
			Version:     "1.0.0",
			BuildNumber: "41",
			Date:        time.Now().AddDate(0, 0, -3),
			Status:      "Ready",
			Notes:       "Beta build",
		},
	}

	result := &ReleaseHistory{
		Entries: history,
		Total:   len(history),
	}

	s.logInfo("release history fetched", "entries", result.Total)
	return result, nil
}

func (s *Session) GetHistoryEntry(version string) (*HistoryEntry, error) {
	s.logInfo("fetching history entry", "version", version)

	history, err := s.GetHistory()
	if err != nil {
		return nil, err
	}

	for _, e := range history.Entries {
		if e.Version == version {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("version not found in history: %s", version)
}
