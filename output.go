package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"os"
	"text/template"

	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
	"github.com/k0kubun/pp/v3"
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
	Pretty   bool
	Template string
	Format   string
	Output   outputFlag
}

func (p *printOptions) install(set *flag.FlagSet) {
	set.BoolVar(&p.Pretty, "pretty", false, "Prettify (if format is json)")
	set.StringVar(&p.Template, "template", "", "Go template for formatting")
	set.StringVar(&p.Format, "format", "json", "Printer format (available: json, pp)")
	set.Var(&p.Output, "output", "output (default: writes to stdout)")
}

func (p printOptions) printData(data interface{}) error {
	if tmpl := p.Template; tmpl != "" {
		t, err := template.New("print").Parse(tmpl)
		if err != nil {
			return err
		}
		return t.Execute(&p.Output, t)
	}

	switch p.Format {
	case "pp":
		_, err := pp.Fprintln(&p.Output, data)
		return err
	case "json":
		e := json.NewEncoder(&p.Output)
		if p.Pretty {
			e.SetIndent("", "\t")
		}
		return e.Encode(data)
	default:
		return errors.Errorf("unknown format %q", p.Format)
	}
}

func printSession(ctx context.Context, data *session.Data, opts printOptions) error {
	if ctx.Err() != nil {
		return nil
	}
	return opts.printData(data)
}
