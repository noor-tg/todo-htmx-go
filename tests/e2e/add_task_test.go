package e2e_test

import (
	"testing"

	"github.com/pioz/faker"
)

func TestAddTask(t *testing.T) {
	g := setup(t)

	p := g.page("http://localhost:3000")
	text := faker.ColorName()

	// NOTE: no need to use type key enter event
	p.MustElement("#new-task").MustInput(text)

	// NOTE: it is better to use this than text equality with got
	p.MustElementR("li", text).MustVisible()
}
