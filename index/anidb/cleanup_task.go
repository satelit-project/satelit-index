package anidb

import (
	"io/ioutil"
	"os"
	"path"

	"satelit-project/satelit-index/config"
	"satelit-project/satelit-index/models"

	"github.com/gobuffalo/pop"
)

type CleanupAnidbIndexTask struct {
	db        *pop.Connection
	anidbCfg  config.Anidb
	serverCfg config.Server
}

func NewCleanupAnidbIndexTask(db *pop.Connection, anidbCfg config.Anidb, serverCfg config.Server) CleanupAnidbIndexTask {
	return CleanupAnidbIndexTask{db, anidbCfg, serverCfg}
}

func (t CleanupAnidbIndexTask) Cleanup() error {
	filesToRetain, err := t.indexFilesToRetain()
	if err != nil {
		return err
	}

	return t.cleanup(filesToRetain)
}

func (t CleanupAnidbIndexTask) indexFilesToRetain() ([]models.AnidbIndexFile, error) {
	var indexFiles []models.AnidbIndexFile
	err := t.db.Q().All(&indexFiles)
	if err != nil {
		return nil, err
	}

	return indexFiles, nil
}

func (t CleanupAnidbIndexTask) cleanup(filesToRetain []models.AnidbIndexFile) error {
	existingFiles, err := ioutil.ReadDir(t.archivesPath())
	if err != nil {
		return err
	}

	expectedFiles := make(map[string]bool)
	for _, file := range filesToRetain {
		expectedFiles[file.Name] = true
	}

	for _, file := range existingFiles {
		if expectedFiles[file.Name()] {
			continue
		}

		err = os.Remove(file.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func (t CleanupAnidbIndexTask) archivesPath() string {
	return path.Join(t.serverCfg.FilesServePath, t.serverCfg.ArchivesServePath)
}
