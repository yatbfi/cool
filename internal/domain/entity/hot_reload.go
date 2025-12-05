package entity

import "time"

// WatchConfig holds the configuration for file watching
type WatchConfig struct {
	Path         string            // Directory to watch for changes
	WorkDir      string            // Working directory where command runs (usually project root with go.mod)
	Patterns     []string          // Regex patterns to match files
	EnvVars      map[string]string // Environment variables
	Command      []string          // Command to execute
	PollInterval time.Duration     // How often to check for changes
	RestartDelay time.Duration     // Delay before restart after change
	Debug        bool              // Enable debug logging
	Quiet        bool              // Disable logging
}

// FileChange represents a file change event
type FileChange struct {
	Path      string    // File path
	EventType EventType // Type of change
	Timestamp time.Time // When the change occurred
}

// EventType represents the type of file change
type EventType int

const (
	EventCreated EventType = iota
	EventModified
	EventDeleted
)

func (e EventType) String() string {
	switch e {
	case EventCreated:
		return "Created"
	case EventModified:
		return "Modified"
	case EventDeleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}

// ProcessInfo holds information about a running process
type ProcessInfo struct {
	PID     int
	Command []string
	Started time.Time
}
