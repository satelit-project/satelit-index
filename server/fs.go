package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/satelit-project/satelit-index/config"
)

type dotFileHidingFileSystem struct {
	http.FileSystem
}

type dotFileHidingFile struct {
	http.File
}

func NewFileServer(config config.Server) http.Handler {
	fs := dotFileHidingFileSystem{http.Dir(config.FilesServePath)}
	return http.FileServer(fs)
}

func (fs dotFileHidingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) {
		return nil, os.ErrPermission
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return dotFileHidingFile{file}, err
}

func (f dotFileHidingFile) Readdir(n int) (fis []os.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}

	return
}

func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}
