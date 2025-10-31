package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	*baseCmd
}

var version = "x-dev"

func NewRootCommand() *RootCmd {
	return &RootCmd{
		baseCmd: newBaseCommand(&cobra.Command{
			Use:     "cool",
			Version: version,
			Short:   "bfi lora tools.",
			Long:    `cool is tools for lora devs.`,
		}),
	}
}
