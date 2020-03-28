package anidb

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
)

type IndexFile struct {
	Hash     string
	FilePath string
}

// AniDB database dump updater.
type IndexUpdater struct {
	indexURL string
	client   *http.Client
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

	valid, err := isGzip(filePath)
	if err != nil {
		return IndexFile{}, err
	}
	if !valid {
		return IndexFile{}, errors.New("index file is not an archive")
	}

	return d.saveIndex(filePath)
}

// Downloads database index and writes it to specified file.
//
// Path to the database index will be returned if it was successfully downloaded.
// Provided file will also be closed.
func (d IndexUpdater) downloadIndex(tmp *os.File) (string, error) {
	resp, err := d.client.Get(d.indexURL)
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

	remotePath, err := d.storage.UploadFile(destPath, "application/gzip")
	if err != nil {
		return idx, err
	}

	idx.Hash = hash
	idx.FilePath = remotePath
	return idx, nil
}

// Returns MD5 hash string for a file at specified path.
func (d IndexUpdater) fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func isGzip(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}

	head := make([]byte, 262)
	f.Read(head)

	return filetype.IsArchive(head), nil
}
