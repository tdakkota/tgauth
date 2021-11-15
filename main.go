package main

import (
	"fmt"
	"os"

	"github.com/cristalhq/acmd"
)

const description = `Simple CLI tool for creating gotd sessions.`

func main() {
	cmds := []acmd.Command{
		botCmd(),
		userCmd(),
		qrCmd(),
		tdesktopCmd(),
		noauthCmd(),
	}
	cfg := acmd.Config{
		AppDescription: description,
	}
	if err := acmd.RunnerOf(cmds, cfg).Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}
