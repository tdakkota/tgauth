package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cristalhq/acmd"
	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
)

func tdesktopCmd() acmd.Command {
	return acmd.Command{
		Name:        "tdesktop",
		Description: "Create session using Telegram Desktop storage",
		Do:          tdesktopDo,
	}
}

func getDefaultTDataPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Downloads", "Telegram", "tdata")
}

var (
	errWrongPasscode = errors.New("wrong passcode")
)

func tdesktopDo(ctx context.Context, args []string) (rErr error) {
	s := flag.NewFlagSet("tdesktop", flag.ContinueOnError)
	var (
		tdata    string
		output   outputFlag
		passcode string
		idx      int
	)
	s.StringVar(&tdata, "tdata", getDefaultTDataPath(), "path to tdata")
	s.Var(&output, "output", "output (default: writes to stdout)")
	s.StringVar(&passcode, "passcode", "", "passcode")
	s.IntVar(&idx, "idx", -1, "account index")

	if err := s.Parse(args); err != nil {
		return err
	}

	accounts, err := tdesktop.Read(tdata, []byte(passcode))
	switch {
	case errors.Is(err, tdesktop.ErrKeyInfoDecrypt):
		return errWrongPasscode
	case err != nil:
		return err
	}

	if idx >= len(accounts) {
		return errors.Errorf("too big index %d, there are only %d account(s)", idx, len(accounts))
	}
	if idx < 0 && len(accounts) > 1 {
		// TODO(tdakkota): choose by username
		options := make([]string, len(accounts))
		for i, a := range accounts {
			options[i] = strconv.FormatUint(a.Authorization.UserID, 10)
		}

		sel := &survey.Select{
			Message: "Choose account",
			Options: options,
		}
		var answer string
		if err := survey.AskOne(sel, &answer); err != nil {
			return err
		}

		for i, option := range options {
			if option == answer {
				idx = i
				break
			}
		}
	}
	data, err := session.TDesktopSession(accounts[idx])
	if err != nil {
		return errors.Wrap(err, "convert")
	}

	return json.NewEncoder(&output).Encode(data)
}
