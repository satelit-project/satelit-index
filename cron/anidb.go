package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/satelit-project/satelit-index/config"
)

type AnidbScheduler struct {
	cron   *cron.Cron
	config config.Anidb
}

func NewAnidbScheduler(cron *cron.Cron, config config.Anidb) AnidbScheduler {
	return AnidbScheduler{cron, config}
}

func (s AnidbScheduler) SetUpdateIndexTask(task Task) {
	schedule, err := cron.ParseStandard(s.config.UpdateInterval)
	if err != nil {
		panic(err)
	}

	s.cron.Schedule(schedule, task)
}

func (s AnidbScheduler) StartJobs() {
	s.cron.Start()
}

func (s AnidbScheduler) StopJobs() {
	s.cron.Stop()
}
