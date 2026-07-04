package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/tg"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"rsc.io/qr"
)

func qrCmd() *cobra.Command {
	var (
		pngPath    string
		gotdFlags  gotdOptions
		printFlags printOptions
	)

	cmd := &cobra.Command{
		Use:   "qr",
		Short: "Create session via QR login flow",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			d := tg.NewUpdateDispatcher()
			loggedIn := qrlogin.OnLoginToken(d)
			data, err := gotdFlags.GetSession(
				ctx, telegram.Options{
					UpdateHandler: d,
				},
				func(ctx context.Context, client *telegram.Client) (rErr error) {
					show := func(_ context.Context, token qrlogin.Token) error {
						qrterminal.Generate(token.URL(), qrterminal.L, os.Stdout)
						return nil
					}
					if pngPath != "" {
						f, err := os.Create(filepath.Clean(pngPath))
						if err != nil {
							return err
						}
						defer multierr.AppendInvoke(&rErr, multierr.Close(f))
						show = func(_ context.Context, token qrlogin.Token) error {
							encoded, err := qr.Encode(token.URL(), qrterminal.L)
							if err != nil {
								return errors.Wrap(err, "encode QR")
							}
							if _, err := f.Write(encoded.PNG()); err != nil {
								return errors.Wrap(err, "write png")
							}
							return nil
						}
					}

					_, err := client.QR().Auth(ctx, loggedIn, show)
					if err != nil {
						return errors.Wrap(err, "QR login")
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

	cmd.Flags().StringVar(&pngPath, "png", "", "Generate path to image of QR")
	gotdFlags.install(cmd.Flags())
	printFlags.install(cmd.Flags())

	return cmd
}
