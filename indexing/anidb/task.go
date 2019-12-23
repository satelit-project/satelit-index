package anidb

import (
	"context"
	"path/filepath"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/db"
	"shitty.moe/satelit-project/satelit-index/logging"
	"shitty.moe/satelit-project/satelit-index/task"
)

// A task to update database index.
type IndexUpdateTask struct {
	downloader IndexDownloader
	cfg        config.AniDB
	db         *db.Queries
	log        *logging.Logger
}

// Downloads new database index and makes it available to external services.
func (t IndexUpdateTask) Run() error {
	downloadPath := t.cfg.Dir
	t.log.Infof("downloading new index to: %v", downloadPath)
	idxPath, err := t.downloader.Download(downloadPath)
	if err != nil {
		t.log.Errorf("index download failed: %v, err")
		return err
	}

	t.log.Infof("trying to save new index to db: %v", idxPath)
	saved, err := t.updateDB(idxPath)
	if err != nil {
		t.log.Errorf("failed to save new index: %v", err)
		return err
	}

	t.log.Infof("new index file saved, db updated: %v", saved)
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
	return true, t.db.AddIndexFile(context.Background(), hash)
}

// Factory for IndexUpdateTask task.
type IndexUpdateTaskFactory struct {
	downloader IndexDownloader
	cfg        config.AniDB
	db         *db.Queries
	log        *logging.Logger
}

// Creates new task.
func (t IndexUpdateTaskFactory) MakeTask() task.Task {
	return IndexUpdateTask{
		downloader: t.downloader,
		cfg: t.cfg,
		db: t.db,
		log: t.log.With("id", t.ID()),
	}
}

// Returns task scheduling interval.
func (t IndexUpdateTaskFactory) Interval() uint64 {
	return t.cfg.UpdateInterval
}

// Returns identificator of the produced task.
func (t IndexUpdateTaskFactory) ID() string {
	return "anidb-upd"
}
