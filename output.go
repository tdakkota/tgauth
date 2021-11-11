package main

import (
	"io"
	"os"
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
