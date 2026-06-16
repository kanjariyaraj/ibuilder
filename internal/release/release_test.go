package release

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	s := NewSession(nil)
	if s == nil {
		t.Fatal("expected session, got nil")
	}
}

func TestUploadValidate(t *testing.T) {
	s := NewSession(nil)
	err := s.validateIPA("nonexistent.ipa")
	if err == nil {
		t.Error("expected error for nonexistent IPA")
	}
}

func TestUploadLatestNoProject(t *testing.T) {
	s := NewSession(nil)
	_, err := s.UploadLatest()
	if err == nil {
		t.Error("expected error with no project dir")
	}
}

func TestCheckStatus(t *testing.T) {
	s := NewSession(nil)
	result, err := s.CheckStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UploadState != StatusReady {
		t.Errorf("expected StatusReady, got %s", result.UploadState)
	}
}

func TestCheckBuildStatus(t *testing.T) {
	s := NewSession(nil)
	result, err := s.CheckBuildStatus("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BuildNumber != "42" {
		t.Errorf("expected build 42, got %s", result.BuildNumber)
	}
}

func TestListGroups(t *testing.T) {
	s := NewSession(nil)
	groups, err := s.ListGroups()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if groups.Total == 0 {
		t.Error("expected at least one group")
	}
}

func TestInspectGroup(t *testing.T) {
	s := NewSession(nil)
	group, err := s.InspectGroup("Internal Testers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "Internal Testers" {
		t.Errorf("expected Internal Testers, got %s", group.Name)
	}
}

func TestInspectGroupNotFound(t *testing.T) {
	s := NewSession(nil)
	_, err := s.InspectGroup("Nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent group")
	}
}

func TestListBuilds(t *testing.T) {
	s := NewSession(nil)
	builds, err := s.ListBuilds()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if builds.Total == 0 {
		t.Error("expected at least one build")
	}
}

func TestGetBuild(t *testing.T) {
	s := NewSession(nil)
	build, err := s.GetBuild("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if build.BuildNumber != "42" {
		t.Errorf("expected build 42, got %s", build.BuildNumber)
	}
}

func TestGetBuildNotFound(t *testing.T) {
	s := NewSession(nil)
	_, err := s.GetBuild("999")
	if err == nil {
		t.Error("expected error for nonexistent build")
	}
}

func TestListTesters(t *testing.T) {
	s := NewSession(nil)
	testers, err := s.ListTesters()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if testers.Total == 0 {
		t.Error("expected at least one tester")
	}
}

func TestAddTester(t *testing.T) {
	s := NewSession(nil)
	err := s.AddTester("test@example.com", "Internal Testers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveTester(t *testing.T) {
	s := NewSession(nil)
	err := s.RemoveTester("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateNotes(t *testing.T) {
	s := NewSession(nil)
	s.SetProjectDir(t.TempDir())
	notes, err := s.GenerateNotes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if notes.Version == "" {
		t.Error("expected non-empty version")
	}
}

func TestGetHistory(t *testing.T) {
	s := NewSession(nil)
	history, err := s.GetHistory()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if history.Total == 0 {
		t.Error("expected at least one history entry")
	}
}

func TestGetHistoryEntry(t *testing.T) {
	s := NewSession(nil)
	entry, err := s.GetHistoryEntry("1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", entry.Version)
	}
}

func TestGetHistoryEntryNotFound(t *testing.T) {
	s := NewSession(nil)
	_, err := s.GetHistoryEntry("99.99.99")
	if err == nil {
		t.Error("expected error for nonexistent version")
	}
}

func TestPrepare(t *testing.T) {
	s := NewSession(nil)
	result, err := s.Prepare()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Log("prepare reported success")
	}
}

func TestCategorizeCommits(t *testing.T) {
	commits := []string{
		"feat: add new feature",
		"fix: resolve bug",
		"BREAKING: major change",
		"chore: update deps",
	}
	features, fixes, breaking := categorizeCommits(commits)
	if len(features) != 1 {
		t.Errorf("expected 1 feature, got %d", len(features))
	}
	if len(fixes) != 1 {
		t.Errorf("expected 1 fix, got %d", len(fixes))
	}
	if len(breaking) != 1 {
		t.Errorf("expected 1 breaking, got %d", len(breaking))
	}
}

func TestReleaseStatusConstants(t *testing.T) {
	if StatusPending != "PENDING" {
		t.Errorf("expected PENDING, got %s", StatusPending)
	}
	if StatusReady != "READY" {
		t.Errorf("expected READY, got %s", StatusReady)
	}
	if StatusFailed != "FAILED" {
		t.Errorf("expected FAILED, got %s", StatusFailed)
	}
}

func TestNotesFormatConstants(t *testing.T) {
	if NotesMarkdown != "markdown" {
		t.Errorf("expected markdown, got %s", NotesMarkdown)
	}
	if NotesJSON != "json" {
		t.Errorf("expected json, got %s", NotesJSON)
	}
	if NotesHTML != "html" {
		t.Errorf("expected html, got %s", NotesHTML)
	}
}

func TestRenderMarkdown(t *testing.T) {
	notes := &ReleaseNotes{
		Version:     "1.0.0",
		BuildNumber: "42",
		Date:        "2026-01-01",
		Summary:     "Test release",
		Features:    []string{"Feature 1"},
		Fixes:       []string{"Fix 1"},
	}
	output := renderMarkdown(notes)
	if output == "" {
		t.Error("expected non-empty markdown output")
	}
}

func TestRenderJSON(t *testing.T) {
	notes := &ReleaseNotes{
		Version: "1.0.0",
		BuildNumber: "42",
		Date:    "2026-01-01",
		Summary: "Test release",
	}
	output := renderJSON(notes)
	if output == "" {
		t.Error("expected non-empty JSON output")
	}
}

func TestRenderHTML(t *testing.T) {
	notes := &ReleaseNotes{
		Version: "1.0.0",
	}
	output := renderHTML(notes)
	if output == "" {
		t.Error("expected non-empty HTML output")
	}
}

func TestTimestamp(t *testing.T) {
	s := NewSession(nil)
	ts := s.Timestamp()
	if ts == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestProjectDir(t *testing.T) {
	s := NewSession(nil)
	if s.ProjectDir() != "" {
		t.Error("expected empty project dir initially")
	}
	s.SetProjectDir("/test")
	if s.ProjectDir() != "/test" {
		t.Errorf("expected /test, got %s", s.ProjectDir())
	}
}
