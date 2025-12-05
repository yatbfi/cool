package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/internal/domain/usecase"
	infraRepo "github.com/yatbfi/cool/internal/infrastructure/repository"
)

// RootCmd is the root command
type RootCmd struct {
	*baseCmd
}

// NewRootCommand creates the root command with all subcommands
func NewRootCommand() *RootCmd {
	rootCmd := &RootCmd{}
	rootCmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "cool",
		Short: "Cool CLI - Developer tools",
	})

	// Initialize repositories and usecase
	historyRepo, err := infraRepo.NewReviewHistoryRepository()
	if err != nil {
		// Handle error gracefully - command will still work but review features may fail
		historyRepo = nil
	}

	gchatUc := usecase.NewGChatUsecase()
	reviewUc := usecase.NewReviewUsecase(historyRepo, gchatUc)

	// Add subcommands
	setupCmd := NewSetupCmd()
	setupCmd.Cmd().AddCommand(
		NewSetupEmailCmd().Cmd(),
		NewSetupWebhookCmd().Cmd(),
		NewSetupEditorCmd().Cmd(),
		NewSetupProjectRootCmd().Cmd(),
	)

	reviewCmd := NewReviewCmd()
	reviewCmd.Cmd().AddCommand(
		NewReviewRequestCmd(reviewUc).Cmd(),
		NewReviewHistoriesCmd(reviewUc).Cmd(),
		NewReviewSubmitCollabCmd(reviewUc).Cmd(),
	)

	configCmd := NewConfigCmd()
	configCmd.Cmd().AddCommand(
		NewConfigPreviewCmd().Cmd(),
	)

	rootCmd.cmd.AddCommand(
		setupCmd.Cmd(),
		reviewCmd.Cmd(),
		configCmd.Cmd(),
		NewRunCmd().Cmd(),
		NewUpdateCmd().Cmd(),
		NewCompletionCmd().Cmd(),
	)

	return rootCmd
}

// Execute executes the root command
func Execute() error {
	rootCmd := NewRootCommand()
	return rootCmd.cmd.Execute()
}
