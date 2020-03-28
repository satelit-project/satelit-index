package anidb

import (
	"context"
	"net/http"

	"github.com/corpix/uarand"

	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/indexing"
)

// Represents remote file storage.
type RemoteStorage interface {
	// Uploads file and returns it's URL.
	UploadFile(path, contentType string) (string, error)
}

// Represents AniDB database queries.
type DBQueries interface {
	// Returns number of index files with specified hash.
	CountIndexFiles(ctx context.Context, hash string) (int64, error)

	// Adds new index file to database.
	//
	// If there's already index file in the database with the same name
	// then nothing will be added or updated. The index file will be ignored.
	AddIndexFile(ctx context.Context, idx IndexFile) error
}

// Adapter for the db.Queries object.
type Queries struct {
	Q *db.Queries
}

func (q Queries) CountIndexFiles(ctx context.Context, hash string) (int64, error) {
	return q.Q.CountIndexFiles(ctx, hash)
}

func (q Queries) AddIndexFile(ctx context.Context, idx IndexFile) error {
	params := db.AddIndexFileParams{Hash: idx.Hash, Source: int32(indexing.Anidb), FilePath: idx.FilePath}
	return q.Q.AddIndexFile(ctx, params)
}

// Returns new HTTP client with fake User-Agent.
func NewFakeClient() *http.Client {
	client := &http.Client{Transport: &fakeUATransport{http.DefaultTransport}}
	return client
}

// Wrapper around HTTP client transport which adds custom User-Agent header.
type fakeUATransport struct {
	T http.RoundTripper
}

func (t *fakeUATransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ua := uarand.GetRandom()
	req.Header.Add("User-Agent", ua)
	return t.T.RoundTrip(req)
}
