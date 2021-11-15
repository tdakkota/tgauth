package main

import (
	"context"
	"flag"

	"github.com/cristalhq/acmd"
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
		gotdFlags  gotdOptions
		printFlags printOptions
	)
	gotdFlags.install(s)
	printFlags.install(s)

	if err := s.Parse(args); err != nil {
		return err
	}

	client, err := gotdFlags.Client(telegram.Options{})
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
