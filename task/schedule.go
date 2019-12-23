package task

import (
	"github.com/jasonlvhit/gocron"
	"shitty.moe/satelit-project/satelit-index/logging"
)

// Scheduler for background tasks
type Scheduler struct {
	inner  *gocron.Scheduler
	cancel chan bool
	log    *logging.Logger
}

// Creates new background scheduler.
func NewScheduler(log *logging.Logger) Scheduler {
	return Scheduler{
		inner:  gocron.NewScheduler(),
		cancel: nil,
		log:    log.With("tasks", "bg"),
	}
}

// Adds new task for background execution.
func (s Scheduler) Add(t TaskFactory) {
	s.inner.Every(t.Interval()).Seconds().DoSafely(func(t TaskFactory) {
		s.log.Infof("running task: %s", t.ID())

		task := t.MakeTask()
		if err := task.Run(); err != nil {
			s.log.Errorf("task %s failed: %s", t.ID(), err)
		}
	}, t)
}

// Starts scheduler.
//
// After the scheduler is started it will start spawning tasks
// based on their desired iterval. The method is not thread-safe.
func (s *Scheduler) Start() {
	if s.cancel != nil {
		return
	}

	s.cancel = s.inner.Start()
}

// Stops scheduler.
//
// The method will stop spawning new tasks but already spawned tasks
// will continue execution until finished. The method is not thread-safe.
func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel <- true
		s.cancel = nil
	}
}
