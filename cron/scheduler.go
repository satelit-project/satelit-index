package cron

type Task func()

type Scheduler interface {
	StartJobs()
	StopJobs()
}

func (t Task) Run() {
	t()
}
