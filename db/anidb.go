package db

import (
	"github.com/gobuffalo/pop"
	"github.com/satelit-project/satelit-index/config"
)

func SetupAnidbTables(db *pop.Connection, cfg config.Anidb) error {
	query := db.RawQuery("select manage_anidb_index_files_limit(?);", cfg.FilesLimit)
	return query.Exec()
}
