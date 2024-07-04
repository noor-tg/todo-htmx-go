package e2e_test

import (
	"testing"

	"github.com/go-rod/rod/lib/input"
	"github.com/pioz/faker"
)

const editInput = `//*[@id="list"]/div[1]/form/input`

func TestEditTask(t *testing.T) {
	g := setup(t)

	p := g.page("/")
	text := faker.ColorName()

	AddNewTaskOp(p, text)

	li := `//*[@id="list"]/div[1]/li`

	p.MustElementX(li).MustClick().Page().MustWaitRequestIdle()()

	newText := faker.ColorName()
	// mustselect to to select existing text. input "" to remove old text
	// then input new text
	p.MustElementX(editInput).MustSelectAllText().MustInput("").MustInput(newText).MustType(input.Enter).Page().MustWaitRequestIdle()()

	g.Eq(p.MustElement("li").MustText(), newText)
}

func TestToggleTaskStatus(t *testing.T) {
	g := setup(t)

	p := g.page("/")
	text := faker.ColorName()

	AddNewTaskOp(p, text)

	statusInput := `//*[@id="list"]/div/input`
	// click the checkbox
	p.MustElementX(statusInput).MustClick().Page().MustWaitRequestIdle()()

	// NOTE: use * to dereference string pointer returned from the function
	g.Eq(*p.MustElementX(statusInput).MustAttribute("checked"), "checked")
}
