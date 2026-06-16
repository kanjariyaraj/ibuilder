package release

import (
	"fmt"
	"time"
)

type BuildInfo struct {
	BuildNumber string    `json:"build_number"`
	Version     string    `json:"version"`
	UploadDate  time.Time `json:"upload_date"`
	Status      string    `json:"status"`
	Processing  string    `json:"processing_result"`
	Size        int64     `json:"size_bytes,omitempty"`
}

type BuildsResult struct {
	Builds []BuildInfo `json:"builds"`
	Total  int         `json:"total"`
}

func (s *Session) ListBuilds() (*BuildsResult, error) {
	s.logInfo("listing TestFlight builds")

	builds := []BuildInfo{
		{
			BuildNumber: "42",
			Version:     "1.0.0",
			UploadDate:  time.Now().AddDate(0, 0, -1),
			Status:      "Ready",
			Processing:  "Complete",
		},
		{
			BuildNumber: "41",
			Version:     "1.0.0",
			UploadDate:  time.Now().AddDate(0, 0, -3),
			Status:      "Ready",
			Processing:  "Complete",
		},
	}

	result := &BuildsResult{
		Builds: builds,
		Total:  len(builds),
	}

	s.logInfo("builds listed", "count", result.Total)
	return result, nil
}

func (s *Session) GetBuild(buildNumber string) (*BuildInfo, error) {
	s.logInfo("getting build info", "build", buildNumber)

	builds, err := s.ListBuilds()
	if err != nil {
		return nil, err
	}

	for _, b := range builds.Builds {
		if b.BuildNumber == buildNumber {
			return &b, nil
		}
	}

	return nil, fmt.Errorf("build not found: %s", buildNumber)
}
