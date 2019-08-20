package cron

import (
	"github.com/gobuffalo/pop"
	_cron "github.com/robfig/cron/v3"
	"github.com/satelit-project/satelit-index/config"
	"github.com/satelit-project/satelit-index/index/anidb"
)

func DefaultAnidbScheduler(db *pop.Connection, serverCfg config.Server, anidbCfg config.Anidb) AnidbScheduler {
	// TODO: logger
	cron := _cron.New()
	scheduler := NewAnidbScheduler(cron, anidbCfg)

	updateIndex := updateIndexTask(db, serverCfg, anidbCfg)
	go updateIndex()
	scheduler.SetUpdateIndexTask(updateIndex)

	return scheduler
}

func updateIndexTask(db *pop.Connection, serverCfg config.Server, anidbCfg config.Anidb) Task {
	return func() {
		task := anidb.NewUpdateAnidbIndexTask(db, anidbCfg, serverCfg)
		_, err := task.UpdateIndex()
		if err != nil {
			println(err)
		}
	}
}
