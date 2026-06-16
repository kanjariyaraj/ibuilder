package releasepipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ValidationResult struct {
	EnvironmentOK bool     `json:"environment_ok"`
	ConfigOK      bool     `json:"config_ok"`
	SigningOK     bool     `json:"signing_ok"`
	RepoOK        bool     `json:"repo_ok"`
	GitHubOK      bool     `json:"github_ok"`
	ProjectOK     bool     `json:"project_ok"`
	Issues        []string `json:"issues"`
}

func (p *Pipeline) Validate() *StageResult {
	p.logInfo("stage 1: validating environment")

	vr := &ValidationResult{}

	vr.EnvironmentOK = p.checkEnvironment()
	if !vr.EnvironmentOK {
		vr.Issues = append(vr.Issues, "Missing required tools (node, go, git)")
	}

	vr.ConfigOK = p.checkConfig()
	if !vr.ConfigOK {
		vr.Issues = append(vr.Issues, "builder.json not found or invalid")
	}

	vr.SigningOK = p.checkSigning()
	if !vr.SigningOK {
		vr.Issues = append(vr.Issues, "Signing configuration incomplete")
	}

	vr.RepoOK = p.checkRepo()
	if !vr.RepoOK {
		vr.Issues = append(vr.Issues, "Git repository not detected")
	}

	vr.GitHubOK = p.checkGitHub()
	if !vr.GitHubOK {
		vr.Issues = append(vr.Issues, "GitHub authentication not configured")
	}

	vr.ProjectOK = p.checkProject()
	if !vr.ProjectOK {
		vr.Issues = append(vr.Issues, "Project structure invalid")
	}

	success := len(vr.Issues) == 0
	msg := "all validations passed"
	if !success {
		msg = fmt.Sprintf("found %d issue(s): %s", len(vr.Issues), strings.Join(vr.Issues, "; "))
	}

	return p.addResult(StageValidate, success, msg, nil)
}

func (p *Pipeline) checkEnvironment() bool {
	for _, tool := range []string{"node", "go", "git"} {
		cmd := exec.Command(tool, "--version")
		if err := cmd.Run(); err != nil {
			return false
		}
	}
	return true
}

func (p *Pipeline) checkConfig() bool {
	cfgPath := filepath.Join(p.ProjectDir(), "builder.json")
	_, err := os.Stat(cfgPath)
	return err == nil
}

func (p *Pipeline) checkSigning() bool {
	cfgPath := filepath.Join(p.ProjectDir(), "builder.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return false
	}
	content := string(data)
	return strings.Contains(content, `"team_id"`) &&
		!strings.Contains(content, `"team_id": ""`)
}

func (p *Pipeline) checkRepo() bool {
	gitDir := filepath.Join(p.ProjectDir(), ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

func (p *Pipeline) checkGitHub() bool {
	cmd := exec.Command("gh", "auth", "status")
	return cmd.Run() == nil
}

func (p *Pipeline) checkProject() bool {
	if p.ProjectDir() == "" {
		return false
	}
	entries, err := os.ReadDir(p.ProjectDir())
	if err != nil {
		return false
	}
	return len(entries) > 0
}
