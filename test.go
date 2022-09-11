package main

import (
	"context"
	"flag"

	"github.com/cristalhq/acmd"
	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
)

func testCmd() acmd.Command {
	return acmd.Command{
		Name:        "test",
		Description: "Create test user session",
		ExecFunc:    testExec,
	}
}

func testExec(ctx context.Context, args []string) (rErr error) {
	s := flag.NewFlagSet("test", flag.ContinueOnError)
	var (
		phone      string
		gotdFlags  gotdOptions
		printFlags printOptions
	)
	s.StringVar(&phone, "phone", "", "Phone to acquire")
	gotdFlags.install(s)
	printFlags.install(s)

	if err := s.Parse(args); err != nil {
		return err
	}

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
}
