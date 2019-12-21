package anidb

import (
	"context"
	"path/filepath"

	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// A task to update database index.
type IndexUpdateTask struct {
	downloader   IndexDownloader
	downloadPath string
	db           *db.Queries
	logger       *logging.Logger
}

// Downloads new database index and makes it available to external services.
func (t IndexUpdateTask) Run() error {
	t.logger.Infof("downloading new index to: %v", t.downloadPath)
	idxPath, err := t.downloader.Download(t.downloadPath)
	if err != nil {
		t.logger.Errorf("index download failed: %v, err")
		return err
	}

	t.logger.Infof("trying to save new index to db: %v", idxPath)
	saved, err := t.updateDB(idxPath)
	if err != nil {
		t.logger.Errorf("failed to save new index: %v", err)
		return err
	}

	t.logger.Infof("new index file saved, db updated: %v", saved)
	return nil
}

// Inserts a record about index file to database if not already exists.
func (t IndexUpdateTask) updateDB(idxPath string) (saved bool, err error) {
	hash := filepath.Base(idxPath)
	count, err := t.db.CountIndexFiles(context.Background(), hash)
	if err != nil || count > 0 {
		return false, err
	}

	// there's a data race here but this is fine since
	// db schema has ON CONFLICT DO NOTHING
	args := db.AddIndexFileParams{
		Name: hash,
		Hash: hash,
	}
	return true, t.db.AddIndexFile(context.Background(), args)
}
