# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`tgauth` is a CLI tool (single `main` package, all files at repo root) for creating
Telegram/[gotd](https://github.com/gotd/td) sessions through various auth flows, and printing
the resulting session data. It's a thin wrapper around `gotd/td`'s auth APIs plus a CLI layer.

## Commands

```sh
go build ./...                 # build
go test ./...                  # run all tests
go run . <subcommand> -h       # run a subcommand, e.g. `go run . user -h`
golangci-lint run              # lint (config in .golangci.yml)
golangci-lint fmt ./...        # format (per user's global Go standards)
```

There is no single-package test suite beyond the standard `go test ./...` — this is a flat,
single-package repo, so all tests run together.

## Architecture

Each subcommand is its own file at the repo root and registers an `acmd.Command` in `main.go`:

- `bot.go` — `bot`: auth via bot token
- `user.go` — `user`: interactive phone/code/password auth (uses `AlecAivazis/survey` for prompts,
  implements `gotd/td`'s `auth.UserAuthenticator` via `surveyAuth`)
- `qr.go` — `qr`: QR-code login flow, can render to terminal (ANSI) or write a PNG
- `tdesktop.go` — `tdesktop`: reads sessions directly out of a Telegram Desktop `tdata` directory
  (optionally passcode-protected), letting you pick an account if there are multiple
- `noauth.go` — `noauth`: creates a session without performing authorization
- `test.go` — `test`: creates a session against Telegram's test DC
- `try.go` — `try`: takes an existing session (file or stdin) and calls `Self()` to verify/print
  the logged-in user — useful for validating a session produced by the other commands

Two shared helpers used by nearly every subcommand:

- `gotd.go` (`gotdOptions`) — common gotd client flags (`-app-id`, `-app-hash`, `-DC`, `-test`,
  logging flags) and `GetSession`, which runs a `telegram.Client` against an in-memory session
  storage, executes a caller-supplied auth callback, then loads and returns the resulting
  `session.Data`.
- `output.go` (`printOptions`) — common output flags (`-format` json/pp, `-pretty`, `-template`,
  `-output`) and `printData`/`printSession` to render whatever the auth flow produced.

The pattern for adding a new subcommand: define `xxxCmd() acmd.Command`, parse its own
`flag.FlagSet` embedding `gotdOptions`/`printOptions` as needed, call `gotdFlags.GetSession(...)`
with an auth callback, then `printSession(...)`, and register it in `main.go`.

## CI

`.github/workflows/x.yml` delegates to reusable workflows from `go-faster/x`
(test/cover/lint/commit/codeql) rather than defining steps inline — check that repo if a CI
workflow needs to be understood in depth.
