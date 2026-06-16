package releasepipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type GitHubReleaseResult struct {
	Tag     string `json:"tag"`
	URL     string `json:"url"`
	Created bool   `json:"created"`
}

func (p *Pipeline) CreateGitHubRelease() *StageResult {
	p.logInfo("stage 6: creating GitHub release")

	dir := p.ProjectDir()
	tag := fmt.Sprintf("v%s", time.Now().Format("2006.01.02"))

	if p.IsDryRun() {
		return p.addResult(StageRelease, true,
			fmt.Sprintf("[DRY RUN] Would create GitHub release: %s", tag), nil)
	}

	ipaPath, err := p.findIPA()
	if err != nil {
		ipaPath = ""
	}

	notesPath := filepath.Join(dir, ".build", "releases")
	notesFile := ""
	entries, _ := os.ReadDir(notesPath)
	for _, e := range entries {
		if !e.IsDir() && len(e.Name()) > 0 {
			notesFile = filepath.Join(notesPath, e.Name())
		}
	}

	args := []string{"release", "create", tag, "--title", tag, "--generate-notes"}
	if notesFile != "" {
		args = append(args, "--notes-file", notesFile)
	}
	if ipaPath != "" {
		args = append(args, ipaPath)
	}

	cmd := exec.Command("gh", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return p.addResult(StageRelease, false, "GitHub release creation failed",
			fmt.Errorf("%s: %s", err.Error(), string(output)))
	}

	result := &GitHubReleaseResult{
		Tag:     tag,
		Created: true,
	}

	if len(output) > 0 {
		result.URL = string(output)
	}

	msg := fmt.Sprintf("GitHub release created: %s", tag)
	return p.addResult(StageRelease, true, msg, nil)
}
