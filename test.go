package main

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
	"github.com/spf13/cobra"
)

func testCmd() *cobra.Command {
	var (
		phone      string
		gotdFlags  gotdOptions
		printFlags printOptions
	)

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Create test user session",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			gotdFlags.Test = true

			data, err := gotdFlags.GetSession(
				ctx, telegram.Options{},
				func(ctx context.Context, client *telegram.Client) error {
					dc := client.Config().ThisDC

					var err error
					if phone != "" {
						err = client.Auth().TestUser(ctx, phone, dc)
					} else {
						err = client.Auth().Test(ctx, dc)
					}
					if err != nil {
						return errors.Wrap(err, "login")
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

	cmd.Flags().StringVar(&phone, "phone", "", "Phone to acquire")
	gotdFlags.install(cmd.Flags())
	printFlags.install(cmd.Flags())

	return cmd
}
