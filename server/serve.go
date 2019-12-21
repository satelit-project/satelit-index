package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// Server which serves anime titles index files.
type IndexServer struct {
	cfg *config.Config
	log *logging.Logger
}

// Creates new server instance with provided configuration and logger.
func New(cfg *config.Config, log *logging.Logger) IndexServer {
	return IndexServer{
		cfg: cfg,
		log: log,
	}
}

// Starts serving anime titles index files.
func (s IndexServer) Run() error {
	dir := s.cfg.Serving.Path
	port := s.cfg.Serving.Port
	s.log.Infof("serving files from %v at %d", dir, port)

	if err := s.createServeDirs(); err != nil {
		s.log.Errorf("failed to create dir %v: %v", dir, err)
		return err
	}

	addr := fmt.Sprintf(":%d", port)
	fs := http.FileServer(http.Dir(dir))

	h := LogRequest(http.StripPrefix("/index/", fs), s.log.With("dir", "index"))
	http.Handle("/index/", h)

	if err := http.ListenAndServe(addr, nil); err != nil {
		s.log.Errorf("error while serving files: %v", err)
		return err
	}

	return nil
}

// Creates required directories for serving if they are not exists.
func (s IndexServer) createServeDirs() error {
	root := s.cfg.Serving.Path
	anidb := filepath.Join(root, s.cfg.AniDB.Dir)
	perm := os.FileMode(0755)

	if err := os.MkdirAll(root, perm); err != nil {
		return err
	}

	if err := os.MkdirAll(anidb, perm); err != nil {
		return err
	}

	return nil
}
