package todo

import "embed"

// must used from here to embed static dir. can not use ./ ../ in embed

//go:embed static
var Static embed.FS

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
