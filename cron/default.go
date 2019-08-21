package cron

import (
	"fmt"

	"github.com/gobuffalo/pop"
	_cron "github.com/robfig/cron/v3"
	"github.com/satelit-project/satelit-index/config"
	"github.com/satelit-project/satelit-index/index/anidb"
)

func DefaultAnidbScheduler(db *pop.Connection, serverCfg config.Server, anidbCfg config.Anidb) AnidbScheduler {
	// TODO: logger
	cron := _cron.New()
	scheduler := NewAnidbScheduler(cron, anidbCfg)

	scheduler.SetUpdateIndexTask(updateIndexTask(db, serverCfg, anidbCfg))
	scheduler.SetCleanupIndexesTask(cleanupIndexesTask(db, serverCfg, anidbCfg))

	return scheduler
}

func updateIndexTask(db *pop.Connection, serverCfg config.Server, anidbCfg config.Anidb) Task {
	return func() {
		task := anidb.NewUpdateAnidbIndexTask(db, anidbCfg, serverCfg)
		_, err := task.UpdateIndex()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanupIndexesTask(db *pop.Connection, serverCfg config.Server, anidbCfg config.Anidb) Task {
	return func() {
		task := anidb.NewCleanupAnidbIndexTask(db, anidbCfg, serverCfg)
		err := task.Cleanup()
		if err != nil {
			fmt.Println(err)
		}
	}
}
