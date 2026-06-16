package releasepipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ReportData struct {
	Version    string         `json:"version"`
	Mode       ReleaseMode    `json:"mode"`
	DryRun     bool           `json:"dry_run"`
	Duration   string         `json:"duration"`
	Stages     []StageResult  `json:"stages"`
	Success    bool           `json:"success"`
}

func (p *Pipeline) GenerateReport(started time.Time) *StageResult {
	p.logInfo("stage 7: generating release report")

	dir := p.ProjectDir()
	reportDir := filepath.Join(dir, ".build", "reports", "release")
	os.MkdirAll(reportDir, 0755)

	completed := time.Now()
	duration := completed.Sub(started).Round(time.Second).String()

	success := true
	for _, r := range p.results {
		if !r.Success {
			success = false
			break
		}
	}

	data := &ReportData{
		Version:  time.Now().Format("2006.01.02"),
		Mode:     p.mode,
		DryRun:   p.dryRun,
		Duration: duration,
		Stages:   p.results,
		Success:  success,
	}

	markdown := p.renderMarkdown(data)
	jsonContent := p.renderJSON(data)

	mdPath := filepath.Join(reportDir, fmt.Sprintf("release_report_%s.md", p.Timestamp()))
	jsonPath := filepath.Join(reportDir, fmt.Sprintf("release_report_%s.json", p.Timestamp()))

	os.WriteFile(mdPath, []byte(markdown), 0644)
	os.WriteFile(jsonPath, []byte(jsonContent), 0644)

	msg := fmt.Sprintf("reports generated: %s, %s", mdPath, jsonPath)
	return p.addResult(StageReport, true, msg, nil)
}

func (p *Pipeline) renderMarkdown(data *ReportData) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Release Report\n\n"))
	b.WriteString(fmt.Sprintf("**Version:** %s  \n", data.Version))
	b.WriteString(fmt.Sprintf("**Mode:** %s  \n", data.Mode))
	b.WriteString(fmt.Sprintf("**Duration:** %s  \n", data.Duration))
	b.WriteString(fmt.Sprintf("**Status:** %s  \n\n", statusText(data.Success)))

	b.WriteString("## Pipeline Stages\n\n")
	b.WriteString("| Stage | Status | Message |\n")
	b.WriteString("|-------|--------|--------|\n")
	for _, r := range data.Stages {
		icon := "✓"
		if !r.Success {
			icon = "✗"
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", r.Stage, icon, r.Message))
	}

	return b.String()
}

func (p *Pipeline) renderJSON(data *ReportData) string {
	stages := make([]string, len(data.Stages))
	for i, r := range data.Stages {
		errField := ""
		if r.Error != "" {
			errField = fmt.Sprintf(", \"error\": %q", r.Error)
		}
		stages[i] = fmt.Sprintf("{\"stage\": %q, \"success\": %v, \"message\": %q%s}",
			r.Stage, r.Success, r.Message, errField)
	}

	return fmt.Sprintf(`{
  "version": %q,
  "mode": %q,
  "dry_run": %v,
  "duration": %q,
  "success": %v,
  "stages": [%s]
}`, data.Version, data.Mode, data.DryRun, data.Duration, data.Success,
		strings.Join(stages, ", "))
}

func statusText(success bool) string {
	if success {
		return "✓ SUCCESS"
	}
	return "✗ FAILED"
}
