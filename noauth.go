package main

import (
	"context"
	"flag"

	"github.com/cristalhq/acmd"
	"github.com/gotd/td/telegram"
)

func noauthCmd() acmd.Command {
	return acmd.Command{
		Name:        "noauth",
		Description: "Create session without authorization",
		Do:          noauthDo,
		Subcommands: nil,
	}
}

func noauthDo(ctx context.Context, args []string) (rErr error) {
	s := flag.NewFlagSet("noauth", flag.ContinueOnError)
	var (
		gotdFlags  gotdOptions
		printFlags printOptions
	)
	gotdFlags.install(s)
	printFlags.install(s)

	if err := s.Parse(args); err != nil {
		return err
	}

	data, err := gotdFlags.GetSession(
		ctx, telegram.Options{},
		func(ctx context.Context, client *telegram.Client) error {
			return nil
		},
	)
	if err != nil {
		return err
	}

	return printSession(data, printFlags)
}
