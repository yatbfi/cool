package main

import (
	"fmt"
	"os"

	"github.com/yatbfi/cool/cmd"
	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/domain/usecases"
)

var (
	rootCmd *cmd.RootCmd
)

func main() {
	cfg := config.GetConfig()

	initRootCmd()
	initCommands(cfg)

	if err := rootCmd.Cmd().Execute(); err != nil {
		fmt.Println("")
		os.Exit(1)
	}
}

func initRootCmd() {
	rootCmd = cmd.NewRootCommand()
	_ = cmd.NewCompletionCmd().Reg(rootCmd)
}

func initCommands(cfg *config.Config) {
	// --- Usecases ---
	gChatUc := usecases.NewGChatUsecase(cfg)
	reviewUc := usecases.NewReviewUsecase(cfg, gChatUc)

	// --- Review command tree ---
	reviewCmd := cmd.NewReviewCmd().Reg(rootCmd)
	cmd.NewReviewRequestCmd(reviewUc).Reg(reviewCmd)

	// --- Setup command (register directly under root) ---
	cmd.NewSetupCmd().Reg(rootCmd)
}
