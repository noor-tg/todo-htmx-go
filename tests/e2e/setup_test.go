package e2e_test

import (
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
func (g G) page(url string) *rod.Page {
	page := g.browser.MustPage(url)
	page.MustWindowFullscreen()

	g.Cleanup(page.MustClose)

	return page
}
