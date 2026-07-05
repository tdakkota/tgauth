# tgauth

[![pkg-img]][pkg-url]
[![coverage-img]][coverage-url]

Simple CLI tool for creating gotd sessions.

## Install

```
go install github.com/tdakkota/tgauth@latest
```

## Examples

### Create user session interactively
```shell
$ tgauth user
? Your phone number 79000427572
? Activation code [? for help] ?

? The code sent by Telegram SMS
? Activation code
```

### Read Telegram Desktop session
```shell
$ tgauth -tdata "Telegram Desktop/tdata"
```

### Read Telegram Desktop session with passcode
```shell
$ tgauth -passcode 12345 -tdata "Telegram Desktop/tdata"
```

### Create session using bot token
```shell
$ tgauth bot -token 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew1
```

### QR login
```shell
# Will print a QR using ANSI in your terminal.
$ tgauth qr
```

[pkg-img]: https://pkg.go.dev/badge/tdakkota/tgauth
[pkg-url]: https://pkg.go.dev/github.com/tdakkota/tgauth
[coverage-img]: https://codecov.io/gh/tdakkota/tgauth/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/tdakkota/tgauth
