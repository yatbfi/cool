package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/pkg/common"
)

type ConfigPreviewCmd struct {
	*baseCmd
}

func NewConfigPreviewCmd() *ConfigPreviewCmd {
	cmd := &ConfigPreviewCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "preview",
		Short: "Preview current configuration",
		Long:  `Display all current configuration settings stored in ~/.cool-cli/config.json`,
		RunE:  cmd.run,
	})
	return cmd
}

func (c *ConfigPreviewCmd) run(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	fmt.Println()
	fmt.Println("📋 Current Configuration")
	fmt.Println("========================")
	fmt.Println()

	// User Information
	fmt.Println("👤 User Information:")
	if cfg.UserName != "" {
		fmt.Printf("   Name  : %s\n", cfg.UserName)
	} else {
		fmt.Println("   Name  : (not set)")
	}

	if cfg.UserEmail != "" {
		fmt.Printf("   Email : %s\n", cfg.UserEmail)
	} else {
		fmt.Println("   Email : (not set)")
	}
	fmt.Println()

	// Webhooks
	fmt.Println("🔗 Google Chat Webhooks:")
	if cfg.GChatReviewWebhookURL != "" {
		fmt.Printf("   Review Webhook : %s\n", cfg.GChatReviewWebhookURL)
	} else {
		fmt.Println("   Review Webhook : (not set)")
	}

	if cfg.GChatCollabWebhookURL != "" {
		fmt.Printf("   Collab Webhook : %s\n", cfg.GChatCollabWebhookURL)
	} else {
		fmt.Println("   Collab Webhook : (not set)")
	}
	fmt.Println()

	// Editor
	fmt.Println("✏️  Editor Settings:")
	if cfg.PreferredEditor != "" {
		editorName := common.GetEditorDisplayName(cfg.PreferredEditor)
		if cfg.PreferredEditor == "auto" {
			// Show what will actually be used
			if detectedCmd, err := common.GetEditorCommand(cfg.PreferredEditor); err == nil {
				fmt.Printf("   Preferred Editor : %s (detected: %s)\n", editorName, detectedCmd)
			} else {
				fmt.Printf("   Preferred Editor : %s\n", editorName)
			}
		} else {
			fmt.Printf("   Preferred Editor : %s\n", editorName)
		}
	} else {
		fmt.Println("   Preferred Editor : (not set)")
	}
	fmt.Println()

	// Project Root
	fmt.Println("📁 Project Settings:")
	if cfg.ProjectRoot != "" {
		fmt.Printf("   Project Root : %s\n", cfg.ProjectRoot)
	} else {
		fmt.Println("   Project Root : (not set)")
	}
	fmt.Println()

	// Configuration file location
	fmt.Println("💾 Configuration File:")
	fmt.Println("   ~/.cool-cli/config.json")
	fmt.Println()

	// Show setup commands for missing config
	hasEmpty := cfg.UserName == "" || cfg.UserEmail == "" ||
		cfg.GChatReviewWebhookURL == "" || cfg.GChatCollabWebhookURL == ""

	if hasEmpty {
		fmt.Println("💡 Some settings are not configured. Use these commands to set them up:")
		if cfg.UserName == "" || cfg.UserEmail == "" {
			fmt.Println("   cool setup email")
		}
		if cfg.GChatReviewWebhookURL == "" || cfg.GChatCollabWebhookURL == "" {
			fmt.Println("   cool setup webhook")
		}
		if cfg.PreferredEditor == "" {
			fmt.Println("   cool setup editor")
		}
		if cfg.ProjectRoot == "" {
			fmt.Println("   cool setup project-root")
		}
		fmt.Println()
	}

	return nil
}
