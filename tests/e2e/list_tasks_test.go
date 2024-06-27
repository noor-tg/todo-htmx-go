package e2e_test

import (
	"testing"

	"github.com/pioz/faker"
)

func TestListTasks(t *testing.T) {
	g := setup(t)

	p := g.page("/")
	list := []string{faker.ColorName(), faker.ColorName(), faker.ColorName(), faker.ColorName()}

	for _, text := range list {
		// NOTE: no need to use type key enter event
		p.MustElement("#new-task").MustInput(text).Page().MustWaitRequestIdle()()
	}

	for idx, li := range p.MustElements("li") {
		li.MustVisible()
		// NOTE: because items added in reverse to html .
		// I check it in reverse
		g.Eq(li.MustText(), list[len(list)-idx-1])
	}
}
