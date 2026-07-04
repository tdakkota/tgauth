package main

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
	"github.com/spf13/cobra"
)

func botCmd() *cobra.Command {
	var (
		token      string
		gotdFlags  gotdOptions
		printFlags printOptions
	)

	cmd := &cobra.Command{
		Use:   "bot",
		Short: "Create session via bot token authorization",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			if token == "" {
				return errors.New("token option is required")
			}

			data, err := gotdFlags.GetSession(
				ctx, telegram.Options{},
				func(ctx context.Context, client *telegram.Client) error {
					_, err := client.Auth().Bot(ctx, token)
					if err != nil {
						return errors.Wrap(err, "bot login")
					}
					return nil
				},
			)
			if err != nil {
				return err
			}

			return printSession(ctx, data, printFlags)
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "Bot token")
	gotdFlags.install(cmd.Flags())
	printFlags.install(cmd.Flags())

	return cmd
}
