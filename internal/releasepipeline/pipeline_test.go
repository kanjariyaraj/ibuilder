package releasepipeline

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPipeline(t *testing.T) {
	p := NewPipeline(nil)
	if p == nil {
		t.Fatal("expected pipeline, got nil")
	}
	if p.Mode() != ModeBeta {
		t.Errorf("expected ModeBeta, got %s", p.Mode())
	}
}

func TestSetMode(t *testing.T) {
	p := NewPipeline(nil)
	p.SetMode(ModeProduction)
	if p.Mode() != ModeProduction {
		t.Errorf("expected ModeProduction, got %s", p.Mode())
	}
}

func TestDryRun(t *testing.T) {
	p := NewPipeline(nil)
	if p.IsDryRun() {
		t.Error("expected dry run false by default")
	}
	p.SetDryRun(true)
	if !p.IsDryRun() {
		t.Error("expected dry run true")
	}
}

func TestProjectDir(t *testing.T) {
	p := NewPipeline(nil)
	if p.ProjectDir() != "" {
		t.Error("expected empty project dir")
	}
	p.SetProjectDir("/test")
	if p.ProjectDir() != "/test" {
		t.Errorf("expected /test, got %s", p.ProjectDir())
	}
}

func TestResults(t *testing.T) {
	p := NewPipeline(nil)
	results := p.Results()
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestAddResult(t *testing.T) {
	p := NewPipeline(nil)
	r := p.addResult(StageValidate, true, "ok", nil)
	if r == nil {
		t.Fatal("expected result, got nil")
	}
	if r.Stage != StageValidate {
		t.Errorf("expected StageValidate, got %s", r.Stage)
	}
	if !r.Success {
		t.Error("expected success")
	}

	results := p.Results()
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestAddResultWithError(t *testing.T) {
	p := NewPipeline(nil)
	r := p.addResult(StageBuild, false, "failed", fmt.Errorf("test error"))
	if r.Success {
		t.Error("expected failure")
	}
	if r.Error == "" {
		t.Error("expected non-empty error")
	}
}

func TestValidateStage(t *testing.T) {
	p := NewPipeline(nil)
	r := p.Validate()
	if r == nil {
		t.Fatal("expected result, got nil")
	}
	if r.Stage != StageValidate {
		t.Errorf("expected StageValidate, got %s", r.Stage)
	}
}

func TestBuildStage(t *testing.T) {
	p := NewPipeline(nil)
	p.SetDryRun(true)
	r := p.Build()
	if r == nil {
		t.Fatal("expected result, got nil")
	}
	if !r.Success {
		t.Errorf("expected success in dry run, got: %s", r.Message)
	}
}

func TestSignStage(t *testing.T) {
	p := NewPipeline(nil)
	p.SetDryRun(true)
	r := p.Sign()
	if r == nil {
		t.Fatal("expected result, got nil")
	}
	if r.Stage != StageSign {
		t.Errorf("expected StageSign, got %s", r.Stage)
	}
}

func TestGenerateNotesStage(t *testing.T) {
	p := NewPipeline(nil)
	p.SetDryRun(true)
	r := p.GenerateNotes()
	if r == nil {
		t.Fatal("expected result, got nil")
	}
	if !r.Success {
		t.Errorf("expected success in dry run, got: %s", r.Message)
	}
}

func TestUploadStage(t *testing.T) {
	p := NewPipeline(nil)
	p.SetDryRun(true)
	r := p.Upload()
	if r == nil {
		t.Fatal("expected result, got nil")
	}
	if !r.Success {
		t.Errorf("expected success in dry run, got: %s", r.Message)
	}
}

func TestGitHubReleaseStage(t *testing.T) {
	p := NewPipeline(nil)
	p.SetDryRun(true)
	r := p.CreateGitHubRelease()
	if r == nil {
		t.Fatal("expected result, got nil")
	}
	if !r.Success {
		t.Errorf("expected success in dry run, got: %s", r.Message)
	}
}

func TestGenerateReportStage(t *testing.T) {
	p := NewPipeline(nil)
	started := time.Now()
	r := p.GenerateReport(started)
	if r == nil {
		t.Fatal("expected result, got nil")
	}
	if !r.Success {
		t.Errorf("expected success, got: %s", r.Message)
	}
}

func TestFullPipeline(t *testing.T) {
	p := NewPipeline(nil)
	p.SetDryRun(true)
	p.StartStatus()

	p.Validate()
	p.Build()
	p.Sign()
	p.GenerateNotes()
	p.Upload()
	p.CreateGitHubRelease()
	p.GenerateReport(time.Now())

	p.FinishStatus()

	results := p.Results()
	if len(results) != 7 {
		t.Errorf("expected 7 stage results, got %d", len(results))
	}

	skipFailure := map[PipelineStage]bool{
		StageValidate: true,
	}
	for _, r := range results {
		if !r.Success && !skipFailure[r.Stage] {
			t.Errorf("stage %s failed: %s", r.Stage, r.Message)
		}
	}
}

func TestStatus(t *testing.T) {
	status := GetStatus()
	if status == nil {
		t.Fatal("expected status, got nil")
	}
}

func TestStatusSummary(t *testing.T) {
	summary := StatusSummary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestTimestamp(t *testing.T) {
	p := NewPipeline(nil)
	ts := p.Timestamp()
	if ts == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestModeConstants(t *testing.T) {
	if ModeBeta != "beta" {
		t.Errorf("expected beta, got %s", ModeBeta)
	}
	if ModeProduction != "production" {
		t.Errorf("expected production, got %s", ModeProduction)
	}
	if ModeInternal != "internal" {
		t.Errorf("expected internal, got %s", ModeInternal)
	}
}

func TestStageConstants(t *testing.T) {
	if StageValidate != "validate" {
		t.Errorf("expected validate, got %s", StageValidate)
	}
	if StageBuild != "build" {
		t.Errorf("expected build, got %s", StageBuild)
	}
	if StageReport != "report" {
		t.Errorf("expected report, got %s", StageReport)
	}
}
