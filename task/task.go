package task

// Representation of a task that can be run in background.
type Task interface {
	// Run the task.
	Run() error
}

// A factory that can create a specific task to run in background.
type TaskFactory interface {
	// Creates new task to run in background.
	MakeTask() Task

	// Returns interval which represents how much seconds to wait
	// before spawning new task.
	Interval() uint64
}

// A function that implements Task interface.
type TaskFunc func() error

func (t TaskFunc) Run() error {
	return t()
}
