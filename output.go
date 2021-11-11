package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"

	"github.com/gotd/td/session"
)

type outputFlag struct {
	Value string
	w     io.Writer
	close func() error
}

func (o outputFlag) String() string {
	return o.Value
}

func (o *outputFlag) Set(s string) error {
	o.Value = s

	o.w = os.Stdout
	if s != "" {
		file, err := os.Create(s)
		if err != nil {
			return err
		}
		o.w = file
		o.close = file.Close
	}
	return nil
}

func (o *outputFlag) Write(data []byte) (int, error) {
	if o.w == nil {
		o.w = os.Stdout
	}

	return o.w.Write(data)
}

func (o *outputFlag) Close() error {
	if o.close == nil {
		return nil
	}
	return o.close()
}

type printOptions struct {
	Pretty bool
	Output outputFlag
}

func (p *printOptions) install(set *flag.FlagSet) {
	set.BoolVar(&p.Pretty, "pretty", false, "pretty json")
	set.Var(&p.Output, "output", "output (default: writes to stdout)")
}

func printSession(data *session.Data, opts printOptions) error {
	e := json.NewEncoder(&opts.Output)
	if opts.Pretty {
		e.SetIndent("", "\t")
	}
	return e.Encode(data)
}
