package releasepipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type BuildStageResult struct {
	WorkflowID string `json:"workflow_id"`
	BuildNum   string `json:"build_number"`
	Artifact   string `json:"artifact"`
}

func (p *Pipeline) Build() *StageResult {
	p.logInfo("stage 2: triggering build")

	dir := p.ProjectDir()
	cfgPath := filepath.Join(dir, "builder.json")

	if p.IsDryRun() {
		msg := "[DRY RUN] Would trigger workflow build via: gh workflow run"
		return p.addResult(StageBuild, true, msg, nil)
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return p.addResult(StageBuild, false, "builder.json not found",
			fmt.Errorf("config file missing"))
	}

	cmd := exec.Command("gh", "workflow", "run", ".github/workflows/build.yml",
		"--ref", "main",
		"-f", "mode=release")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return p.addResult(StageBuild, false, "workflow trigger failed",
			fmt.Errorf("%s: %s", err.Error(), strings.TrimSpace(string(output))))
	}

	p.logInfo("build workflow triggered")

	time.Sleep(5 * time.Second)

	ipaPath, err := p.findIPA()
	if err != nil {
		return p.addResult(StageBuild, false, "no artifact found after build", err)
	}

	msg := fmt.Sprintf("build complete, artifact: %s", ipaPath)
	return p.addResult(StageBuild, true, msg, nil)
}
