package main

import (
	"context"
	"flag"

	"github.com/cristalhq/acmd"
	"github.com/go-faster/errors"
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
		token      string
		gotdFlags  gotdOptions
		printFlags printOptions
	)
	s.StringVar(&token, "token", "", "Bot token")
	gotdFlags.install(s)
	printFlags.install(s)

	if err := s.Parse(args); err != nil {
		return err
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

	return printSession(data, printFlags)
}
