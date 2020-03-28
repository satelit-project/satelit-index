package indexing

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v6"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// Remote storage for anime index files.
type IndexStorage struct {
	// Storage configuration.
	cfg *config.Storage

	// Name of directory where to store index files.
	dir string

	// S3 client.
	client *minio.Client

	// Logger
	log *logging.Logger
}

// Creates and returns new storage object or error if initialization failed.
func NewIndexStorage(cfg *config.Storage, dir string, log *logging.Logger) (IndexStorage, error) {
	secure := true
	if strings.HasPrefix(cfg.Host, "localhost") || strings.HasPrefix(cfg.Host, "127.0.0.1") {
		secure = false
	}

	client, err := minio.New(cfg.Host, cfg.Key, cfg.Secret, secure)
	if err != nil {
		return IndexStorage{}, err
	}

	return IndexStorage{cfg, dir, client, log.With("storage", cfg.Bucket)}, nil
}

// Uploads file at path to remote storage and returns it's path or error if upload failed.
func (s IndexStorage) UploadFile(localPath, contentType string) (string, error) {
	name := filepath.Base(localPath)
	remotePath := filepath.Join(s.storageDir(), name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.UploadTimeout)*time.Second)
	defer cancel()

	s.log.Infof("uploading file %s to %s", name, remotePath)
	_, err := s.client.FPutObjectWithContext(ctx, s.cfg.Bucket, remotePath, localPath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		s.log.Errorf("failed to upload %s: %v", name, err)
		return "", err
	}

	return remotePath, nil
}

func (s IndexStorage) storageDir() string {
	return s.dir
}

// Storage that does not save anything
type NoStorage struct {
	// Logger
	log *logging.Logger
}

// Creates and returns new storage object.
func NewNoStorage(log *logging.Logger) NoStorage {
	return NoStorage{log.With("storage", "void")}
}

// Pretends that uploads file at path to remote storage and returns empty URL.
func (s NoStorage) UploadFile(localPath, contentType string) (string, error) {
	s.log.Infof("uploading file %s into void", localPath)
	return "", nil
}
