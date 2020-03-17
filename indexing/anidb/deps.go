package anidb

import (
	"context"

	"shitty.moe/satelit-project/satelit-index/db"
)

// Represents remote file storage.
type RemoteStorage interface {
	// Uploads file and returns it's URL.
	UploadFile(path, contentType string) (string, error)
}

// Represents AniDB database queries.
type DBQueries interface {
	CountIndexFiles(ctx context.Context, hash string) (int64, error)
	AddIndexFile(ctx context.Context, idx AniDBIndex) error
}

// Adapter for the db.Queries object.
type AniDBQueries struct {
	Q *db.Queries
}

func (q AniDBQueries) CountIndexFiles(ctx context.Context, hash string) (int64, error) {
	return q.Q.CountIndexFiles(ctx, hash)
}

func (q AniDBQueries) AddIndexFile(ctx context.Context, idx AniDBIndex) error {
	params := db.AddIndexFileParams{Hash: idx.Hash, Url: idx.URL}
	return q.Q.AddIndexFile(ctx, params)
}
