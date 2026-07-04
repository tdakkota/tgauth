package main

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/spf13/cobra"
)

func noauthCmd() *cobra.Command {
	var (
		gotdFlags  gotdOptions
		printFlags printOptions
	)

	cmd := &cobra.Command{
		Use:   "noauth",
		Short: "Create session without authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			data, err := gotdFlags.GetSession(
				ctx, telegram.Options{},
				func(ctx context.Context, client *telegram.Client) error {
					return nil
				},
			)
			if err != nil {
				return err
			}

			return printSession(ctx, data, printFlags)
		},
	}

	gotdFlags.install(cmd.Flags())
	printFlags.install(cmd.Flags())

	return cmd
}
