package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/pkg/common"
)

// Command defines interface for subcommands
type Command interface {
	Cmd() *cobra.Command
}

// baseCmd is the shared struct for cobra commands
type baseCmd struct {
	cmd *cobra.Command
}

func newBaseCommand(cmd *cobra.Command) *baseCmd {
	return &baseCmd{cmd: cmd}
}

func (c *baseCmd) Cmd() *cobra.Command {
	c.cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := validateEnvironment(); err != nil {
			return err
		}

		if shouldSkipValidation(cmd.Name()) {
			return nil
		}

		if err := validateUserSetup(); err != nil {
			return err
		}

		if err := validateCommandSpecificConfig(cmd.Name()); err != nil {
			return err
		}

		return nil
	}
	return c.cmd
}

// validateEnvironment checks if the OS and required binaries are available
func validateEnvironment() error {
	if err := common.OsArchSupported(); err != nil {
		return err
	}
	if _, err := common.GetGobin(); err != nil {
		return err
	}
	return nil
}

// shouldSkipValidation returns true for commands that don't require setup validation
func shouldSkipValidation(cmdName string) bool {
	return slices.Index([]string{"setup", "update"}, cmdName) >= 0
}

// validateUserSetup checks if user name and email are configured
func validateUserSetup() error {
	cfg := config.GetConfig()
	if cfg.UserName == "" || cfg.UserEmail == "" {
		fmt.Println("⚠️  User setup is incomplete.")
		fmt.Println("Please run the setup command first:")
		fmt.Println("   \"cool setup\"")
		fmt.Println()
		return fmt.Errorf("missing user setup (name/email)")
	}
	return nil
}

// validateCommandSpecificConfig validates configuration required for specific commands
func validateCommandSpecificConfig(cmdName string) error {
	cfg := config.GetConfig()

	switch cmdName {
	case "request":
		// review request command needs review webhook
		if cfg.GChatReviewWebhookURL == "" {
			fmt.Println("⚠️  GChat review webhook URL is not configured.")
			fmt.Println("Please run the setup command to configure it:")
			fmt.Println("   cool setup webhook")
			fmt.Println()
			return fmt.Errorf("missing GChat review webhook URL")
		}
	case "submit-collab":
		// review submit-collab command needs collab webhook
		if cfg.GChatCollabWebhookURL == "" {
			fmt.Println("⚠️  GChat collaboration webhook URL is not configured.")
			fmt.Println("Please run the setup command to configure it:")
			fmt.Println("   cool setup webhook")
			fmt.Println()
			return fmt.Errorf("missing GChat collaboration webhook URL")
		}
	}

	return nil
}

func (c *baseCmd) Reg(p Command) Command {
	p.Cmd().AddCommand(c.Cmd())
	return c
}
