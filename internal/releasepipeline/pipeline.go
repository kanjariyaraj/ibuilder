package releasepipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kanjariyaraj/Builder/internal/logger"
)

type ReleaseMode string

const (
	ModeBeta       ReleaseMode = "beta"
	ModeProduction ReleaseMode = "production"
	ModeInternal   ReleaseMode = "internal"
	ModeCustom     ReleaseMode = "custom"
)

type PipelineStage string

const (
	StageValidate  PipelineStage = "validate"
	StageBuild     PipelineStage = "build"
	StageSign      PipelineStage = "sign"
	StageNotes     PipelineStage = "notes"
	StageUpload    PipelineStage = "upload"
	StageRelease   PipelineStage = "release"
	StageReport    PipelineStage = "report"
)

type StageResult struct {
	Stage   PipelineStage `json:"stage"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Started time.Time     `json:"started_at"`
	Done    time.Time     `json:"done_at"`
	Error   string        `json:"error,omitempty"`
}

type PipelineResult struct {
	Success     bool                   `json:"success"`
	Mode        ReleaseMode            `json:"mode"`
	DryRun      bool                   `json:"dry_run"`
	Stages      []StageResult          `json:"stages"`
	Started     time.Time              `json:"started_at"`
	Completed   time.Time              `json:"completed_at"`
	Artifact    string                 `json:"artifact,omitempty"`
	ReleaseTag  string                 `json:"release_tag,omitempty"`
	Summary     string                 `json:"summary"`
}

type Pipeline struct {
	mu         sync.RWMutex
	log        *logger.Logger
	projectDir string
	mode       ReleaseMode
	dryRun     bool
	results    []StageResult
}

func NewPipeline(log *logger.Logger) *Pipeline {
	return &Pipeline{
		log:     log,
		mode:    ModeBeta,
		dryRun:  false,
		results: []StageResult{},
	}
}

func (p *Pipeline) SetProjectDir(dir string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.projectDir = dir
}

func (p *Pipeline) ProjectDir() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.projectDir
}

func (p *Pipeline) SetMode(mode ReleaseMode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.mode = mode
}

func (p *Pipeline) Mode() ReleaseMode {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.mode
}

func (p *Pipeline) SetDryRun(dryRun bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.dryRun = dryRun
}

func (p *Pipeline) IsDryRun() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dryRun
}

func (p *Pipeline) logInfo(msg string, args ...any) {
	if p.log == nil {
		return
	}
	p.log.Info(msg, args...)
}

func (p *Pipeline) logWarn(msg string, args ...any) {
	if p.log == nil {
		return
	}
	p.log.Warn(msg, args...)
}

func (p *Pipeline) Timestamp() string {
	return time.Now().Format("20060102_150405")
}

func (p *Pipeline) addResult(stage PipelineStage, success bool, message string, err error) *StageResult {
	r := &StageResult{
		Stage:   stage,
		Success: success,
		Message: message,
		Started: time.Now(),
		Done:    time.Now(),
	}
	if err != nil {
		r.Error = err.Error()
	}
	p.mu.Lock()
	p.results = append(p.results, *r)
	p.mu.Unlock()
	return r
}

func (p *Pipeline) findIPA() (string, error) {
	dir := p.ProjectDir()
	if dir == "" {
		return "", fmt.Errorf("no project directory set")
	}

	buildDir := filepath.Join(dir, ".build")
	var ipaPath string
	err := filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".ipa") {
			ipaPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if ipaPath == "" {
		return "", fmt.Errorf("no .ipa file found in .build directory")
	}
	return ipaPath, nil
}

func (p *Pipeline) Results() []StageResult {
	p.mu.RLock()
	defer p.mu.RUnlock()
	res := make([]StageResult, len(p.results))
	copy(res, p.results)
	return res
}

func (p *Pipeline) ensureDirs() error {
	dirs := []string{
		filepath.Join(p.ProjectDir(), ".build", "logs"),
		filepath.Join(p.ProjectDir(), ".build", "reports", "release"),
		filepath.Join(p.ProjectDir(), ".build", "releases"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}
	return nil
}
