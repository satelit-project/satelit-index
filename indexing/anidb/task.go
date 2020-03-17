package anidb

import (
	"context"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
	"shitty.moe/satelit-project/satelit-index/task"
)

// A task to update database index.
type IndexUpdateTask struct {
	updater IndexUpdater
	cfg     *config.AniDB
	db      DBQueries
	log     *logging.Logger
}

// Downloads new database index and makes it available to external services.
func (t IndexUpdateTask) Run() error {
	t.log.Infof("updating anime index")
	idx, err := t.updater.Update()
	if err != nil {
		t.log.Errorf("index download failed: %v", err)
		return err
	}

	t.log.Infof("trying to save new index to db: %v", idx)
	saved, err := t.updateDB(idx)
	if err != nil {
		t.log.Errorf("failed to save new index: %v", err)
		return err
	}

	t.log.Infof("new index file saved, db updated: %v", saved)
	return nil
}

// Inserts a record about index file to database if not already exists.
func (t IndexUpdateTask) updateDB(idx AniDBIndex) (saved bool, err error) {
	count, err := t.db.CountIndexFiles(context.Background(), idx.Hash)
	if err != nil || count > 0 {
		return false, err
	}

	// there's a data race here but this is fine since
	// db schema has ON CONFLICT DO NOTHING
	return true, t.db.AddIndexFile(context.Background(), idx)
}

// Factory for IndexUpdateTask task.
type IndexUpdateTaskFactory struct {
	Cfg     *config.AniDB
	DB      DBQueries
	Storage RemoteStorage
	Log     *logging.Logger
}

// Creates new task.
func (t IndexUpdateTaskFactory) MakeTask() task.Task {
	d := IndexUpdater{t.Cfg.IndexURL, t.Storage}
	return IndexUpdateTask{
		updater: d,
		cfg:     t.Cfg,
		db:      t.DB,
		log:     t.Log.With("id", t.ID()),
	}
}

// Returns task scheduling interval.
func (t IndexUpdateTaskFactory) Interval() uint64 {
	return t.Cfg.UpdateInterval
}

// Returns identificator of the produced task.
func (t IndexUpdateTaskFactory) ID() string {
	return "anidb-upd"
}
