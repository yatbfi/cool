package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Editor represents an available editor on the system
type Editor struct {
	Name        string // Display name
	Command     string // Command to execute
	IsAvailable bool   // Whether the editor is installed
}

// DetectAvailableEditors scans the system for common text editors
func DetectAvailableEditors() []Editor {
	editors := []Editor{
		{Name: "Vim", Command: "vim"},
		{Name: "Nano", Command: "nano"},
		{Name: "Vi", Command: "vi"},
		{Name: "Emacs", Command: "emacs -nw"},
		{Name: "VS Code (wait mode)", Command: "code --wait"},
	}

	// Add Windows-specific editors
	if runtime.GOOS == "windows" {
		editors = append(editors, Editor{Name: "Notepad", Command: "notepad"})
	}

	// Check which editors are available
	available := []Editor{}
	for _, editor := range editors {
		// Extract the base command (first word)
		baseCmd := strings.Fields(editor.Command)[0]
		if _, err := exec.LookPath(baseCmd); err == nil {
			editor.IsAvailable = true
			available = append(available, editor)
		}
	}

	return available
}

// GetEditorCommand returns the editor command to use
// Priority: 1. configuredEditor, 2. $EDITOR, 3. $VISUAL, 4. first available editor
func GetEditorCommand(configuredEditor string) (string, error) {
	// If user has configured a preferred editor, use it
	if configuredEditor != "" && configuredEditor != "auto" {
		return configuredEditor, nil
	}

	// Check $EDITOR environment variable
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}

	// Check $VISUAL environment variable
	if visual := os.Getenv("VISUAL"); visual != "" {
		return visual, nil
	}

	// Find first available editor
	available := DetectAvailableEditors()
	if len(available) > 0 {
		return available[0].Command, nil
	}

	return "", fmt.Errorf("no text editor found. Please install vim, nano, or set the $EDITOR environment variable")
}

// OpenEditor opens a text editor for the user to input multiline text
func OpenEditor(editorCmd, prompt string) (string, error) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "cool-input-*.txt")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath) // Clean up

	// Write prompt to file
	if prompt != "" {
		promptText := fmt.Sprintf(";; %s\n;; Lines starting with ;; will be ignored\n\n", prompt)
		if _, err := tmpFile.WriteString(promptText); err != nil {
			tmpFile.Close()
			return "", fmt.Errorf("write prompt: %w", err)
		}
	}
	tmpFile.Close()

	// Parse editor command (might have flags like "code --wait")
	cmdParts := strings.Fields(editorCmd)
	if len(cmdParts) == 0 {
		return "", fmt.Errorf("invalid editor command: %s", editorCmd)
	}

	// Prepare command
	cmd := exec.Command(cmdParts[0], append(cmdParts[1:], tmpPath)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run editor
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	// Read the file content
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	// Filter out comment lines and trim
	lines := strings.Split(string(content), "\n")
	var resultLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, ";;") {
			resultLines = append(resultLines, line)
		}
	}

	result := strings.TrimSpace(strings.Join(resultLines, "\n"))
	return result, nil
}

// OpenEditorWithContent opens a text editor with initial content
func OpenEditorWithContent(editorCmd, initialContent string) (string, error) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "cool-edit-*.txt")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write initial content
	if _, err := tmpFile.WriteString(initialContent); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("write content: %w", err)
	}
	tmpFile.Close()

	// Parse editor command
	cmdParts := strings.Fields(editorCmd)
	if len(cmdParts) == 0 {
		return "", fmt.Errorf("invalid editor command: %s", editorCmd)
	}

	// Run editor
	cmd := exec.Command(cmdParts[0], append(cmdParts[1:], tmpPath)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	// Read result
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

// GetEditorDisplayName returns a friendly name for the editor command
func GetEditorDisplayName(editorCmd string) string {
	if editorCmd == "" || editorCmd == "auto" {
		return "Auto-detect"
	}

	// Extract base command
	baseCmd := filepath.Base(strings.Fields(editorCmd)[0])

	// Map to friendly names
	nameMap := map[string]string{
		"vim":     "Vim",
		"nvim":    "Neovim",
		"nano":    "Nano",
		"vi":      "Vi",
		"emacs":   "Emacs",
		"code":    "VS Code",
		"notepad": "Notepad",
	}

	if name, ok := nameMap[baseCmd]; ok {
		return name
	}

	return baseCmd
}
