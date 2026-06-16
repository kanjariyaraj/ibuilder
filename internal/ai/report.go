package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ReportFormat string

const (
	FormatMarkdown ReportFormat = "markdown"
	FormatJSON     ReportFormat = "json"
	FormatHTML     ReportFormat = "html"
)

type Report struct {
	Title       string           `json:"title"`
	GeneratedAt time.Time        `json:"generated_at"`
	Healthy     bool             `json:"healthy"`
	Issues      []AnalysisResult `json:"issues"`
	Checks      []DoctorCheck    `json:"checks"`
	Summary     string           `json:"summary"`
}

func (s *Session) GenerateReport(dir string, format ReportFormat) (string, error) {
	s.log.Info("generating AI report", "format", format)

	doctorReport := s.Doctor()
	logs, _ := s.collectLogs(dir)

	result := s.analyzer.Analyze(logs)

	var issues []AnalysisResult
	if result.Category != CatUnknown {
		issues = append(issues, *result)
	}

	report := &Report{
		Title:       "iBuilder AI Diagnostic Report",
		GeneratedAt: time.Now(),
		Healthy:     doctorReport.Healthy && len(issues) == 0,
		Issues:      issues,
		Checks:      doctorReport.Checks,
	}

	if report.Healthy {
		report.Summary = "No critical issues detected"
	} else {
		report.Summary = fmt.Sprintf("Found %d issue(s) and %d warning(s)", len(issues), warningCount(doctorReport.Checks))
	}

	outputDir := filepath.Join(dir, ".build", "reports", "ai")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create report directory: %w", err)
	}

	var output string
	var ext string
	switch format {
	case FormatJSON:
		output = s.renderJSON(report)
		ext = "json"
	case FormatHTML:
		output = s.renderHTML(report)
		ext = "html"
	default:
		output = s.renderMarkdown(report)
		ext = "md"
	}

	filename := fmt.Sprintf("ai_diagnostic_%s.%s", s.Timestamp(), ext)
	path := filepath.Join(outputDir, filename)

	if err := os.WriteFile(path, []byte(output), 0644); err != nil {
		return "", fmt.Errorf("failed to write report: %w", err)
	}

	s.log.Info("report generated", "path", path)
	return path, nil
}

func (s *Session) renderMarkdown(report *Report) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s\n\n", report.Title))
	b.WriteString(fmt.Sprintf("**Generated:** %s\n\n", report.GeneratedAt.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("**Overall Status:** %s\n\n", statusBadge(report.Healthy)))

	if report.Summary != "" {
		b.WriteString(fmt.Sprintf("**Summary:** %s\n\n", report.Summary))
	}

	if len(report.Issues) > 0 {
		b.WriteString("## Issues\n\n")
		for _, issue := range report.Issues {
			b.WriteString(fmt.Sprintf("### %s\n\n", issue.Problem))
			b.WriteString(fmt.Sprintf("- **Category:** %s\n", issue.Category))
			b.WriteString(fmt.Sprintf("- **Cause:** %s\n", issue.Cause))
			b.WriteString(fmt.Sprintf("- **Impact:** %s\n", issue.Impact))
			b.WriteString(fmt.Sprintf("- **Fix:** %s\n", issue.Fix))
			b.WriteString(fmt.Sprintf("- **Confidence:** %.0f%%\n\n", issue.Confidence*100))
		}
	}

	if len(report.Checks) > 0 {
		b.WriteString("## System Checks\n\n")
		b.WriteString("| Check | Status | Message |\n")
		b.WriteString("|-------|--------|--------|\n")
		for _, check := range report.Checks {
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", check.Name, check.Status, check.Message))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (s *Session) renderJSON(report *Report) string {
	return fmt.Sprintf(`{
  "title": %q,
  "generated_at": %q,
  "healthy": %v,
  "summary": %q,
  "issues_count": %d,
  "checks_count": %d
}`,
		report.Title,
		report.GeneratedAt.Format(time.RFC3339),
		report.Healthy,
		report.Summary,
		len(report.Issues),
		len(report.Checks),
	)
}

func (s *Session) renderHTML(report *Report) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><head><title>")
	b.WriteString(report.Title)
	b.WriteString("</title></head><body>")
	b.WriteString(fmt.Sprintf("<h1>%s</h1>", report.Title))
	b.WriteString(fmt.Sprintf("<p>Generated: %s</p>", report.GeneratedAt.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("<p>Status: <strong>%s</strong></p>", statusBadge(report.Healthy)))
	b.WriteString("</body></html>")
	return b.String()
}

func statusBadge(healthy bool) string {
	if healthy {
		return "✅ HEALTHY"
	}
	return "❌ ISSUES FOUND"
}

func warningCount(checks []DoctorCheck) int {
	count := 0
	for _, c := range checks {
		if c.Status == "WARNING" {
			count++
		}
	}
	return count
}
