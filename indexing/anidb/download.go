package anidb

import (
	"crypto/md5"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// AniDB database dump downloader.
type IndexDownloader struct {
	indexURL string
}

// Downloads database dump to a directory with provided path.
func (d IndexDownloader) Download(path string) (string, error) {
	filePath, err := d.downloadIndex()
	if err != nil {
		return "", err
	}

	return d.moveIndex(filePath, path)
}

// Downloads database index to temporary directory.
//
// Path to the database index will be returned if it was successfully downloaded.
func (d IndexDownloader) downloadIndex() (string, error) {
	tmp, err := ioutil.TempFile("", "anidb_index")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(d.indexURL)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tmp, resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return "", err
	}

	if err = tmp.Close(); err != nil {
		return "", err
	}

	return filepath.Abs(filepath.Dir(tmp.Name()))
}

// Moves database index from indexPath to a destDir directory. The file will be
// renamed to it's MD5 hash.
func (d IndexDownloader) moveIndex(indexPath, destDir string) (string, error) {
	destPath := filepath.Join(destDir, filepath.Base(indexPath))
	err := os.Rename(indexPath, destPath)
	return destPath, err
}

// Returns MD5 hash string for a file at specified path.
func (d IndexDownloader) fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}
