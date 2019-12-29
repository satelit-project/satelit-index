package server

import (
	"context"
	"encoding/json"
	"net/http"

	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/logging"
)

type latestAniDBIndexService struct {
	q   *db.Queries
	log *logging.Logger
}

func (s latestAniDBIndexService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		idx, err := s.q.LatestIndexFile(context.Background())
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

func (s latestAniDBIndexService) reportError(w http.ResponseWriter, err error) {
	s.log.Errorf("failed to query db: %v", err)
	http.Error(w, `{"error": "could not retireve data"}`, http.StatusInternalServerError)
}
