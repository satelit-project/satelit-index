package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Application configuration.
type Config struct {
	Serving  *Serving  `yaml:"serving"`
	Database *Database `yaml:"db"`
	Storage  *Storage  `yaml:"storage"`
	AniDB    *AniDB    `yaml:"anidb"`
	Logging  *Logging  `yaml:"logging"`
}

// Server configuration.
type Serving struct {
	// Port to listen for incoming connections.
	Port uint `yaml:"port"`

	// Timeout for graceful shutdown.
	HaltTimeout uint64 `yaml:"halt-timeout"`
}

// S3 configuration for the service data storage.
type Storage struct {
	// Storage access key.
	Key string `yaml:"key"`

	// Storage access secret.
	Secret string `yaml:"secret"`

	// Host to store service artifacts.
	Host string `yaml:"host"`

	// S3 bucket name.
	Bucket string `yaml:"bucket"`

	// Timeout for files uploading.
	UploadTimeout uint64 `yaml:"upload-timeout"`
}

// Database configuration.
type Database struct {
	// Database connection URL.
	URL string `yaml:"url"`
}

// AniDB specific configuration.
type AniDB struct {
	// Path to a directory with AniDB index files relative to storage path.
	StorageDir string `yaml:"storage-dir"`

	// URL from where to download database index files.
	IndexURL string `yaml:"index-url"`

	// How many seconds to wait before database index update.
	UpdateInterval uint64 `yaml:"update-interval"`
}

// Logging configuration.
type Logging struct {
	// Logging profile.
	Profile string `yaml:"profile"`
}

// Returns default app configuration or error if failed to read it.
func Default() (Config, error) {
	data := makeData(os.Environ())
	return AtPath("config/default.yml", data)
}

// Returns app configuration parsed from template with provided data.
func AtPath(path string, data map[string]string) (Config, error) {
	var cfg Config
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	content, err = render(content, data)
	if err != nil {
		return cfg, err
	}

	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Checks if storage configuration is valid and can be used.
func (s *Storage) IsValid() bool {
	if len(s.Key) == 0 || len(s.Secret) == 0 {
		return false
	}

	if len(s.Host) == 0 || len(s.Bucket) == 0 {
		return false
	}

	return true
}

// Renders template with provided data.
func render(cfg []byte, data map[string]string) ([]byte, error) {
	t, err := template.New("config").Parse(string(cfg))
	if err != nil {
		return nil, err
	}

	t.Option("missingkey=zero")

	var b bytes.Buffer
	err = t.Execute(&b, data)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Maps provided environment varialbes into template data.
func makeData(env []string) map[string]string {
	data := make(map[string]string, 8)
	for _, env := range env {
		sp := strings.Split(env, "=")
		data[sp[0]] = strings.Join(sp[1:], "=")
	}

	return data
}
