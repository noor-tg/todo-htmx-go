package e2e_test

import (
	"testing"

	"github.com/pioz/faker"
)

func TestAddTask(t *testing.T) {
	g := setup(t)

	p := g.page("/")
	text := faker.ColorName()

	// NOTE: no need to use type key enter event
	p.MustElement("#new-task").MustInput(text).Page().MustWaitRequestIdle()()

	g.Eq(p.MustElement("li").MustText(), text)
}
