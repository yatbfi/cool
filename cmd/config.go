package cmd

import (
	"github.com/spf13/cobra"
)

// ConfigCmd handles config operations
type ConfigCmd struct {
	*baseCmd
}

// NewConfigCmd creates a new config command
func NewConfigCmd() *ConfigCmd {
	cmd := &ConfigCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long:  `View and manage Cool CLI configuration.`,
	})
	return cmd
}
