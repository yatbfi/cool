package repository

import (
	"context"

	"github.com/yatbfi/cool/internal/domain/entity"
)

// FileWatcher defines the interface for watching file changes
type FileWatcher interface {
	// Watch starts watching for file changes
	Watch(ctx context.Context) (<-chan entity.FileChange, <-chan error)

	// Stop stops the file watcher
	Stop() error
}

// ProcessRunner defines the interface for running and managing processes
type ProcessRunner interface {
	// Start starts a new process
	Start(ctx context.Context) error

	// Stop stops the running process
	Stop() error

	// IsRunning checks if the process is currently running
	IsRunning() bool

	// Info returns information about the running process
	Info() *entity.ProcessInfo
}
