package ai

import (
	"fmt"
	"time"
)

type ExplainResult struct {
	Problem       string          `json:"problem"`
	Category      FailureCategory `json:"category"`
	Cause         string          `json:"cause"`
	Impact        string          `json:"impact"`
	Fix            string          `json:"fix"`
	Confidence    float64         `json:"confidence"`
	ConfidenceLabel string        `json:"confidence_label"`
	Timestamp     time.Time       `json:"timestamp"`
}

func (s *Session) Explain(dir string) (*ExplainResult, error) {
	s.log.Info("running AI explain", "dir", dir)

	logs, err := s.collectLogs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to collect logs: %w", err)
	}

	if len(logs) == 0 {
		return &ExplainResult{
			Problem:         "No data to analyze",
			Category:        CatUnknown,
			Cause:           "No build logs, workflow logs, or error reports found",
			Impact:          "Nothing to explain",
			Fix:             "Run a build or workflow first, then run 'builder ai explain'",
			Confidence:      0.0,
			ConfidenceLabel: "NONE",
			Timestamp:       time.Now(),
		}, nil
	}

	analysis := s.analyzer.Analyze(logs)

	confidenceLabel := "LOW"
	if analysis.Confidence >= 0.8 {
		confidenceLabel = "HIGH"
	} else if analysis.Confidence >= 0.5 {
		confidenceLabel = "MEDIUM"
	}

	result := &ExplainResult{
		Problem:         analysis.Problem,
		Category:        analysis.Category,
		Cause:           analysis.Cause,
		Impact:          analysis.Impact,
		Fix:             analysis.Fix,
		Confidence:      analysis.Confidence,
		ConfidenceLabel: confidenceLabel,
		Timestamp:       time.Now(),
	}

	s.log.Info("explain analysis complete",
		"category", result.Category,
		"confidence", result.ConfidenceLabel,
	)

	return result, nil
}
