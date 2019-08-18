package server

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/gobuffalo/pop"
	"github.com/satelit-project/satelit-index/config"
)

type JSONError struct {
	err error
}

type Server struct {
	config config.Server
}

func NewServer(config config.Server) *Server {
	return &Server{config: config}
}

func (s *Server) Serve() {
	dbConnection, err := pop.Connect(config.Environment())
	if err != nil {
		panic(err)
	}

	serveIndexFiles(s.config, dbConnection)
	serveStaticFiles(s.config)

	port := fmt.Sprintf(":%d", s.config.Port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

func serveIndexFiles(config config.Server, dbConnection *pop.Connection) {
	filesURL, err := url.Parse(config.FilesServeURL)
	if err != nil {
		panic(err)
	}

	filesURL.Path = config.ArchivesServePath
	service := NewIndexFilesService(dbConnection, *filesURL)
	http.HandleFunc("/index_files", service.GetIndexFile)
}

func serveStaticFiles(config config.Server) {
	fs := NewFileServer(config)
	pattern := path.Join("/", config.ArchivesServePath) + "/"
	http.Handle(pattern, fs)
}

func (e JSONError) Error() string {
	return fmt.Sprintf("{\"error\": \"%v\"}", e.err.Error())
}
