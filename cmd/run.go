package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/domain/entity"
	"github.com/yatbfi/cool/internal/domain/usecase"
	"github.com/yatbfi/cool/internal/infrastructure/repository"
	"github.com/yatbfi/cool/internal/pkg/logger"
)

type RunCmd struct {
	*baseCmd
	path      string
	regex     string
	env       string
	debug     bool
	quiet     bool
	debugPort int
	command   []string
}

func NewRunCmd() *RunCmd {
	cmd := &RunCmd{}
	cobraCmd := &cobra.Command{
		Use:   "run [flags] -- [command]",
		Short: "Run Go application with hot reload",
		Long: `Run a Go application with automatic hot reload on file changes.

This command monitors file changes and automatically restarts your Go application. 
Perfect for development workflows.

Automatically loads .env file if present in the working directory.

Debugger Support:
  Use --debug-port to enable Delve debugger for IDE attachment (Goland/VSCode).
  The app will wait for debugger connection before starting.

Examples:
  # Basic hot reload
  cool run -- go run .

  # Run specific main.go
  cool run -- go run cmd/http/main.go

  # Run with debugger (Goland/VSCode)
  cool run --debug-port=2345 -- go run .

  # Run tests with hot reload  
  cool run -- go test -v ./...

  # Watch specific directory
  cool run --path=./cmd -- go run ./cmd

  # With additional environment variables
  cool run --env="PORT=8080,DEBUG=true" -- go run .

  # Custom file patterns
  cool run --regex=".*\\.go$,.*\\.html$" -- go run .`,
		RunE: cmd.run,
		Example: `  cool run -- go run .
  cool run -- go run cmd/http/main.go
  cool run --debug-port=2345 -- go run .
  cool run --path=./app -- go run ./app
  cool run --env="DB_HOST=localhost" -- go run .`,
	}

	cobraCmd.Flags().StringVarP(&cmd.path, "path", "p", ".", "Path to watch for changes")
	cobraCmd.Flags().StringVarP(&cmd.regex, "regex", "r", "", "Comma-separated regex patterns for files to watch")
	cobraCmd.Flags().StringVarP(&cmd.env, "env", "e", "", "Comma-separated KEY=VALUE pairs for environment variables")
	cobraCmd.Flags().BoolVarP(&cmd.debug, "debug", "d", false, "Enable debug mode")
	cobraCmd.Flags().BoolVarP(&cmd.quiet, "quiet", "q", false, "Disable logging output")
	cobraCmd.Flags().IntVar(&cmd.debugPort, "debug-port", 0, "Enable Delve debugger on specified port (e.g., 2345) for IDE attachment")

	cmd.baseCmd = newBaseCommand(cobraCmd)
	return cmd
}

