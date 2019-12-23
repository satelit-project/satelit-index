package server

import (
	"context"
	"time"
	"fmt"
	"net/http"
	"os"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// Server which serves anime titles index files.
type IndexServer struct {
	inner *http.Server
	cfg   *config.Config
	q     *db.Queries
	log   *logging.Logger
}

// Creates new server instance with provided configuration and logger.
func New(cfg *config.Config, q *db.Queries, log *logging.Logger) IndexServer {
	return IndexServer{
		cfg: cfg,
		q:   q,
		log: log.With("srv", "index"),
	}
}

// Starts serving anime titles index files.
//
// The method is not thread-safe and should only be called once.
func (s *IndexServer) Run() error {
	if s.inner != nil {
		return nil
	}

	port := s.cfg.Serving.Port
	s.log.Infof("serving at %d", port)

	mux := http.NewServeMux()
	mux.Handle("/index/", s.makeFsHandler())
	mux.Handle("/latest/anidb/", s.makeAniDBHandler())

	addr := fmt.Sprintf(":%d", port)
	s.inner = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := s.inner.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("error while serving files: %v", err)
		}
	}()

	return nil
}

// Tries to gracefully shutdown the server.
//
// The method is not thread-safe and should be called once.
func (s *IndexServer) Shutdown() error {
	if s.inner != nil {
		timeout := time.Duration(s.cfg.Serving.HaltTimeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		s.log.Infof("trying to gracefully shutdown server")
		return s.inner.Shutdown(ctx)
	}

	return nil
}

// Returns new handler for serving files.
func (s IndexServer) makeFsHandler() http.Handler {
	dir := s.cfg.Serving.Path
	fs := http.StripPrefix("/index/", http.FileServer(http.Dir(dir)))
	if err := s.createServeDirs(); err != nil {
		s.log.Fatalf("failed to create dir %v: %v", dir, err)
	}

	return LogRequest(fs, s.log)
}

// Returns new handler for getting latest anidb data.
func (s IndexServer) makeAniDBHandler() http.Handler {
	log := s.log.With("service", "anidb")
	h := latestAniDBIndexService{
		q:   s.q,
		log: log,
	}

	return LogRequest(h, s.log)
}

// Creates required directories for serving if they are not exists.
func (s IndexServer) createServeDirs() error {
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
