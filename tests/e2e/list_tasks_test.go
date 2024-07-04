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
		AddNewTaskOp(p, text)
	}

	for idx, li := range p.MustElements("li") {
		li.MustVisible()
		// NOTE: because items added in reverse to html .
		// I check it in reverse
		g.Eq(li.MustText(), list[len(list)-idx-1])
	}
}

func TestFilterByDescriptionTasks(t *testing.T) {
	g := setup(t)

	p := g.page("/")
	list := []string{faker.ColorName(), faker.ColorName(), faker.ColorName(), faker.ColorName()}

	for _, text := range list {
		AddNewTaskOp(p, text)
	}

	p.MustElement("input[name=description]").MustInput(list[1]).Page().MustWaitRequestIdle()()

	p.MustElement("li").MustVisible()
	g.Eq(p.MustElement("li").MustText(), list[1])
}

func TestFilterByStatusTasks(t *testing.T) {
	// setup
	g := setup(t)
	p := g.page("/")

	// prepare tasks descriptions
	list := []string{faker.ColorName(), faker.ColorName(), faker.ColorName(), faker.ColorName()}

	// add new task for each description in list . wait for network idle
	for _, text := range list {
		AddNewTaskOp(p, text)
	}

	// edit input checkbox
	editInput := `//*[@id="list"]/div[1]/input`
	// click the checkbox (to make it scheduled)
	p.MustElementX(editInput).MustClick().Page().MustWaitRequestIdle()()

	// filter by scheduled
	scheduleStateBtn := "/html/body/div/form[1]/div[2]/label[2]"
	p.MustElementX(scheduleStateBtn).MustClick().Page().MustWaitRequestIdle()()

	// it must be visible
	g.Eq(p.MustElement("li").MustVisible(), true)
	// last item was checked as complete
	g.Eq(p.MustElement("li").MustText(), list[len(list)-1])
	// there should be only one li in tasks list
	g.Eq(len(p.MustElements("#list li")), 1)
}

func TestFilterByDescriptionAndStatusTasks(t *testing.T) {
	// setup
	g := setup(t)
	p := g.page("/")

	// prepare tasks descriptions
	task1 := faker.ColorName()
	task2 := faker.ColorName()

	AddNewTaskOp(p, task1)
	AddNewTaskOp(p, task1)

	AddNewTaskOp(p, task2)
	AddNewTaskOp(p, task2)
	AddNewTaskOp(p, task2)

	// edit input checkbox
	editInput := `//*[@id="list"]/div[1]/input`
	// click the checkbox (to make it scheduled)
	p.MustElementX(editInput).MustClick().Page().MustWaitRequestIdle()()

	// filter by scheduled
	scheduleStateBtn := "/html/body/div/form[1]/div[2]/label[3]"
	p.MustElementX(scheduleStateBtn).MustClick().Page().MustWaitRequestIdle()()
	// show only tasks contain task2 description
	p.MustElement("input[name=description]").MustInput(task2).Page().MustWaitRequestIdle()()

	// it must be visible
	g.Eq(p.MustElement("li").MustVisible(), true)
	// last item was checked as complete
	g.Eq(p.MustElement("li").MustText(), task2)
	// there should be only one li in tasks list
	g.Eq(len(p.MustElements("#list li")), 2)
}
