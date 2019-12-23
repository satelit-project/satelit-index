package task

// Representation of a task that can be run in background.
type Task interface {
	// Run the task.
	Run() error
}

// A factory that can create a specific task to run in background.
type Factory interface {
	// Returns identificator of tasks produced by the factory.
	ID() string

	// Creates new task to run in background.
	MakeTask() Task

	// Returns interval which represents how much seconds to wait
	// before spawning new task.
	Interval() uint64
}
