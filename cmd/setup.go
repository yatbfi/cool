package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
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
		Short: "Initial setup (save your name and email)",
		RunE:  cmd.run,
	})
	return cmd
}

func (c *SetupCmd) run(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	// Ask for name
	namePrompt := promptui.Prompt{
		Label:   "What is your name?",
		Default: cfg.UserName,
	}
	name, err := namePrompt.Run()
	if err != nil {
		return err
	}

	// Ask for email
	emailPrompt := promptui.Prompt{
		Label:   "What is your email?",
		Default: cfg.UserEmail,
	}
	email, err := emailPrompt.Run()
	if err != nil {
		return err
	}

	cfg.UserName = name
	cfg.UserEmail = email

	if err = config.SaveLocalConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("\n✅ Setup complete!\n")
	fmt.Printf("Name : %s\n", name)
	fmt.Printf("Email: %s\n", email)
	fmt.Println("Configuration saved successfully at ~/.cool-cli/config.json")
	return nil
}