func (c *RunCmd) run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		args = []string{"go", "run", "."}
	}

	// Transform command to use dlv if debug port is specified
	if c.debugPort > 0 {
		args = c.wrapWithDelve(args)
	}

	c.command = args

	// Get project root if available
	cfg := config.GetConfig()
	workDir := c.path
	if cfg.ProjectRoot != "" && c.path == "." {
		workDir = cfg.ProjectRoot
	}

	// Expand to absolute path
	absPath, err := filepath.Abs(workDir)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Get current working directory for running the command
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Parse environment variables from flags
	envVars := c.parseEnvVars()

	// Auto-load .env file from working directory (project root)
	envVars = c.loadEnvFile(currentDir, envVars)

	// Parse regex patterns
	patterns := c.parseRegexPatterns()

	// Show what we're doing
	fmt.Println()
	fmt.Println("🔥 Starting hot reload...")
	if c.debugPort > 0 {
		fmt.Printf("🐛 Debugger enabled on port %d\n", c.debugPort)
		fmt.Printf("   Connect your IDE debugger to localhost:%d\n", c.debugPort)
	}
	fmt.Printf("   Watching: %s\n", absPath)
	fmt.Printf("   Working directory: %s\n", currentDir)
	fmt.Printf("   Command: %s\n", strings.Join(args, " "))
	if len(patterns) > 0 {
		fmt.Printf("   Patterns: %s\n", strings.Join(patterns, ", "))
	}
	fmt.Println()
	if c.debugPort > 0 {
		fmt.Println("⏳ Waiting for debugger to attach...")
	}
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("=====================================")
	fmt.Println()

	// Create watch configuration
	watchConfig := &entity.WatchConfig{
		Path:         absPath,    // Watch this directory for changes
		WorkDir:      currentDir, // Run command from current directory
		Patterns:     patterns,
		EnvVars:      envVars,
		Command:      c.command,
		PollInterval: 500 * time.Millisecond,
		RestartDelay: time.Second,
		Debug:        c.debug,
		Quiet:        c.quiet,
	}

	// Create logger
	log := logger.NewSimple(c.quiet)

	// Create infrastructure implementations
	fileWatcher, err := repository.NewPollingFileWatcher(watchConfig)
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	processRunner := repository.NewCommandProcessRunner(watchConfig)

	// Create usecase
	hotReload := usecase.NewHotReloadUsecase(watchConfig, fileWatcher, processRunner, log)

	// Run hot reload
	return hotReload.Run(context.Background())
}

func (c *RunCmd) parseEnvVars() map[string]string {
	envVars := make(map[string]string)
	if c.env == "" {
		return envVars
	}

	pairs := strings.Split(c.env, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}
	return envVars
}

func (c *RunCmd) parseRegexPatterns() []string {
	if c.regex == "" {
		// Default patterns for Go projects
		return []string{`\.go$`, `go\.mod$`, `go\.sum$`}
	}

	patterns := strings.Split(c.regex, ",")
	result := make([]string, 0, len(patterns))
	for _, p := range patterns {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// loadEnvFile loads .env file if present and merges with existing env vars
// Priority: flag env vars > .env file
func (c *RunCmd) loadEnvFile(basePath string, flagEnvVars map[string]string) map[string]string {
	envFile := filepath.Join(basePath, ".env")

	// Check if .env exists
	content, err := os.ReadFile(envFile)
	if err != nil {
		// .env file not found or cannot be read, just return flag vars
		return flagEnvVars
	}

	if !c.quiet {
		fmt.Printf("📄 Loading .env file from %s\n", envFile)
	}

	// Parse .env file
	envVars := make(map[string]string)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}

			envVars[key] = value
		}
	}

	// Merge with flag vars (flag vars have priority)
	for k, v := range flagEnvVars {
		envVars[k] = v
	}

	return envVars
}

// wrapWithDelve wraps the go command with Delve debugger
func (c *RunCmd) wrapWithDelve(args []string) []string {
	// Check if command is "go run"
	if len(args) >= 2 && args[0] == "go" && args[1] == "run" {
		// Transform: go run <path> -> dlv debug <path> --headless --listen=:port --api-version=2 --accept-multiclient
		dlvArgs := []string{
			"dlv",
			"debug",
		}

		// Add the package path (skip "go run")
		if len(args) > 2 {
			dlvArgs = append(dlvArgs, args[2:]...)
		}

		// Add dlv flags
		dlvArgs = append(dlvArgs,
			"--headless",
			fmt.Sprintf("--listen=:%d", c.debugPort),
			"--api-version=2",
			"--accept-multiclient",
			"--continue",
		)

		return dlvArgs
	}

	// For other commands, wrap with dlv exec
	// Example: dlv exec --headless --listen=:port --api-version=2 <binary>
	return []string{
		"dlv",
		"exec",
		"--headless",
		fmt.Sprintf("--listen=:%d", c.debugPort),
		"--api-version=2",
		"--accept-multiclient",
		"--continue",
		"--",
		strings.Join(args, " "),
	}
}
