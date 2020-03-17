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

type IndexFile struct {
	Hash string
	URL  string
}

// AniDB database dump updater.
type IndexUpdater struct {
	indexURL string
	storage  RemoteStorage
}

// Downloads and stores latest AniDB anime dump.
func (d IndexUpdater) Update() (IndexFile, error) {
	tmp, err := ioutil.TempFile("", "anidb_index")
	if err != nil {
		return IndexFile{}, err
	}

	filePath, err := d.downloadIndex(tmp)
	if err != nil {
		return IndexFile{}, err
	}

	return d.saveIndex(filePath)
}

// Downloads database index and writes it to specified file.
//
// Path to the database index will be returned if it was successfully downloaded.
// Provided file will also be closed.
func (d IndexUpdater) downloadIndex(tmp *os.File) (string, error) {
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

// Uploads index file to remote storage and returns it's info.
func (d IndexUpdater) saveIndex(idxPath string) (IndexFile, error) {
	idx := IndexFile{}
	hash, err := d.fileHash(idxPath)
	if err != nil {
		return idx, err
	}

	destPath := filepath.Join(filepath.Dir(idxPath), hash)
	if err := os.Rename(idxPath, destPath); err != nil {
		return idx, err
	}

	url, err := d.storage.UploadFile(destPath, "application/gzip")
	if err != nil {
		return idx, err
	}

	idx.Hash = hash
	idx.URL = url
	return idx, nil
}

// Returns MD5 hash string for a file at specified path.
func (d IndexUpdater) fileHash(path string) (string, error) {
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
