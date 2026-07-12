package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/spf13/cobra"
)

func getDefaultTDataPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "Telegram Desktop", "tdata")
	}
	return ""
}

var (
	errWrongPasscode = errors.New("wrong passcode")
	errTDataRequired = errors.New("argument tdata is required")
)

func tdesktopCmd() *cobra.Command {
	var (
		tdata      string
		passcode   string
		idx        int
		printFlags printOptions
	)

	cmd := &cobra.Command{
		Use:   "tdesktop",
		Short: "Create session using Telegram Desktop storage",
		RunE: func(cmd *cobra.Command, _ []string) (rErr error) {
			ctx := cmd.Context()

			if tdata == "" {
				return errTDataRequired
			}
			dir, err := os.Stat(tdata)
			switch {
			case errors.Is(err, fs.ErrNotExist):
				return errors.Errorf("can't find tdata (path: %q)", tdata)
			case err != nil:
				return err
			case !dir.IsDir():
				return errors.Errorf("%q is not a directory", tdata)
			}

			accounts, err := tdesktop.Read(tdata, []byte(passcode))
			switch {
			case errors.Is(err, tdesktop.ErrKeyInfoDecrypt):
				return errWrongPasscode
			case err != nil:
				return err
			}

			switch {
			case idx >= len(accounts):
				return errors.Errorf("too big index %d, there are only %d account(s)", idx, len(accounts))
			case idx < 0 && len(accounts) > 1:
				// TODO(tdakkota): choose by username
				options := make([]string, len(accounts))
				for i, a := range accounts {
					auth := a.Authorization
					options[i] = fmt.Sprintf("User %d (test: %t)", auth.UserID, a.Config.Environment.Test())
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
			case idx < 0:
				idx = 0
			}

			data, err := session.TDesktopSession(accounts[idx])
			if err != nil {
				return errors.Wrap(err, "convert")
			}

			return printSession(ctx, data, &printFlags)
		},
	}

	cmd.Flags().StringVar(&tdata, "tdata", getDefaultTDataPath(), "path to tdata")
	cmd.Flags().StringVar(&passcode, "passcode", "", "passcode")
	cmd.Flags().IntVar(&idx, "idx", -1, "account index")
	printFlags.install(cmd.Flags())

	return cmd
}
