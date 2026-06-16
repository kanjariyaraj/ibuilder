package ai

import (
	"fmt"
	"os"
	"strings"
)

type ClassificationReport struct {
	Category   FailureCategory `json:"category"`
	Matches    int             `json:"matches"`
	Lines      []string        `json:"lines,omitempty"`
	Confidence float64         `json:"confidence"`
}

func (a *Analyzer) ClassifyLogs(logs []string) []ClassificationReport {
	reports := make(map[FailureCategory]*ClassificationReport)

	for _, log := range logs {
		for _, line := range strings.Split(log, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			for _, c := range a.classifiers {
				if match := c.Match(line); match != nil {
					if _, ok := reports[*match]; !ok {
						reports[*match] = &ClassificationReport{
							Category: *match,
						}
					}
					reports[*match].Matches++
					if len(reports[*match].Lines) < 5 {
						reports[*match].Lines = append(reports[*match].Lines, line)
					}
				}
			}
		}
	}

	total := 0
	for _, r := range reports {
		total += r.Matches
	}

	var result []ClassificationReport
	for _, r := range reports {
		if total > 0 {
			r.Confidence = float64(r.Matches) / float64(total)
		}
		result = append(result, *r)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func (a *Analyzer) ClassifyLogFile(path string) ([]ClassificationReport, error) {
	logContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	return a.ClassifyLogs([]string{string(logContent)}), nil
}
