package cron

import (
	"satelit-project/satelit-index/config"
	"satelit-project/satelit-index/index/anidb"
	"satelit-project/satelit-index/logging"

	"github.com/gobuffalo/pop"
	_cron "github.com/robfig/cron/v3"
)

func DefaultAnidbScheduler(db *pop.Connection, serverCfg config.Server, anidbCfg config.Anidb) AnidbScheduler {
	cron := _cron.New()
	scheduler := NewAnidbScheduler(cron, anidbCfg)

	scheduler.SetUpdateIndexTask(updateIndexTask(db, serverCfg, anidbCfg))
	scheduler.SetCleanupIndexesTask(cleanupIndexesTask(db, serverCfg, anidbCfg))

	return scheduler
}

func updateIndexTask(db *pop.Connection, serverCfg config.Server, anidbCfg config.Anidb) Task {
	return func() {
		log := logging.DefaultLogger().With("anidb-update")
		log.Info("Starting AniDB index update")

		task := anidb.NewUpdateAnidbIndexTask(db, anidbCfg, serverCfg)
		_, err := task.UpdateIndex()
		if err != nil {
			log.Error("AniDB index update failed: %v", err)
		}
	}
}

func cleanupIndexesTask(db *pop.Connection, serverCfg config.Server, anidbCfg config.Anidb) Task {
	return func() {
		log := logging.DefaultLogger().With("anidb-cleanup")
		log.Info("Starting AniDB index cleanup")

		task := anidb.NewCleanupAnidbIndexTask(db, anidbCfg, serverCfg)
		err := task.Cleanup()
		if err != nil {
			log.Error(err)
		}
	}
}
