package cron

import (
	"satelit-project/satelit-index/config"

	"github.com/robfig/cron/v3"
)

type AnidbScheduler struct {
	cron   *cron.Cron
	config config.Anidb
}

func NewAnidbScheduler(cron *cron.Cron, config config.Anidb) AnidbScheduler {
	return AnidbScheduler{cron, config}
}

func (s AnidbScheduler) SetUpdateIndexTask(task Task) {
	s.schedule(task, s.config.UpdateTime)
}

func (s AnidbScheduler) SetCleanupIndexesTask(task Task) {
	s.schedule(task, s.config.CleanupTime)
}

func (s AnidbScheduler) StartJobs() {
	if s.config.RunTasksOnStartup {
		for _, e := range s.cron.Entries() {
			go e.Job.Run()
		}
	}

	s.cron.Start()
}

func (s AnidbScheduler) StopJobs() {
	s.cron.Stop()
}

func (s AnidbScheduler) schedule(task Task, time string) {
	schedule, err := cron.ParseStandard(time)
	if err != nil {
		panic(err)
	}

	s.cron.Schedule(schedule, task)
}
