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
	Serving  Serving  `yaml:"serving"`
	Database Database `yaml:"db"`
	AniDB    AniDB    `yaml:"anidb"`
}

// Server configuration.
type Serving struct {
	// Port to listen for incoming connections.
	Port int `yaml:"port"`

	// Path to a directory to serve files from.
	Path string `yaml:"serve-path"`
}

type Database struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
	Passwd  string `yaml:"passwd"`
	SSLMode string `yaml:"ssl-mode"`
}

// AniDB specific configuration.
type AniDB struct {
	// Path to a directory with AniDB index files. The path is related
	// to server's serving directory.
	Dir string `yaml:"directory"`

	// URL from where to download database index files.
	IndexURL string `yaml:"index-url"`

	// How many seconds to wait before database index update.
	UpdateInterval uint64 `yaml:"update-interval"`
}

// Returns default app configuration or error if failed to read it.
func Default() (*Config, error) {
	data := makeData(os.Environ())
	return AtPath("config/default.yml", data)
}

// Returns app configuration parsed from template with provided data.
func AtPath(path string, data map[string]string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content, err = render(content, data)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Renders template with provided data.
func render(cfg []byte, data map[string]string) ([]byte, error) {
	t, err := template.New("config").Parse(string(cfg))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	t.Execute(&b, data)

	return b.Bytes(), nil
}

// Maps provided environment varialbes into template data.
func makeData(env []string) map[string]string {
	data := make(map[string]string, 8)
	for _, env := range env {
		sp := strings.Split(env, "=")
		if len(sp) != 2 {
			continue
		}

		data[sp[0]] = sp[1]
	}

	return data
}
