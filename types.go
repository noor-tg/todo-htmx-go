package todo

import "embed"

// must used from here to embed static dir. can not use ./ ../ in embed

//go:embed static
var Static embed.FS

//go:embed certs
var Certs embed.FS

type Task struct {
	Id          int
	Description string
	Status      string
}

type Config struct {
	Cleanup bool
	LogHttp bool
	DB      string
}

var ProductionCfg = Config{
	Cleanup: false,
	LogHttp: true,
	DB:      "./todo.db",
}
var TestCfg = Config{
	Cleanup: true,
	LogHttp: false,
	DB:      "./test.db",
}

type Counts struct {
	Total     int
	Completed int
}

type Store interface {
	Migrate() error
	InsertTask(description string) (Task, error)
	GetTasks(filter Task) ([]Task, error)
	UpdateTask(id int, description string) (Task, error)
	GetTaskById(id int) (Task, error)
	DeleteTask(id int) error
	ToggleTaskStatus(id int) (Task, error)
	GetTasksByStatus(status string) ([]Task, error)
	GetTasksCount() (int, error)
	GetCompletedTasksCount() (int, error)
	GetTasksCounters() (int, int, error)
	ToggleAndAnimationData(id int) (Counts, Task, int, error)
}
