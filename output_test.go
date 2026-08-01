package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrintData(t *testing.T) {
	data := sessionFile{Version: 1, Data: *testSessionData()}

	tests := []struct {
		name    string
		opts    printOptions
		want    string
		wantErr bool
	}{
		{
			name: "Template",
			opts: printOptions{Template: "{{ .Data.DC }}"},
			want: "4",
		},
		{
			// Template takes precedence over format.
			name: "TemplateOverridesFormat",
			opts: printOptions{Template: "dc={{ .Data.DC }}", Format: "pp"},
			want: "dc=4",
		},
		{
			name: "JSON",
			opts: printOptions{Format: "json"},
			want: `"Version":1`,
		},
		{
			name:    "BadTemplate",
			opts:    printOptions{Template: "{{ .Data.DC "},
			wantErr: true,
		},
		{
			name:    "UnknownFormat",
			opts:    printOptions{Format: "yaml"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			opts := tt.opts
			opts.Output.w = &buf

			err := opts.printData(data)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Contains(t, buf.String(), tt.want)
		})
	}
}
