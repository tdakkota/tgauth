package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"rsc.io/qr"
)

func qrCmd() *cobra.Command {
	var (
		pngPath       string
		password      string
		passwordStdin bool
		gotdFlags     gotdOptions
		printFlags    printOptions
	)

	cmd := &cobra.Command{
		Use:   "qr",
		Short: "Create session via QR login flow",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			if passwordStdin {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return errors.Wrap(err, "read password from stdin")
				}
				password = strings.TrimSpace(string(data))
			}

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
					if err == nil {
						return nil
					}
					if !tgerr.Is(err, "SESSION_PASSWORD_NEEDED") {
						return errors.Wrap(err, "QR login")
					}
					if password == "" {
						return errors.New("2FA password required but not provided; use --password or --password-stdin")
					}
					_, err = client.Auth().Password(ctx, password)
					if err != nil {
						return errors.Wrap(err, "password login")
					}
					return nil
				},
			)
			if err != nil {
				return err
			}

			return printSession(ctx, data, &printFlags)
		},
	}

	cmd.Flags().StringVar(&pngPath, "png", "", "Generate path to image of QR")
	cmd.Flags().StringVar(&password, "password", "", "2FA password")
	cmd.Flags().BoolVar(&passwordStdin, "password-stdin", false, "Read 2FA password from stdin")
	gotdFlags.install(cmd.Flags())
	printFlags.install(cmd.Flags())

	return cmd
}
