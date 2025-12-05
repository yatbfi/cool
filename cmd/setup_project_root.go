package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/pkg/common"
)

type SetupProjectRootCmd struct {
	*baseCmd
}

func NewSetupProjectRootCmd() *SetupProjectRootCmd {
	cmd := &SetupProjectRootCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "project-root",
		Short: "Setup project root directory",
		Long: `Configure the project root directory for your workspace.

This is typically the base directory where your projects are located,
such as $GOPATH/src/bfi-finance or similar.

The project root will be used by other features that need to know
the base location of your projects.`,
		RunE: cmd.run,
	})
	return cmd
}

func (c *SetupProjectRootCmd) run(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	// Detect default project root
	defaultRoot := detectDefaultProjectRoot()
	if cfg.ProjectRoot != "" {
		defaultRoot = cfg.ProjectRoot
	}

	fmt.Println("📁 Setup Project Root Directory")
	fmt.Println()

	if cfg.ProjectRoot != "" {
		fmt.Printf("Current project root: %s\n", cfg.ProjectRoot)
		fmt.Println()
	}

	// Prompt for project root
	projectRootPrompt := promptui.Prompt{
		Label:   "Project Root Directory",
		Default: defaultRoot,
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("project root cannot be empty")
			}
			// Expand ~ to home directory
			if strings.HasPrefix(input, "~") {
				home, err := os.UserHomeDir()
				if err == nil {
					input = filepath.Join(home, input[1:])
				}
			}
			// Check if directory exists
			if _, err := os.Stat(input); os.IsNotExist(err) {
				return fmt.Errorf("directory does not exist: %s", input)
			}
			return nil
		},
	}

	projectRoot, err := projectRootPrompt.Run()
	if err != nil {
		return err
	}

	// Expand ~ to home directory
	if strings.HasPrefix(projectRoot, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home directory: %w", err)
		}
		projectRoot = filepath.Join(home, projectRoot[1:])
	}

	// Clean the path
	projectRoot = filepath.Clean(projectRoot)

	// Save to config
	cfg.ProjectRoot = projectRoot
	if err := config.SaveLocalConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("\n✅ Project root setup complete!\n")
	fmt.Printf("Project Root: %s\n", projectRoot)

	return nil
}

// detectDefaultProjectRoot attempts to detect a sensible default project root
func detectDefaultProjectRoot() string {
	// Try $GOPATH/src/bfi-finance
	if gopath, err := common.GetGopath(); err == nil {
		bfiPath := filepath.Join(gopath, "src", "bfi-finance")
		if _, err := os.Stat(bfiPath); err == nil {
			return bfiPath
		}

		// Try just $GOPATH/src
		srcPath := filepath.Join(gopath, "src")
		if _, err := os.Stat(srcPath); err == nil {
			return srcPath
		}
	}

	// Try ~/projects
	if home, err := os.UserHomeDir(); err == nil {
		projectsPath := filepath.Join(home, "projects")
		if _, err := os.Stat(projectsPath); err == nil {
			return projectsPath
		}
	}

	// Default to current directory
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}

	return ""
}
