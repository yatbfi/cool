package usecase

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/yatbfi/cool/internal/domain/entity"
	"github.com/yatbfi/cool/internal/domain/repository"
)

// HotReload defines the interface for hot reload usecase
type HotReload interface {
	// Run starts the hot reload loop
	Run(ctx context.Context) error
}

// hotReloadUsecase implements HotReload interface
type hotReloadUsecase struct {
	config        *entity.WatchConfig
	fileWatcher   repository.FileWatcher
	processRunner repository.ProcessRunner
	logger        Logger
}

// Logger defines the logging interface
type Logger interface {
	Info(message string)
	Success(message string)
	Error(message string)
	Debug(message string)
}

// NewHotReloadUsecase creates a new hot reload usecase
func NewHotReloadUsecase(
	config *entity.WatchConfig,
	fileWatcher repository.FileWatcher,
	processRunner repository.ProcessRunner,
	logger Logger,
) HotReload {
	return &hotReloadUsecase{
		config:        config,
		fileWatcher:   fileWatcher,
		processRunner: processRunner,
		logger:        logger,
	}
}

// Run starts the hot reload loop
func (uc *hotReloadUsecase) Run(ctx context.Context) error {
	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	// Start initial process
	uc.logger.Success(fmt.Sprintf("▶️  Running: %s", strings.Join(uc.config.Command, " ")))
	if err := uc.processRunner.Start(ctx); err != nil {
		return fmt.Errorf("failed to start initial process: %w", err)
	}

	// Start file watcher
	changeChan, errorChan := uc.fileWatcher.Watch(ctx)

	var lastRestart time.Time
	restartDelay := uc.config.RestartDelay
	if restartDelay == 0 {
		restartDelay = time.Second
	}

	for {
		select {
		case change, ok := <-changeChan:
			if !ok {
				return nil
			}

			// Debounce restarts
			if time.Since(lastRestart) < restartDelay {
				continue
			}

			if uc.config.Debug {
				uc.logger.Debug(fmt.Sprintf("%s: %s", change.EventType, change.Path))
			}

			// Restart process
			uc.logger.Info("📝 Changes detected, restarting...")
			if err := uc.restartProcess(ctx); err != nil {
				uc.logger.Error(fmt.Sprintf("Failed to restart: %v", err))
			} else {
				lastRestart = time.Now()
			}

		case err := <-errorChan:
			if err != nil {
				uc.logger.Error(fmt.Sprintf("Watch error: %v", err))
			}

		case <-sigChan:
			uc.logger.Info("🛑 Stopping...")
			if err := uc.processRunner.Stop(); err != nil {
				uc.logger.Error(fmt.Sprintf("Failed to stop process: %v", err))
			}
			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// restartProcess stops and starts the process
func (uc *hotReloadUsecase) restartProcess(ctx context.Context) error {
	// Stop current process
	if err := uc.processRunner.Stop(); err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}

	// Start new process
	uc.logger.Success(fmt.Sprintf("▶️  Running: %s", strings.Join(uc.config.Command, " ")))
	if err := uc.processRunner.Start(ctx); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	return nil
}
