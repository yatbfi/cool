package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
)

type SetupCmd struct {
	*baseCmd
}

func NewSetupCmd() *SetupCmd {
	cmd := &SetupCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "setup",
		Short: "Initial setup - configure user info and webhook URLs",
		Long:  "Run initial setup. Automatically checks if email or webhook is empty and runs the corresponding setup.",
		RunE:  cmd.run,
	})
	return cmd
}

func (c *SetupCmd) run(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	needsEmailSetup := cfg.UserName == "" || cfg.UserEmail == ""
	needsWebhookSetup := cfg.GChatReviewWebhookURL == "" || cfg.GChatCollabWebhookURL == ""

	// If everything is already configured
	if !needsEmailSetup && !needsWebhookSetup {
		fmt.Println("✅ Configuration is already complete!")
		fmt.Printf("\nCurrent settings:\n")
		fmt.Printf("Name                     : %s\n", cfg.UserName)
		fmt.Printf("Email                    : %s\n", cfg.UserEmail)
		fmt.Printf("GChat Review Webhook URL : %s\n", cfg.GChatReviewWebhookURL)
		fmt.Printf("GChat Collab Webhook URL : %s\n", cfg.GChatCollabWebhookURL)
		fmt.Println("\nUse 'cool setup email' or 'cool setup webhook' to update specific settings.")
		return nil
	}

	fmt.Println("🚀 Starting setup...\n")

	// Run email setup if needed
	if needsEmailSetup {
		fmt.Println("📧 Setting up user info...")
		emailCmd := NewSetupEmailCmd()
		if err := emailCmd.run(cmd, args); err != nil {
			return err
		}
		fmt.Println()
	}

	// Run webhook setup if needed
	if needsWebhookSetup {
		fmt.Println("🔗 Setting up webhook URLs...")
		webhookCmd := NewSetupWebhookCmd()
		if err := webhookCmd.run(cmd, args); err != nil {
			return err
		}
		fmt.Println()
	}

	fmt.Println("✅ All setup complete!")
	fmt.Println("Configuration saved successfully at ~/.cool-cli/config.json")
	return nil
}
