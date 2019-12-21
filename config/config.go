package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Application configuration.
type Config struct {
	Serving Serving `yaml:"serving"`
	AniDB   AniDB   `yaml:"anidb"`
}

// AniDB specific configuration.
type AniDB struct {
	// Path to a directory with AniDB index files. The path is related
	// to server's serving directory.
	Dir string `yaml:"directory"`

	// URL from where to download database index files.
	IndexURL string `yaml:"index-url"`
}

// Server configuration.
type Serving struct {
	// Port to listen for incoming connections.
	Port int `yaml:"port"`

	// Path to a directory to serve files from.
	Path string `yaml:"serve-path"`
}

// Returns default app configuration or error if failed to read it.
func Default() (*Config, error) {
	return AtPath("config/default.yml")
}

// Returns app configuration parsed from file at provided file
// or error if failed to read it.
func AtPath(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
