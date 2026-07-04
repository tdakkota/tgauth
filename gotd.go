package main

import (
	"context"
	"flag"

	"github.com/go-faster/errors"
	"github.com/gotd/td/constant"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type gotdOptions struct {
	AppID     int
	AppHash   string
	DC        int
	Test      bool
	Log       bool
	LogLevel  zapcore.Level
	LogFormat string
}

func (p *gotdOptions) install(set *flag.FlagSet) {
	set.IntVar(&p.AppID, "app-id", constant.TestAppID, "AppID (default: Telegram Desktop test)")
	set.StringVar(&p.AppHash, "app-hash", constant.TestAppHash, "AppHash (default: Telegram Desktop test)")
	set.IntVar(&p.DC, "DC", 2, "DC ID to use")
	set.BoolVar(&p.Test, "test", false, "Use test server")
	set.BoolVar(&p.Log, "log", false, "enable logging")
	set.Var(&p.LogLevel, "loglevel", "logging level")
	set.StringVar(&p.LogFormat, "logformat", "console", "log format (json or console)")
}

func (p gotdOptions) Client(opts telegram.Options) (*telegram.Client, error) {
	opts.DC = p.DC
	if opts.Logger == nil && p.Log {
		var zapCfg zap.Config
		switch p.LogFormat {
		case "console":
			zapCfg = zap.NewDevelopmentConfig()
		case "json":
			zapCfg = zap.NewProductionConfig()
		}

		logger, err := zapCfg.Build(zap.IncreaseLevel(p.LogLevel))
		if err != nil {
			return nil, errors.Wrap(err, "create logger")
		}
		opts.Logger = logger
	}
	if p.Test {
		opts.DCList = dcs.Test()
	}
	return telegram.NewClient(p.AppID, p.AppHash, opts), nil
}

func (p gotdOptions) GetSession(
	ctx context.Context,
	opts telegram.Options,
	cb func(ctx context.Context, client *telegram.Client) error,
) (*session.Data, error) {
	var storage session.StorageMemory
	opts.SessionStorage = &storage

	client, err := p.Client(opts)
	if err != nil {
		return nil, errors.Wrap(err, "initialize")
	}

	if err := client.Run(ctx, func(ctx context.Context) error {
		return cb(ctx, client)
	}); err != nil {
		return nil, err
	}

	data, err := (&session.Loader{
		Storage: &storage,
	}).Load(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "gotd generated invalid session"))
	}

	return data, nil
}
