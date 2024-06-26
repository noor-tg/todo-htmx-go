package e2e_test

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/server"
	"os"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/google/uuid"
	"github.com/ysmood/got"
)

// test context.
type G struct {
	got.G

	browser *rod.Browser
}

// setup for tests.
var setup = func() func(t *testing.T) G {
	u := launcher.New().Headless(true).Bin("brave").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	return func(t *testing.T) G {
		// NOTE: run in Parallel has problem with testing db.
		// so stop using it until make new way to make db for each test
		// FIXED: by using uuid for db name
		t.Parallel() // run each test concurrently

		return G{got.New(t), browser}
	}
}()

// a helper function to create an incognito page.
func (g G) page(path string) *rod.Page {
	// db file for test
	db := "tmp/" + uuid.NewString() + ".db"
	router := serve(g, db)

	page := g.browser.MustPage(router.URL(path))
	page.MustWindowFullscreen()

	g.Cleanup(func() {
		page.MustClose()
		// remove db file
		os.Remove(db)
	})

	return page
}

func serve(g G, db string) *got.Router {
	router := g.Serve()
	cfg := todo.TestCfg
	cfg.DB = db
	serve := server.NewTasksServer(cfg)
	router.Server.Handler = serve.Router
	return router
}
