package anidb

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	pathutils "path"

	"satelit-project/satelit-index/config"
	"satelit-project/satelit-index/models"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

type UpdateAnidbIndexTask struct {
	db        *pop.Connection
	anidbCfg  config.Anidb
	serverCfg config.Server
}

func NewUpdateAnidbIndexTask(db *pop.Connection, anidbCfg config.Anidb, serverCfg config.Server) UpdateAnidbIndexTask {
	return UpdateAnidbIndexTask{db: db, anidbCfg: anidbCfg, serverCfg: serverCfg}
}

func (t UpdateAnidbIndexTask) UpdateIndex() (*models.AnidbIndexFile, error) {
	name, err := downloadIndexFile(t.anidbCfg.IndexURL, t.saveDirectory())
	if err != nil {
		return nil, err
	}

	hash, err := fileHash(t.filePath(name))
	if err != nil {
		return nil, err
	}

	return t.addUpdatedIndex(name, hash)
}

func (t UpdateAnidbIndexTask) addUpdatedIndex(name string, hash string) (*models.AnidbIndexFile, error) {
	var indexFile *models.AnidbIndexFile
	var err error

	err = t.db.Transaction(func(conn *pop.Connection) error {
		indexFile, err = createIndexIfNeeded(conn, name, hash)
		return err
	})

	if err != nil {
		return nil, err
	}

	return indexFile, err
}

func (t UpdateAnidbIndexTask) saveDirectory() string {
	return pathutils.Join(t.serverCfg.FilesServePath, t.serverCfg.ArchivesServePath)
}

func (t UpdateAnidbIndexTask) filePath(name string) string {
	return pathutils.Join(t.saveDirectory(), name)
}

func createIndexIfNeeded(conn *pop.Connection, name string, hash string) (*models.AnidbIndexFile, error) {
	count, err := conn.Q().Where("hash = ?", hash).Count(models.AnidbIndexFile{})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return createNewIndex(conn, name, hash)
	}

	err = updateExistingIndex(conn, name, hash)
	return nil, err
}

func createNewIndex(conn *pop.Connection, name string, hash string) (*models.AnidbIndexFile, error) {
	indexFile := models.AnidbIndexFile{Name: name, Hash: hash}
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	indexFile.ID = id
	err = conn.Create(&indexFile)
	if err != nil {
		return nil, err
	}

	return &indexFile, nil
}

func updateExistingIndex(conn *pop.Connection, name string, hash string) error {
	var indexFile models.AnidbIndexFile
	err := conn.Q().Where("hash = ?", hash).First(&indexFile)
	if err != nil {
		return err
	}

	indexFile.Name = name
	return conn.Update(&indexFile)
}

func fileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
