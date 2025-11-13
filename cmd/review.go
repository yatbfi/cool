package cmd

import (
	"github.com/spf13/cobra"
)

// ReviewCmd is the parent command for review operations
type ReviewCmd struct {
	*baseCmd
}

// NewReviewCmd creates a new review command
func NewReviewCmd() *ReviewCmd {
	cmd := &ReviewCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "review",
		Short: "Manage code review requests",
		Long: `Manage code review requests to tech lead and architect.

This command provides subcommands to:
- Submit review requests to tech lead
- View review history
- Submit approved reviews to collaboration channel`,
	})
	return cmd
}
