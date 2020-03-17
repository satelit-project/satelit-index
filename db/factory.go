package db

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq" // we only need to register the driver
	"shitty.moe/satelit-project/satelit-index/config"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// Factory for database access classes.
type Factory struct {
	cfg  *config.Database
	log  *logging.Logger
	conn *Queries
	mx   *sync.Mutex
}

// Creates new factory instance with provided configuration.
func NewFactory(cfg *config.Database, log *logging.Logger) Factory {
	return Factory{
		cfg: cfg,
		log: log,
		mx:  &sync.Mutex{},
	}
}

// Creates and returns new object for database queries execution.
func (f Factory) MakeQueries() (*Queries, error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	if f.conn != nil {
		return f.conn, nil
	}

	url := f.cfg.URL
	f.log.Infof("connecting to db: %s", url)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	f.conn = New(db)
	return f.conn, nil
}
