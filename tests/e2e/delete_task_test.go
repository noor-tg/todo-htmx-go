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
	p.MustElement("#new-task").MustInput(text).Page().MustWaitRequestIdle()()

	g.Eq(p.MustElement("li").MustText(), text)

	button := `//*[@id="list"]/div[1]/button`
	p.MustElementX(button).MustClick().Page().MustWaitRequestIdle()()

	// NOTE: check for button not exist in dom
	g.Eq(len(p.MustElementsX(button)), 0)
	// NOTE: check for li not exist in list
	g.Eq(p.MustElement("#list").MustHas("li"), false)
}
