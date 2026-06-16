package release

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type NotesFormat string

const (
	NotesMarkdown NotesFormat = "markdown"
	NotesJSON     NotesFormat = "json"
	NotesHTML     NotesFormat = "html"
)

type ReleaseNotes struct {
	Version       string   `json:"version"`
	BuildNumber   string   `json:"build_number"`
	Date          string   `json:"date"`
	Summary       string   `json:"summary"`
	Features      []string `json:"features"`
	Fixes         []string `json:"fixes"`
	Breaking      []string `json:"breaking_changes"`
	UpgradeNotes  []string `json:"upgrade_notes"`
}

func (s *Session) GenerateNotes() (*ReleaseNotes, error) {
	s.logInfo("generating release notes from git history")

	dir := s.ProjectDir()
	if dir == "" {
		return nil, fmt.Errorf("no project directory set")
	}

	version := "1.0.0"
	buildNum := fmt.Sprintf("%d", time.Now().Unix())

	gitLog, _ := s.getGitLog(dir)
	features, fixes, breaking := categorizeCommits(gitLog)

	notes := &ReleaseNotes{
		Version:      version,
		BuildNumber:  buildNum,
		Date:         time.Now().Format("2006-01-02"),
		Summary:      fmt.Sprintf("Release %s (build %s)", version, buildNum),
		Features:     features,
		Fixes:        fixes,
		Breaking:     breaking,
		UpgradeNotes: []string{},
	}

	s.logInfo("release notes generated", "version", version)
	return notes, nil
}

func (s *Session) SaveNotes(notes *ReleaseNotes, format NotesFormat) (string, error) {
	dir := s.ProjectDir()
	outputDir := filepath.Join(dir, ".build", "releases")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := fmt.Sprintf("release_notes_%s_%s", notes.Version, s.Timestamp())
	var content string

	switch format {
	case NotesJSON:
		filename += ".json"
		content = renderJSON(notes)
	case NotesHTML:
		filename += ".html"
		content = renderHTML(notes)
	default:
		filename += ".md"
		content = renderMarkdown(notes)
	}

	path := filepath.Join(outputDir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write release notes: %w", err)
	}

	s.logInfo("release notes saved", "path", path)
	return path, nil
}

func (s *Session) getGitLog(dir string) ([]string, error) {
	cmd := exec.Command("git", "log", "--oneline", "-30")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func categorizeCommits(commits []string) (features, fixes, breaking []string) {
	for _, c := range commits {
		lower := strings.ToLower(c)
		switch {
		case strings.Contains(lower, "breaking") || strings.Contains(lower, "major"):
			breaking = append(breaking, c)
		case strings.Contains(lower, "feat") || strings.Contains(lower, "feature") || strings.Contains(lower, "add"):
			features = append(features, c)
		case strings.Contains(lower, "fix") || strings.Contains(lower, "bug") || strings.Contains(lower, "patch"):
			fixes = append(fixes, c)
		}
	}
	return
}

func renderMarkdown(notes *ReleaseNotes) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Release %s\n\n", notes.Version))
	b.WriteString(fmt.Sprintf("**Build:** %s  \n", notes.BuildNumber))
	b.WriteString(fmt.Sprintf("**Date:** %s  \n\n", notes.Date))
	b.WriteString(fmt.Sprintf("%s\n\n", notes.Summary))

	if len(notes.Features) > 0 {
		b.WriteString("## Features\n\n")
		for _, f := range notes.Features {
			b.WriteString(fmt.Sprintf("- %s\n", f))
		}
		b.WriteString("\n")
	}

	if len(notes.Fixes) > 0 {
		b.WriteString("## Bug Fixes\n\n")
		for _, f := range notes.Fixes {
			b.WriteString(fmt.Sprintf("- %s\n", f))
		}
		b.WriteString("\n")
	}

	if len(notes.Breaking) > 0 {
		b.WriteString("## Breaking Changes\n\n")
		for _, br := range notes.Breaking {
			b.WriteString(fmt.Sprintf("- %s\n", br))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func renderJSON(notes *ReleaseNotes) string {
	return fmt.Sprintf(`{
  "version": %q,
  "build_number": %q,
  "date": %q,
  "summary": %q,
  "features": [%s],
  "fixes": [%s],
  "breaking_changes": [%s]
}`,
		notes.Version, notes.BuildNumber, notes.Date, notes.Summary,
		joinQuoted(notes.Features), joinQuoted(notes.Fixes), joinQuoted(notes.Breaking))
}

func renderHTML(notes *ReleaseNotes) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><head><title>")
	b.WriteString(fmt.Sprintf("Release %s", notes.Version))
	b.WriteString("</title></head><body>")
	b.WriteString(fmt.Sprintf("<h1>Release %s</h1>", notes.Version))
	b.WriteString(fmt.Sprintf("<p>Build: %s | Date: %s</p>", notes.BuildNumber, notes.Date))
	b.WriteString("</body></html>")
	return b.String()
}

func joinQuoted(items []string) string {
	if len(items) == 0 {
		return ""
	}
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("%q", item)
	}
	return strings.Join(quoted, ", ")
}
