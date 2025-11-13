package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
)

type SetupEmailCmd struct {
	*baseCmd
}

func NewSetupEmailCmd() *SetupEmailCmd {
	cmd := &SetupEmailCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "email",
		Short: "Setup user name and email",
		RunE:  cmd.run,
	})
	return cmd
}

func (c *SetupEmailCmd) run(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("\n✅ Email setup complete!\n")
	fmt.Printf("Name  : %s\n", name)
	fmt.Printf("Email : %s\n", email)
	return nil
}
