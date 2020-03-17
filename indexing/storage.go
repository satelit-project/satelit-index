package indexing

import (
	"context"
	"path"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v6"

	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// Remote storage for anime index files.
type IndexStorage struct {
	// Storage configuration.
	cfg    *config.Storage

	// Name of directory where to store index files.
	dir    string

	// S3 client.
	client *minio.Client

	// Logger
	log *logging.Logger
}

// Creates and returns new storage object or error if initialization failed.
func NewIndexStorage(cfg *config.Storage, dir string, log *logging.Logger) (IndexStorage, error) {
	client, err := minio.New(cfg.Host, cfg.Key, cfg.Secret, true)
	if err != nil {
		return IndexStorage{}, err
	}

	return IndexStorage{cfg, dir, client, log.With("index-bucket", cfg.Bucket)}, nil
}

// Uploads file at path to remote storage and returns it's URL or error if upload failed.
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

	return s.fileURL(name), nil
}

func (s IndexStorage) storageDir() string {
	return s.dir
}

func (s IndexStorage) fileURL(name string) string {
	return path.Join(s.cfg.Host, s.storageDir(), name)
}
