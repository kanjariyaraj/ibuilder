package flutter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func (s *Session) WatchMode() (*FileWatcher, error) {
	s.mu.RLock()
	if s.projectDir == "" {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no project directory set")
	}
	s.mu.RUnlock()

	debounceMs := s.cfg.DebounceMs
	if debounceMs <= 0 {
		debounceMs = 500
	}

	watcher := NewFileWatcher(s.projectDir, debounceMs)

	if err := watcher.Start(); err != nil {
		return nil, fmt.Errorf("failed to start file watcher: %w", err)
	}

	s.mu.Lock()
	s.state = SessionWatchMode
	s.mu.Unlock()

	s.log.Info("watch mode started", "dir", s.projectDir, "debounce_ms", debounceMs)

	go func() {
		for event := range watcher.Events() {
			if len(event.Changes) > 0 {
				s.log.Info("file changes detected",
					"count", len(event.Changes),
					"first", event.Changes[0].Path)
			}
		}
	}()

	return watcher, nil
}

type FileChange struct {
	Path      string    `json:"path"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}

type WatchEvent struct {
	Changes  []FileChange `json:"changes"`
	Reloaded bool         `json:"reloaded"`
	Error    string       `json:"error,omitempty"`
}

type FileWatcher struct {
	mu         sync.Mutex
	dir        string
	debounceMs int
	events     chan WatchEvent
	stopChan   chan struct{}
	running    bool
	ignoreDirs []string
	lastChange time.Time
	timer      *time.Timer
}

func NewFileWatcher(dir string, debounceMs int) *FileWatcher {
	return &FileWatcher{
		dir:        dir,
		debounceMs: debounceMs,
		events:     make(chan WatchEvent, 100),
		stopChan:   make(chan struct{}),
		ignoreDirs: []string{".dart_tool", ".pub-cache", "build", ".git", "ios/Pods", "android/.gradle", "android/build", "node_modules"},
	}
}

func (w *FileWatcher) Start() error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return fmt.Errorf("watcher already running")
	}
	w.running = true
	w.mu.Unlock()

	go w.watchLoop()

	return nil
}

func (w *FileWatcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return
	}

	close(w.stopChan)
	w.running = false

	if w.timer != nil {
		w.timer.Stop()
	}
}

func (w *FileWatcher) Events() chan WatchEvent {
	return w.events
}

func (w *FileWatcher) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

func (w *FileWatcher) watchLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	lastSnapshot := w.snapshotFiles()

	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			currentSnapshot := w.snapshotFiles()
			changes := w.detectChanges(lastSnapshot, currentSnapshot)

			if len(changes) > 0 {
				now := time.Now()
				w.mu.Lock()
				w.lastChange = now
				if w.timer != nil {
					w.timer.Stop()
				}
				w.timer = time.AfterFunc(time.Duration(w.debounceMs)*time.Millisecond, func() {
					w.mu.Lock()
					w.events <- WatchEvent{
						Changes:  changes,
						Reloaded: false,
					}
					w.mu.Unlock()
				})
				w.mu.Unlock()
			}

			lastSnapshot = currentSnapshot
		}
	}
}

func (w *FileWatcher) snapshotFiles() map[string]time.Time {
	snapshot := make(map[string]time.Time)

	filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, err := filepath.Rel(w.dir, path)
		if err != nil {
			return nil
		}

		if info.IsDir() {
			for _, ignore := range w.ignoreDirs {
				if strings.HasPrefix(rel, ignore) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".dart", ".yaml", ".xml", ".plist", ".json", ".gradle", ".properties":
			snapshot[rel] = info.ModTime()
		}

		return nil
	})

	return snapshot
}

func (w *FileWatcher) detectChanges(old, new map[string]time.Time) []FileChange {
	var changes []FileChange

	for path, newTime := range new {
		oldTime, exists := old[path]
		if !exists {
			changes = append(changes, FileChange{
				Path:      path,
				Operation: "created",
				Timestamp: newTime,
			})
		} else if !newTime.Equal(oldTime) {
			changes = append(changes, FileChange{
				Path:      path,
				Operation: "modified",
				Timestamp: newTime,
			})
		}
	}

	for path := range old {
		if _, exists := new[path]; !exists {
			changes = append(changes, FileChange{
				Path:      path,
				Operation: "deleted",
				Timestamp: time.Now(),
			})
		}
	}

	return changes
}
