package main

import (
	"io"
	"os"
)

type OutputFlag struct {
	Value string
	w     io.Writer
	close func() error
}

func (o OutputFlag) String() string {
	return o.Value
}

func (o *OutputFlag) Set(s string) error {
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

func (o *OutputFlag) Write(data []byte) (int, error) {
	if o.w == nil {
		o.w = os.Stdout
	}

	return o.w.Write(data)
}

func (o *OutputFlag) Close() error {
	if o.close == nil {
		return nil
	}
	return o.close()
}
