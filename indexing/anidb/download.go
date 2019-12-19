package anidb

import (
	"crypto/md5"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
)

type IndexDownloader struct {
	config config.Config
	log    *logging.Logger
}

func (d IndexDownloader) Download(path string) (string, error) {
	filePath, err := d.downloadIndex()
	if err != nil {
		return "", err
	}

	return d.moveIndex(filePath, path)
}

func (d IndexDownloader) downloadIndex() (string, error) {
	tmp, err := ioutil.TempFile("", "anidb_index")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(d.config.AniDB.IndexURL)
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

func (d IndexDownloader) moveIndex(indexPath, destDir string) (string, error) {
	destPath := filepath.Join(destDir, filepath.Base(indexPath))
	err := os.Rename(indexPath, destPath)
	return destPath, err
}

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
