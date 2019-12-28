package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// Server which serves anime titles index files.
type IndexServer struct {
	inner *http.Server
	cfg   config.Config
	q     *db.Queries
	log   *logging.Logger
}

// Creates new server instance with provided configuration and logger.
func New(cfg config.Config, q *db.Queries, log *logging.Logger) (*IndexServer, error) {
	s := IndexServer{
		cfg: cfg,
		q:   q,
		log: log.With("srv", "index"),
	}

	mux := http.NewServeMux()
	fs, err := s.makeFsHandler()
	if err != nil {
		return nil, err
	}

	mux.Handle("/index/", fs)
	mux.Handle("/latest/anidb/", s.makeAniDBHandler())

	addr := fmt.Sprintf(":%d", cfg.Serving.Port)
	s.inner = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &s, nil
}

// Starts serving anime titles index files.
func (s *IndexServer) Run() error {
	port := s.cfg.Serving.Port

	s.log.Infof("serving at %d", port)
	if err := s.inner.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Tries to gracefully shutdown the server.
func (s *IndexServer) Shutdown() error {
	timeout := time.Duration(s.cfg.Serving.HaltTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.log.Infof("trying to gracefully shutdown server")
	return s.inner.Shutdown(ctx)
}

// Returns new handler for serving files.
func (s *IndexServer) makeFsHandler() (http.Handler, error) {
	dir := s.cfg.Serving.Path
	fs := http.StripPrefix("/index/", http.FileServer(http.Dir(dir)))
	if err := s.createServeDirs(); err != nil {
		return nil, err
	}

	return LogRequest(fs, s.log), nil
}

// Returns new handler for getting latest anidb data.
func (s *IndexServer) makeAniDBHandler() http.Handler {
	log := s.log.With("service", "anidb")
	h := latestAniDBIndexService{
		q:   s.q,
		log: log,
	}

	return LogRequest(h, s.log)
}

// Creates required directories for serving if they are not exists.
func (s *IndexServer) createServeDirs() error {
	root := s.cfg.Serving.Path
	anidb := s.cfg.AniDB.Dir
	perm := os.FileMode(0755)

	if err := os.MkdirAll(root, perm); err != nil {
		return err
	}

	if err := os.MkdirAll(anidb, perm); err != nil {
		return err
	}

	return nil
}
