package cmd

import (
	"github.com/spf13/cobra"
)

type ReviewCmd struct {
	*baseCmd
}

func NewReviewCmd() *ReviewCmd {
	cmd := &ReviewCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "review",
		Short: "Review helper command",
	})
	return cmd
}
