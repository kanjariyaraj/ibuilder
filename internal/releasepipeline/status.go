package releasepipeline

import (
	"fmt"
	"sync"
	"time"
)

type PipelineStatus struct {
	mu         sync.RWMutex
	Running    bool          `json:"running"`
	Mode       ReleaseMode   `json:"mode"`
	DryRun     bool          `json:"dry_run"`
	Current    PipelineStage `json:"current_stage"`
	Progress   string        `json:"progress"`
	StartedAt  time.Time     `json:"started_at"`
	StagesDone int           `json:"stages_done"`
	TotalStages int          `json:"total_stages"`
}

var globalStatus = &PipelineStatus{
	TotalStages: 7,
}

func GetStatus() *PipelineStatus {
	globalStatus.mu.RLock()
	defer globalStatus.mu.RUnlock()
	return globalStatus
}

func (p *Pipeline) StartStatus() {
	globalStatus.mu.Lock()
	globalStatus.Running = true
	globalStatus.Mode = p.mode
	globalStatus.DryRun = p.dryRun
	globalStatus.Current = StageValidate
	globalStatus.Progress = "starting"
	globalStatus.StartedAt = time.Now()
	globalStatus.StagesDone = 0
	globalStatus.mu.Unlock()
}

func (p *Pipeline) UpdateStatus(stage PipelineStage, progress string) {
	globalStatus.mu.Lock()
	globalStatus.Current = stage
	globalStatus.Progress = progress
	globalStatus.StagesDone++
	globalStatus.mu.Unlock()
}

func (p *Pipeline) FinishStatus() {
	globalStatus.mu.Lock()
	globalStatus.Running = false
	globalStatus.Progress = "completed"
	globalStatus.mu.Unlock()
}

func StatusSummary() string {
	s := GetStatus()
	if !s.Running {
		return "No release pipeline currently running."
	}
	return fmt.Sprintf(`Release Pipeline Status:
  Running:   yes
  Mode:      %s
  Dry Run:   %v
  Stage:     %s (%d/%d)
  Progress:  %s
  Started:   %s`,
		s.Mode, s.DryRun, s.Current, s.StagesDone, s.TotalStages,
		s.Progress, s.StartedAt.Format(time.RFC3339))
}
