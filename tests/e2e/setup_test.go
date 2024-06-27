package e2e_test

import (
	"alnoor/todo-go-htmx/server"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
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
		t.Parallel() // run each test concurrently

		return G{got.New(t), browser}
	}
}()

// a helper function to create an incognito page.
func (g G) page(path string) *rod.Page {
	router := serve(g)

	page := g.browser.MustPage(router.URL(path))
	page.MustWindowFullscreen()

	g.Cleanup(func() {
		page.MustClose()
	})

	return page
}

func serve(g G) *got.Router {
	router := g.Serve()
	serve := server.NewTasksServer(server.Config{Cleanup: true, LogHttp: false})
	router.Server.Handler = serve.Router
	return router
}
