package cmd

import (
	"fmt"

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
		// --- Common environment checks ---
		if err := common.OsArchSupported(); err != nil {
			return err
		}
		if _, err := common.GetGobin(); err != nil {
			return err
		}

		// --- Skip setup validation for setup command itself ---
		if cmd.Name() == "setup" {
			return nil
		}

		// --- User setup check ---
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
	return c.cmd
}

func (c *baseCmd) Reg(p Command) Command {
	p.Cmd().AddCommand(c.Cmd())
	return c
}
