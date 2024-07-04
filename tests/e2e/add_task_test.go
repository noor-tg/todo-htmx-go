package e2e_test

import (
	"testing"

	"github.com/pioz/faker"
)

func TestAddTask(t *testing.T) {
	g := setup(t)

	p := g.page("/")
	text := faker.ColorName()

	AddNewTaskOp(p, text)

	g.Eq(p.MustElement("li").MustText(), text)
}
