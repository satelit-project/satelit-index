package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // we only need to register the driver
	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
)

const _urlFmt = "postgres://%s:%s@%s:%d/%s?sslmode=%s"

// Factory for database access classes.
type Factory struct {
	cfg *config.Database
	log *logging.Logger
}

// Creates new factory instance with provided configuration.
func NewFactory(cfg *config.Database, log *logging.Logger) Factory {
	return Factory{
		cfg: cfg,
		log: log,
	}
}

// Creates and returns new object for database queries execution.
func (f Factory) MakeQueries() (*Queries, error) {
	url := makeURL(f.cfg)
	f.log.Infof("connecting to db: %s", url)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return New(db), nil
}

// Returns Postgres connection URL built for given configuration.
func makeURL(cfg *config.Database) string {
	return fmt.Sprintf(_urlFmt, cfg.User, cfg.Passwd, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
}
