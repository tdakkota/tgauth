package main

import (
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/cristalhq/acmd"
	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

func tryCmd() acmd.Command {
	return acmd.Command{
		Name:        "try",
		Description: "Print user info",
		Do:          tryDo,
		Subcommands: nil,
	}
}

func tryDo(ctx context.Context, args []string) error {
	s := flag.NewFlagSet("try", flag.ContinueOnError)
	var (
		sessionFile string
		gotdFlags   gotdOptions
		printFlags  printOptions
	)
	s.StringVar(&sessionFile, "session", "", "Path to session file (default: reads from stdin)")
	gotdFlags.install(s)
	printFlags.install(s)

	if err := s.Parse(args); err != nil {
		return err
	}

	var (
		data []byte
	)
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
}
