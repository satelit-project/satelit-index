package server

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"go.uber.org/zap"

	"satelit-project/satelit-index/config"

	"github.com/gobuffalo/pop"
)

type JSONError struct {
	err error
}

type Server struct {
	db     *pop.Connection
	config config.Server
	logger *zap.SugaredLogger
}

func NewServer(db *pop.Connection, config config.Server) *Server {
	return &Server{db, config, nil}
}

func (s *Server) Serve() {
	router := http.NewServeMux()

	serveIndexFiles(router, s.config, s.db)
	serveStaticFiles(router, s.config)

	port := fmt.Sprintf(":%d", s.config.Port)
	err := http.ListenAndServe(port, NewLoggingHandler(s.logger, router))
	if err != nil {
		panic(err)
	}
}

func (s *Server) SetLogger(l *zap.SugaredLogger) {
	s.logger = l
}

func serveIndexFiles(r *http.ServeMux, config config.Server, dbConnection *pop.Connection) {
	filesURL, err := url.Parse(config.FilesServeURL)
	if err != nil {
		panic(err)
	}

	filesURL.Path = config.ArchivesServePath
	service := NewIndexFilesService(dbConnection, *filesURL)
	r.HandleFunc("/index/anidb", service.GetIndexFile)
}

func serveStaticFiles(r *http.ServeMux, config config.Server) {
	fs := NewFileServer(config)
	pattern := path.Join("/", config.ArchivesServePath) + "/"
	r.Handle(pattern, fs)
}

func (e JSONError) Error() string {
	return fmt.Sprintf("{\"error\": \"%v\"}", e.err.Error())
}
