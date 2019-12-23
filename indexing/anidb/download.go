package anidb

import (
	"crypto/md5"
	"encoding/hex"
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
	tmp, err := ioutil.TempFile("", "anidb_index")
	if err != nil {
		return "", err
	}

	filePath, err := d.downloadIndex(tmp)
	if err != nil {
		return "", err
	}

	return d.moveIndex(filePath, path)
}

// Downloads database index and writes it to specified file.
//
// Path to the database index will be returned if it was successfully downloaded.
// Provided file will also be closed.
func (d IndexDownloader) downloadIndex(tmp *os.File) (string, error) {
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

	return tmp.Name(), nil
}

// Moves database index from idxPath to a destDir directory. The file will be
// renamed to it's MD5 hash.
func (d IndexDownloader) moveIndex(idxPath, destDir string) (string, error) {
	hash, err := d.fileHash(idxPath)
	if err != nil {
		return "", err
	}

	destPath := filepath.Join(destDir, hash)
	dest, err := os.Create(destPath)
	if err != nil {
		return "", err
	}

	src, err := os.Open(idxPath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(src, dest)
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

	return hex.EncodeToString(h.Sum(nil)), nil
}
