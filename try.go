package main

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/spf13/cobra"
)

func tryCmd() *cobra.Command {
	var (
		sessionFile string
		gotdFlags   gotdOptions
		printFlags  printOptions
	)

	cmd := &cobra.Command{
		Use:   "try",
		Short: "Print user info",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			var data []byte
			if sessionFile == "" {
				d, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				data = d
			} else {
				d, err := os.ReadFile(filepath.Clean(sessionFile))
				if err != nil {
					return err
				}
				data = d
			}

			storage := &session.StorageMemory{}
			if err := storage.StoreSession(ctx, data); err != nil {
				return errors.Wrap(err, "invalid session")
			}

			client, err := gotdFlags.Client(telegram.Options{
				SessionStorage: storage,
			})
			if err != nil {
				return err
			}

			return client.Run(ctx, func(ctx context.Context) error {
				self, err := client.Self(ctx)
				if err != nil {
					return err
				}
				return printFlags.printData(self)
			})
		},
	}

	cmd.Flags().StringVar(&sessionFile, "session", "", "Path to session file (default: reads from stdin)")
	gotdFlags.install(cmd.Flags())
	printFlags.install(cmd.Flags())

	return cmd
}
