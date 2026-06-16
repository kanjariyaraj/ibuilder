package release

import (
	"os"
	"path/filepath"
	"strings"
)

type PrepareResult struct {
	Success      bool     `json:"success"`
	SigningOK    bool     `json:"signing_ok"`
	BuildOK      bool     `json:"build_ok"`
	IPAOK        bool     `json:"ipa_ok"`
	MetadataOK   bool     `json:"metadata_ok"`
	NotesOK      bool     `json:"notes_ok"`
	GitStateOK   bool     `json:"git_state_ok"`
	Issues       []string `json:"issues,omitempty"`
	Warnings     []string `json:"warnings,omitempty"`
}

func (s *Session) Prepare() (*PrepareResult, error) {
	s.logInfo("preparing release")

	result := &PrepareResult{}

	result.SigningOK = s.checkSigning()
	result.BuildOK = s.checkBuild()
	result.IPAOK = s.checkIPA()
	result.MetadataOK = s.checkMetadata()
	result.NotesOK = s.checkReleaseNotes()
	result.GitStateOK = s.checkGitState()

	if !result.SigningOK {
		result.Issues = append(result.Issues, "Signing configuration is incomplete")
	}
	if !result.BuildOK {
		result.Issues = append(result.Issues, "Build artifacts not found")
	}
	if !result.IPAOK {
		result.Issues = append(result.Issues, "IPA file not found or invalid")
	}
	if !result.MetadataOK {
		result.Warnings = append(result.Warnings, "App metadata may be incomplete")
	}
	if !result.NotesOK {
		result.Warnings = append(result.Warnings, "Release notes not generated")
	}
	if !result.GitStateOK {
		result.Warnings = append(result.Warnings, "Git working directory is not clean")
	}

	result.Success = result.SigningOK && result.BuildOK && result.IPAOK

	s.logInfo("release preparation complete",
		"success", result.Success,
		"issues", len(result.Issues),
	)
	return result, nil
}

func (s *Session) checkSigning() bool {
	dir := s.ProjectDir()
	cfgPath := filepath.Join(dir, "builder.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return false
	}
	content := string(data)
	return strings.Contains(content, `"team_id"`) &&
		!strings.Contains(content, `"team_id": ""`)
}

func (s *Session) checkBuild() bool {
	dir := s.ProjectDir()
	buildDir := filepath.Join(dir, ".build")
	if _, err := os.Stat(buildDir); err != nil {
		return false
	}
	entries, _ := os.ReadDir(buildDir)
	return len(entries) > 0
}

func (s *Session) checkIPA() bool {
	_, err := s.findIPA()
	return err == nil
}

func (s *Session) checkMetadata() bool {
	dir := s.ProjectDir()
	cfgPath := filepath.Join(dir, "builder.json")
	_, err := os.Stat(cfgPath)
	return err == nil
}

func (s *Session) checkReleaseNotes() bool {
	dir := s.ProjectDir()
	notesDir := filepath.Join(dir, ".build", "releases")
	entries, err := os.ReadDir(notesDir)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

func (s *Session) checkGitState() bool {
	dir := s.ProjectDir()
	gitDir := filepath.Join(dir, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}
