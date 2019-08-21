package db

import (
	"satelit-project/satelit-index/config"

	"github.com/gobuffalo/pop"
)

func SetupAnidbTables(db *pop.Connection, cfg config.Anidb) error {
	query := db.RawQuery("select manage_anidb_index_files_limit(?);", cfg.FilesLimit)
	return query.Exec()
}
