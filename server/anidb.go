package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// Service for retrieving latest AniDB anime index.
type aniDBIndexService struct {
	path string
	q    *db.Queries
	log  *logging.Logger
}

func (s aniDBIndexService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		tag := strings.TrimPrefix(r.URL.Path, s.path)
		idx, err := s.fetchIndex(tag)
		if err != nil {
			s.reportError(w, err)
			return
		}

		data, err := json.Marshal(idx)
		if err != nil {
			s.reportError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(data); err != nil {
			s.log.Errorf("failed to send response: %v", err)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// Fetches AniDB index file by it's tag.
//
// Tag can be either `latest` to fetch latest index file or hash of index
// file to retrieve.
func (s aniDBIndexService) fetchIndex(tag string) (db.AnidbIndexFile, error) {
	switch tag {
	case "latest":
		return s.q.LatestIndexFile(context.Background())

	default:
		return s.q.IndexFileByHash(context.Background(), tag)
	}
}

func (s aniDBIndexService) reportError(w http.ResponseWriter, err error) {
	s.log.Errorf("failed to query db: %v", err)
	http.Error(w, `{"error": "could not retireve data"}`, http.StatusInternalServerError)
}
