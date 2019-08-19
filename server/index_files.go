package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/gobuffalo/pop"
	"github.com/satelit-project/satelit-index/models"
)

type IndexFilesService struct {
	dbConnection *pop.Connection
	serveURL     url.URL
}

type indexFileWithURL struct {
	models.IndexFile
	IndexURL string `json:"url"`
}

func NewIndexFilesService(dbConnection *pop.Connection, serveURL url.URL) IndexFilesService {
	return IndexFilesService{dbConnection, serveURL}
}

func (s IndexFilesService) GetIndexFile(w http.ResponseWriter, r *http.Request) {
	var indexFile models.IndexFile
	err := pop.Q(s.dbConnection).Order("created_at desc").First(&indexFile)

	if err != nil {
		writeError(err, w, http.StatusInternalServerError)
		return
	}

	indexURL := s.serveURL
	indexURL.Path = path.Join(indexURL.Path, indexFile.Name)
	result := indexFileWithURL{indexFile, indexURL.String()}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		writeError(err, w, http.StatusInternalServerError)
		return
	}
}

func writeError(err error, w http.ResponseWriter, statusCode int) {
	err = JSONError{err}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintln(w, err.Error())
}