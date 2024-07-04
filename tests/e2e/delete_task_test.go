package e2e_test

import (
	"testing"

	"github.com/pioz/faker"
)

func TestDeleteTask(t *testing.T) {
	g := setup(t)

	p := g.page("/")
	text := faker.ColorName()

	// NOTE: no need to use type key enter event
	AddNewTaskOp(p, text)

	button := DeleteTaskOp(g, p, text)
	AssertElNotExist(g, p, button)
	// NOTE: check for li not exist in list
	g.Eq(p.MustElement("#list").MustHas("li"), false)
}
