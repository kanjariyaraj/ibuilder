package flutter

import (
	"fmt"
	"time"
)

type RecoveryResult struct {
	Success   bool          `json:"success"`
	Action    string        `json:"action"`
	Duration  time.Duration `json:"duration_ms"`
	Attempts  int           `json:"attempts"`
	Error     string        `json:"error,omitempty"`
}

func (s *Session) Recover() *RecoveryResult {
	s.mu.Lock()
	s.state = SessionRecovering
	s.mu.Unlock()

	s.log.Info("starting flutter session recovery")
	start := time.Now()

	var lastErr error
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		s.log.Info("recovery attempt", "attempt", attempt, "max", maxAttempts)

		if err := s.ResolveDependencies(); err != nil {
			lastErr = err
			s.log.Error("recovery dependency resolution failed", "attempt", attempt, "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if _, err := s.DevMode(); err != nil {
			lastErr = err
			s.log.Error("recovery dev mode failed", "attempt", attempt, "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		s.mu.Lock()
		s.state = SessionActive
		s.mu.Unlock()

		duration := time.Since(start)
		result := &RecoveryResult{
			Success:  true,
			Action:   "recovered",
			Duration: duration,
			Attempts: attempt,
		}

		s.log.Info("flutter session recovered", "attempts", attempt, "duration", duration)
		return result
	}

	s.mu.Lock()
	s.state = SessionInactive
	s.mu.Unlock()

	duration := time.Since(start)
	errMsg := fmt.Sprintf("recovery failed after %d attempts", maxAttempts)
	if lastErr != nil {
		errMsg = fmt.Sprintf("%s: %v", errMsg, lastErr)
	}

	result := &RecoveryResult{
		Success:  false,
		Action:   "recovery_failed",
		Duration: duration,
		Attempts: maxAttempts,
		Error:    errMsg,
	}

	s.log.Error("flutter session recovery failed", "attempts", maxAttempts, "error", lastErr)
	return result
}
