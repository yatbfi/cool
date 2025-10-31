package cmd

import (
	"github.com/spf13/cobra"
)

type CompletionCmd struct {
	*baseCmd
}

func NewCompletionCmd() *CompletionCmd {
	return &CompletionCmd{
		baseCmd: newBaseCommand(&cobra.Command{
			Use:    "completion",
			Hidden: true,
		}),
	}
}
