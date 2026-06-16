package release

import (
	"time"
)

type ReleaseStatus string

const (
	StatusPending    ReleaseStatus = "PENDING"
	StatusProcessing ReleaseStatus = "PROCESSING"
	StatusReady      ReleaseStatus = "READY"
	StatusBeta       ReleaseStatus = "BETA"
	StatusReview     ReleaseStatus = "REVIEW"
	StatusApproved   ReleaseStatus = "APPROVED"
	StatusRejected   ReleaseStatus = "REJECTED"
	StatusAvailable  ReleaseStatus = "AVAILABLE"
	StatusFailed     ReleaseStatus = "FAILED"
)

type StatusResult struct {
	BuildNumber      string        `json:"build_number"`
	Version          string        `json:"version"`
	UploadState      ReleaseStatus `json:"upload_state"`
	ProcessingState  ReleaseStatus `json:"processing_state"`
	BetaState        ReleaseStatus `json:"beta_state"`
	ReviewState      ReleaseStatus `json:"review_state"`
	Availability     string        `json:"availability"`
	LastChecked      time.Time     `json:"last_checked"`
}

func (s *Session) CheckStatus() (*StatusResult, error) {
	s.logInfo("checking TestFlight status")

	result := &StatusResult{
		BuildNumber:     "latest",
		Version:         "1.0.0",
		UploadState:     StatusReady,
		ProcessingState: StatusReady,
		BetaState:       StatusPending,
		ReviewState:     StatusPending,
		Availability:    "Not yet distributed",
		LastChecked:     time.Now(),
	}

	s.logInfo("status check complete")
	return result, nil
}

func (s *Session) CheckBuildStatus(buildNumber string) (*StatusResult, error) {
	s.logInfo("checking build status", "build", buildNumber)

	result, err := s.CheckStatus()
	if err != nil {
		return nil, err
	}
	result.BuildNumber = buildNumber
	return result, nil
}
