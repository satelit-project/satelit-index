package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"

	"gopkg.in/yaml.v3"
)

func init() {
	if len(envy.Get(envEnvironment, "")) > 0 {
		var err error
		env, err = envy.MustGet(envEnvironment)
		assertNoError(err)
	}

	readConfig("config/server.yml", &serverConfig)
	readConfig("config/anidb.yml", &anidbConfig)
}

const envEnvironment = "SI_ENVIRONMENT"

var (
	env          = "development"
	serverConfig map[string]*Server
	anidbConfig  map[string]*Anidb
)

type Server struct {
	Port              int    `yaml:"port"`
	FilesServePath    string `yaml:"files-serve-path"`
	ArchivesServePath string `yaml:"archives-serve-path"`
	FilesServeURL     string `yaml:"files-serve-url"`
	LimitFiles        int    `yaml:"limit-files"`
}

type Anidb struct {
	IndexURL string `yaml:"index-url"`
}

func Environment() string {
	return env
}

func ServerConfig() *Server {
	return serverConfig[env]
}

func AnidbConfig() *Anidb {
	return anidbConfig[env]
}

func readConfig(path string, dst interface{}) {
	r, err := os.Open(path)
	assertNoError(err)

	config, err := parseConfig(r)
	assertNoError(err)

	err = yaml.Unmarshal(config, dst)
	assertNoError(err)
	assertConfigEnv(dst)
}

func parseConfig(r io.Reader) ([]byte, error) {
	tmpl := template.New("config")
	tmpl.Funcs(map[string]interface{}{
		"envOr": envy.Get,
		"env": func(s1 string) string {
			return envy.Get(s1, "")
		},
	})

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	pt, err := tmpl.Parse(string(buf))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse config template")
	}

	var bb bytes.Buffer
	err = pt.Execute(&bb, nil)
	return bb.Bytes(), err
}

func assertConfigEnv(config interface{}) {
	r := reflect.ValueOf(config)
	if r.Kind() != reflect.Ptr || r.Type().Elem().Kind() != reflect.Map {
		panic("configuration object is not a map")
	}

	for _, key := range r.Elem().MapKeys() {
		if key.String() == env {
			return
		}
	}

	panic(fmt.Sprintf("configuration not found for env: %v", env))
}

func assertNoError(e error) {
	if e != nil {
		panic(e)
	}
}
