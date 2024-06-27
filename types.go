package todo

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
