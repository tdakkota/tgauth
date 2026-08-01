package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gotd/td/session"
	"github.com/stretchr/testify/require"
)

func testSessionData() *session.Data {
	key := bytes.Repeat([]byte{0xde, 0xad}, 128)
	return &session.Data{
		Config:    session.Config{ThisDC: 4},
		DC:        4,
		AuthKey:   key,
		AuthKeyID: key[:8],
		Salt:      42,
	}
}

// TestSessionRoundTrip ensures a printed session can be read back by gotd,
// which is what `tgauth try` does.
func TestSessionRoundTrip(t *testing.T) {
	ctx := context.Background()
	data := testSessionData()

	f, err := encodeSession(ctx, data)
	require.NoError(t, err)
	require.Equal(t, 1, f.Version)

	raw, err := json.Marshal(f)
	require.NoError(t, err)

	// gotd must accept the printed bytes as-is.
	var storage session.StorageMemory
	require.NoError(t, storage.StoreSession(ctx, raw))
	loaded, err := (&session.Loader{Storage: &storage}).Load(ctx)
	require.NoError(t, err)
	require.Equal(t, data, loaded)

	decoded, err := decodeSession(ctx, raw)
	require.NoError(t, err)
	require.Equal(t, data, decoded)
}

func TestDecodeSession(t *testing.T) {
	ctx := context.Background()
	data := testSessionData()

	current, err := json.Marshal(sessionFile{Version: 1, Data: *data})
	require.NoError(t, err)
	legacy, err := json.Marshal(data)
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   []byte
		want    *session.Data
		wantErr bool
	}{
		{"Current", current, data, false},
		{"Legacy", legacy, data, false},
		{"Empty", nil, nil, true},
		{"EmptyObject", []byte(`{}`), nil, true},
		{"NoAuthKey", []byte(`{"DC":2}`), nil, true},
		{"UnknownVersion", []byte(`{"Version":999,"Data":{"DC":2}}`), nil, true},
		{"NotJSON", []byte(`hello`), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeSession(ctx, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func FuzzDecodeSession(f *testing.F) {
	ctx := context.Background()
	data := testSessionData()

	current, err := json.Marshal(sessionFile{Version: 1, Data: *data})
	require.NoError(f, err)
	legacy, err := json.Marshal(data)
	require.NoError(f, err)

	for _, seed := range [][]byte{
		current,
		legacy,
		nil,
		[]byte(`{}`),
		[]byte(`{"DC":2}`),
		[]byte(`{"Version":999,"Data":{"DC":2}}`),
		[]byte(`hello`),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		decoded, err := decodeSession(ctx, input)
		if err != nil {
			return
		}
		// Anything decodeSession accepts must survive an encode round-trip.
		encoded, err := encodeSession(ctx, decoded)
		require.NoError(t, err)
		raw, err := json.Marshal(encoded)
		require.NoError(t, err)
		again, err := decodeSession(ctx, raw)
		require.NoError(t, err)
		require.Equal(t, decoded, again)
	})
}
