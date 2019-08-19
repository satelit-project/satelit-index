package anidb

import (
	"fmt"
	"io"
	"net/http"
	"os"
	pathutils "path"
	"time"
)

func downloadIndexFile(url string, dir string) (string, error) {
	path, name := indexFilePath(dir)
	err := downloadFile(url, path)
	if err != nil {
		return "", err
	}

	return name, nil
}

func downloadFile(url string, path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}

	defer func() {
		fileErr := file.Close()
		if err == nil {
			err = fileErr
		}
	}()

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}

	_ = resp.Body.Close()
	return
}

func indexFilePath(dir string) (path string, name string) {
	timestamp := time.Now().Unix()
	name = fmt.Sprintf("%v-anidb.xml.gz", timestamp)
	path = pathutils.Join(dir, name)
	return
}
