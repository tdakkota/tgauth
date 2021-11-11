package main

import (
	"context"
	"flag"

	"github.com/cristalhq/acmd"
	"github.com/gotd/td/constant"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

func botCmd() acmd.Command {
	return acmd.Command{
		Name:        "bot",
		Description: "Create session via bot token authorization",
		Do:          botDo,
		Subcommands: nil,
	}
}

func botDo(ctx context.Context, args []string) (rErr error) {
	s := flag.NewFlagSet("bot", flag.ContinueOnError)
	var (
		token   string
		output  outputFlag
		log     bool
		appID   int
		appHash string
	)
	s.StringVar(&token, "token", "", "Bot token")
	s.Var(&output, "output", "output (default: writes to stdout)")
	s.IntVar(&appID, "app-id", constant.TestAppID, "App id (default: Telegram Desktop test)")
	s.StringVar(&appHash, "app-hash", constant.TestAppHash, "App hash (default: Telegram Desktop test)")
	s.BoolVar(&log, "log", false, "Verbose log")

	if err := s.Parse(args); err != nil {
		return err
	}

	var storage session.StorageMemory
	client := telegram.NewClient(appID, appHash, telegram.Options{
		SessionStorage: &storage,
	})
	if err := client.Run(ctx, func(ctx context.Context) error {
		_, err := client.Auth().Bot(ctx, token)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return storage.Dump(&output)
}
