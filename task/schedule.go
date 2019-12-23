package task

import "github.com/jasonlvhit/gocron"

import "shitty.moe/satelit-project/satelit-index/logging"

// Scheduler for background tasks
type Scheduler struct {
	inner *gocron.Scheduler
	log   *logging.Logger
}

// Creates new background scheduler.
func NewScheduler(log *logging.Logger) Scheduler {
	return Scheduler{
		inner: gocron.NewScheduler(),
		log: log.With("tasks", "bg"),
	}
}

// Adds new task for background execution.
func (s Scheduler) Add(t TaskFactory) {
	s.inner.Every(t.Interval()).DoSafely(func() {
		s.log.Infof("running task: %s", t.ID())

		task := t.MakeTask()
		if err := task.Run(); err != nil {
			s.log.Errorf("task %s failed: %s", t.ID(), err)
		}
	})
}
