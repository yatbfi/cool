package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

type UpdateCmd struct {
	*baseCmd
}

func NewUpdateCmd() *UpdateCmd {
	cmd := &UpdateCmd{}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "update",
		Short: "Update cool CLI to the latest version",
		RunE:  cmd.run,
	})
	return cmd
}

func (c *UpdateCmd) run(_ *cobra.Command, _ []string) error {
	fmt.Println("🚀 Updating cool CLI to the latest version...")

	goBin, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("Go binary not found in PATH: %w", err)
	}

	updateCmd := exec.Command(goBin, "install", "github.com/yatbfi/cool@latest")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	updateCmd.Env = os.Environ()

	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update cool: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ cool has been successfully updated to the latest version!")

	// Print helpful info about PATH
	fmt.Printf("💡 Binary location: %s\n", os.ExpandEnv("$GOBIN"))
	if os.Getenv("GOBIN") == "" {
		fmt.Printf("💡 Default GOPATH bin: %s/bin\n", os.ExpandEnv("$GOPATH"))
	}

	fmt.Printf("💡 OS: %s | Arch: %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("🎉 Run 'cool version' to verify your update.")
	return nil
}
