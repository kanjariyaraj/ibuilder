package releasepipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type NotesStageResult struct {
	NotesPath string `json:"notes_path"`
	Summary   string `json:"summary"`
	Version   string `json:"version"`
}

func (p *Pipeline) GenerateNotes() *StageResult {
	p.logInfo("stage 4: generating release notes")

	dir := p.ProjectDir()
	version := time.Now().Format("2006.01.02")
	buildNum := fmt.Sprintf("%d", time.Now().Unix())

	if p.IsDryRun() {
		return p.addResult(StageNotes, true,
			fmt.Sprintf("[DRY RUN] Would generate release notes for v%s (build %s)", version, buildNum), nil)
	}

	outputDir := filepath.Join(dir, ".build", "releases")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return p.addResult(StageNotes, false, "cannot create notes directory", err)
	}

	content := fmt.Sprintf(`# Release v%s

**Build:** %s
**Date:** %s
**Status:** %s

## Summary

Automated release via iBuilder release pipeline.

## Changes

See git log for detailed changelog.
`, version, buildNum, time.Now().Format("2006-01-02"), p.Mode())

	path := filepath.Join(outputDir, fmt.Sprintf("release_notes_%s.md", p.Timestamp()))
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return p.addResult(StageNotes, false, "cannot write release notes", err)
	}

	result := &NotesStageResult{
		NotesPath: path,
		Summary:   fmt.Sprintf("Release v%s", version),
		Version:   version,
	}

	msg := fmt.Sprintf("notes generated: %s (v%s)", result.NotesPath, result.Version)
	return p.addResult(StageNotes, true, msg, nil)
}
