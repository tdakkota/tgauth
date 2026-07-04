// Command tgauth is a simple CLI tool to create Telegram/gotd sessions.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "tgauth",
		Short:         "Simple CLI tool for creating gotd sessions.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(
		botCmd(),
		userCmd(),
		testCmd(),
		qrCmd(),
		tdesktopCmd(),
		noauthCmd(),
		tryCmd(),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := root.ExecuteContext(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
