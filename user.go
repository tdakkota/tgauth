package main

import (
	"context"
	"flag"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cristalhq/acmd"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

func userCmd() acmd.Command {
	return acmd.Command{
		Name:        "user",
		Description: "Create session via plain user authorization",
		ExecFunc:    userExec,
	}
}

type surveyAuth struct {
	phone    string
	password string
}

func (s *surveyAuth) install(set *flag.FlagSet) {
	set.StringVar(&s.phone, "phone", "", "Phone number")
	set.StringVar(&s.password, "password", "", "Password")
}

func (s surveyAuth) askOneString(original, msg, help string, length int) (code string, err error) {
	if original != "" {
		return original, nil
	}

	options := []survey.AskOpt{
		survey.WithValidator(survey.Required),
	}
	if length > 0 {
		options = append(options,
			survey.WithValidator(survey.MinLength(length)),
			survey.WithValidator(survey.MaxLength(length)),
		)
	}

	err = survey.AskOne(&survey.Input{
		Message: msg,
		Help:    help,
	}, &code, options...)
	return code, err
}

func (s surveyAuth) Phone(ctx context.Context) (string, error) {
	return s.askOneString(s.phone, "Your phone number", "", 0)
}

func (s surveyAuth) Password(ctx context.Context) (string, error) {
	return s.askOneString(s.password, "Your password", "", 0)
}

func (s surveyAuth) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	// TODO(tdakkota): support signup?
	return &auth.SignUpRequired{TermsOfService: tos}
}

func (s surveyAuth) SignUp(ctx context.Context) (auth.UserInfo, error) {
	firstName, err := s.askOneString("", "First Name", "", 0)
	if err != nil {
		return auth.UserInfo{}, err
	}
	lastName, err := s.askOneString("", "Last Name", "", 0)
	if err != nil {
		return auth.UserInfo{}, err
	}
	return auth.UserInfo{
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

func (s surveyAuth) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	var (
		length int
		via    string
	)

	switch codeType := sentCode.Type.(type) {
	case *tg.AuthSentCodeTypeApp:
		length = codeType.Length
		via = "via app, check your private messages"
	case *tg.AuthSentCodeTypeSMS:
		length = codeType.Length
		via = "SMS"
	case *tg.AuthSentCodeTypeCall:
		length = codeType.Length
		via = "a phone call"
	}

	return s.askOneString("", "Activation code", "The code sent by Telegram "+via, length)
}

func userExec(ctx context.Context, args []string) (rErr error) {
	s := flag.NewFlagSet("user", flag.ContinueOnError)
	var (
		ua         surveyAuth
		gotdFlags  gotdOptions
		printFlags printOptions
	)
	ua.install(s)
	gotdFlags.install(s)
	printFlags.install(s)

	if err := s.Parse(args); err != nil {
		return err
	}

	data, err := gotdFlags.GetSession(
		ctx, telegram.Options{},
		func(ctx context.Context, client *telegram.Client) error {
			return auth.NewFlow(ua, auth.SendCodeOptions{}).Run(ctx, client.Auth())
		},
	)
	if err != nil {
		return err
	}

	return printSession(ctx, data, printFlags)
}
