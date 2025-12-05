package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/pkg/common"
)

type SetupEditorCmd struct {
	*baseCmd
}

func NewSetupEditorCmd() *SetupEditorCmd {
	cmd := &SetupEditorCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "editor",
		Short: "Setup preferred text editor for multiline input",
		Long: `Configure your preferred text editor for entering multiline descriptions.

The editor will be used when submitting review requests to allow 
multiline descriptions instead of single-line terminal input.

Available options:
- Auto-detect: Use $EDITOR environment variable or first available editor
- Specific editors: Choose from detected editors on your system`,
		RunE: cmd.run,
	})
	return cmd
}

func (c *SetupEditorCmd) run(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	fmt.Println("🔍 Detecting available text editors...")
	fmt.Println()

	// Detect available editors
	availableEditors := common.DetectAvailableEditors()

	if len(availableEditors) == 0 {
		fmt.Println("⚠️  No common text editors found on your system.")
		fmt.Println("\nPlease install one of the following:")
		fmt.Println("  - vim")
		fmt.Println("  - nano")
		fmt.Println("  - vi")
		fmt.Println("  - emacs")
		fmt.Println("\nOr set the $EDITOR environment variable to your preferred editor.")
		return fmt.Errorf("no text editors available")
	}

	// Build selection items
	items := []string{"Auto-detect (use $EDITOR or first available)"}
	editorCommands := []string{"auto"}

	for _, editor := range availableEditors {
		items = append(items, fmt.Sprintf("%s (%s)", editor.Name, editor.Command))
		editorCommands = append(editorCommands, editor.Command)
	}

	// Current selection
	currentEditor := cfg.PreferredEditor
	if currentEditor == "" {
		currentEditor = "auto"
	}

	fmt.Printf("Current editor: %s\n\n", common.GetEditorDisplayName(currentEditor))

	// Show selection prompt
	selectPrompt := promptui.Select{
		Label: "Select your preferred text editor",
		Items: items,
		Size:  10,
	}

	index, _, err := selectPrompt.Run()
	if err != nil {
		return err
	}

	selectedEditor := editorCommands[index]

	// Save to config
	cfg.PreferredEditor = selectedEditor
	if err := config.SaveLocalConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("\n✅ Editor setup complete!\n")
	fmt.Printf("Selected editor: %s\n", common.GetEditorDisplayName(selectedEditor))

	if selectedEditor == "auto" {
		// Show what will be used
		detectedCmd, err := common.GetEditorCommand(selectedEditor)
		if err == nil {
			fmt.Printf("Will use: %s\n", detectedCmd)
		}
	}

	fmt.Println("\nYour editor will be used for multiline input when submitting review requests.")

	return nil
}
