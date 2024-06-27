package e2e_test

import (
	"testing"

	"github.com/pioz/faker"
)

func TestEditTask(t *testing.T) {
	g := setup(t)

	p := g.page("/")
	text := faker.ColorName()

	// NOTE: no need to use type key enter event
	// NOTE: ()() to start wait for networkIlde instead of preparing one
	p.MustElement("#new-task").MustInput(text).Page().MustWaitRequestIdle()()

	li := `//*[@id="list"]/div[1]/li`
	editInput := `//*[@id="list"]/div[1]/input`

	p.MustElementX(li).MustClick().Page().MustWaitRequestIdle()()

	newText := faker.ColorName()
	// mustselect to to select existing text. input "" to remove old text
	// then input new text
	p.MustElementX(editInput).MustSelectAllText().MustInput("").MustInput(newText).Page().MustWaitRequestIdle()()

	g.Eq(p.MustElement("li").MustText(), newText)
}
