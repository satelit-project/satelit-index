package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AniDB AniDB `yaml:"anidb"`
}

type AniDB struct {
	IndexURL string `yaml:"index-url"`
}

func Default() (*Config, error) {
	return AtPath("config/default.yml")
}

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
