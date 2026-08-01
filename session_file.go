package main

import (
	"context"
	"encoding/json"

	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
)

// sessionFile is the versioned envelope [session.Loader] reads and writes.
//
// gotd keeps the version constant unexported, so encodeSession derives it by
// asking [session.Loader] to marshal the data instead of hardcoding it here.
type sessionFile struct {
	Version int
	Data    session.Data
}

// encodeSession converts data into the envelope gotd expects on disk, so that
// the printed session can be fed back into [session.Loader].
func encodeSession(ctx context.Context, data *session.Data) (sessionFile, error) {
	var storage session.StorageMemory
	if err := (&session.Loader{Storage: &storage}).Save(ctx, data); err != nil {
		return sessionFile{}, errors.Wrap(err, "save")
	}

	raw, err := storage.LoadSession(ctx)
	if err != nil {
		return sessionFile{}, errors.Wrap(err, "load")
	}

	var f sessionFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return sessionFile{}, errors.Wrap(err, "unmarshal")
	}
	return f, nil
}

// decodeSession parses a session file.
//
// It accepts the [sessionFile] envelope and, for backwards compatibility, the
// bare [session.Data] tgauth used to print.
func decodeSession(ctx context.Context, raw []byte) (*session.Data, error) {
	var storage session.StorageMemory
	if err := storage.StoreSession(ctx, raw); err != nil {
		return nil, errors.Wrap(err, "store")
	}

	data, err := (&session.Loader{Storage: &storage}).Load(ctx)
	if err == nil {
		return data, nil
	}
	// Loader reports both an absent and a version-mismatched session as
	// ErrNotFound, so a mismatch is the only hint that this may be a legacy file.
	if !errors.Is(err, session.ErrNotFound) {
		return nil, errors.Wrap(err, "load")
	}

	var legacy session.Data
	if legacyErr := json.Unmarshal(raw, &legacy); legacyErr != nil || len(legacy.AuthKey) == 0 {
		// Not a legacy file either: the envelope error describes it better.
		return nil, errors.Wrap(err, "load")
	}
	return &legacy, nil
}
