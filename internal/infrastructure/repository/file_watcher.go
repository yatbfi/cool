package repository

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/yatbfi/cool/internal/domain/entity"
	"github.com/yatbfi/cool/internal/domain/repository"
)

// PollingFileWatcher implements FileWatcher using polling strategy
type PollingFileWatcher struct {
	config       *entity.WatchConfig
	fileStates   map[string]time.Time
	regexps      []*regexp.Regexp
	mu           sync.RWMutex
	pollInterval time.Duration
}

// NewPollingFileWatcher creates a new polling-based file watcher
func NewPollingFileWatcher(config *entity.WatchConfig) (repository.FileWatcher, error) {
	// Compile regex patterns
	regexps := make([]*regexp.Regexp, 0, len(config.Patterns))
	for _, pattern := range config.Patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		regexps = append(regexps, re)
	}

	pollInterval := config.PollInterval
	if pollInterval == 0 {
		pollInterval = 500 * time.Millisecond
	}

	return &PollingFileWatcher{
		config:       config,
		fileStates:   make(map[string]time.Time),
		regexps:      regexps,
		pollInterval: pollInterval,
	}, nil
}

// Watch starts watching for file changes
func (w *PollingFileWatcher) Watch(ctx context.Context) (<-chan entity.FileChange, <-chan error) {
	changeChan := make(chan entity.FileChange, 10)
	errorChan := make(chan error, 1)

	// Initial scan
	w.scanFiles()

	go func() {
		defer close(changeChan)
		defer close(errorChan)

		ticker := time.NewTicker(w.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				changes := w.checkChanges()
				for _, change := range changes {
					select {
					case changeChan <- change:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return changeChan, errorChan
}

// Stop stops the file watcher
func (w *PollingFileWatcher) Stop() error {
	// Polling watcher stops via context cancellation
	return nil
}

// scanFiles scans all files and records their modification times
func (w *PollingFileWatcher) scanFiles() {
	w.mu.Lock()
	defer w.mu.Unlock()

	filepath.Walk(w.config.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Skip hidden files and directories
		if w.isHidden(path) {
			return nil
		}

		// Check if file matches patterns
		if w.matchesPattern(path) {
			w.fileStates[path] = info.ModTime()
		}

		return nil
	})
}

// checkChanges checks for file modifications
func (w *PollingFileWatcher) checkChanges() []entity.FileChange {
	w.mu.Lock()
	defer w.mu.Unlock()

	var changes []entity.FileChange
	currentFiles := make(map[string]time.Time)

	// Scan for new and modified files
	filepath.Walk(w.config.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if w.isHidden(path) {
			return nil
		}

		if w.matchesPattern(path) {
			currentFiles[path] = info.ModTime()

			lastMod, exists := w.fileStates[path]
			if !exists {
				changes = append(changes, entity.FileChange{
					Path:      path,
					EventType: entity.EventCreated,
					Timestamp: time.Now(),
				})
			} else if info.ModTime().After(lastMod) {
				changes = append(changes, entity.FileChange{
					Path:      path,
					EventType: entity.EventModified,
					Timestamp: time.Now(),
				})
			}
		}

		return nil
	})

	// Check for deleted files
	for path := range w.fileStates {
		if _, exists := currentFiles[path]; !exists {
			changes = append(changes, entity.FileChange{
				Path:      path,
				EventType: entity.EventDeleted,
				Timestamp: time.Now(),
			})
		}
	}

	// Update file states
	w.fileStates = currentFiles

	return changes
}

// matchesPattern checks if file matches any of the patterns
func (w *PollingFileWatcher) matchesPattern(path string) bool {
	for _, re := range w.regexps {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

// isHidden checks if a path contains hidden files or directories
func (w *PollingFileWatcher) isHidden(path string) bool {
	return strings.Contains(path, "/.") || strings.Contains(path, "\\.")
}
