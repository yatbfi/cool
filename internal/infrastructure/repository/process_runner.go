package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/yatbfi/cool/internal/domain/entity"
	"github.com/yatbfi/cool/internal/domain/repository"
)

// CommandProcessRunner implements ProcessRunner for executing commands
type CommandProcessRunner struct {
	config  *entity.WatchConfig
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	mu      sync.Mutex
	started time.Time
}

// NewCommandProcessRunner creates a new command process runner
func NewCommandProcessRunner(config *entity.WatchConfig) repository.ProcessRunner {
	return &CommandProcessRunner{
		config: config,
	}
}

// Start starts a new process
func (r *CommandProcessRunner) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Stop existing process if running
	if r.cmd != nil && r.cmd.Process != nil {
		if r.cancel != nil {
			r.cancel()
		}
		r.cmd.Process.Kill()
		r.cmd.Wait()
	}

	// Create new context
	cmdCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel

	// Create command
	r.cmd = exec.CommandContext(cmdCtx, r.config.Command[0], r.config.Command[1:]...)
	// Use WorkDir if specified, otherwise use Path
	if r.config.WorkDir != "" {
		r.cmd.Dir = r.config.WorkDir
	} else {
		r.cmd.Dir = r.config.Path
	}
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr
	r.cmd.Stdin = os.Stdin

	// Set environment variables
	r.cmd.Env = os.Environ()
	for key, val := range r.config.EnvVars {
		r.cmd.Env = append(r.cmd.Env, fmt.Sprintf("%s=%s", key, val))
	}

	// Start the command
	r.started = time.Now()
	if err := r.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	return nil
}

// Stop stops the running process
func (r *CommandProcessRunner) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cancel != nil {
		r.cancel()
	}

	if r.cmd != nil && r.cmd.Process != nil {
		if err := r.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
		r.cmd.Wait()
	}

	return nil
}

// IsRunning checks if the process is currently running
func (r *CommandProcessRunner) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.cmd != nil && r.cmd.Process != nil && r.cmd.ProcessState == nil
}

// Info returns information about the running process
func (r *CommandProcessRunner) Info() *entity.ProcessInfo {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cmd == nil || r.cmd.Process == nil {
		return nil
	}

	return &entity.ProcessInfo{
		PID:     r.cmd.Process.Pid,
		Command: r.config.Command,
		Started: r.started,
	}
}
