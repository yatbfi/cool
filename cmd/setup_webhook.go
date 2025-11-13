package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
)

type SetupWebhookCmd struct {
	*baseCmd
}

func NewSetupWebhookCmd() *SetupWebhookCmd {
	cmd := &SetupWebhookCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "webhook",
		Short: "Setup GChat webhook URLs",
		RunE:  cmd.run,
	})
	return cmd
}

func (c *SetupWebhookCmd) run(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	// Ask for GChat Review Webhook URL
	reviewWebhookPrompt := promptui.Prompt{
		Label:   "GChat Review Webhook URL",
		Default: cfg.GChatReviewWebhookURL,
	}
	reviewWebhookURL, err := reviewWebhookPrompt.Run()
	if err != nil {
		return err
	}

	// Ask for GChat Collab Webhook URL
	collabWebhookPrompt := promptui.Prompt{
		Label:   "GChat Collab Webhook URL",
		Default: cfg.GChatCollabWebhookURL,
	}
	collabWebhookURL, err := collabWebhookPrompt.Run()
	if err != nil {
		return err
	}

	cfg.GChatReviewWebhookURL = reviewWebhookURL
	cfg.GChatCollabWebhookURL = collabWebhookURL

	if err = config.SaveLocalConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("\n✅ Webhook setup complete!\n")
	fmt.Printf("GChat Review Webhook URL : %s\n", reviewWebhookURL)
	fmt.Printf("GChat Collab Webhook URL : %s\n", collabWebhookURL)
	return nil
}
