package main_test

import (
	"alnoor/todo-go-htmx/server"
	"fmt"
	"net/http"
	"testing"

	"github.com/a-h/templ"
)

func TestTasksServer(t *testing.T) {
	serve := server.NewTasksServer()

	t.Run("post task and return task element", func(t *testing.T) {
		task, err, response := PostNewTask(serve)

		AssertResponseCode(response, http.StatusCreated, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContain(response, task.Description, t)
	})

	t.Run("show edit form when go to /tasks/id", func(t *testing.T) {
		task, _, _ := PostNewTask(serve)

		fmt.Printf("task id: %d\n", task.Id)
		response, err := GetUpdateTaskForm(task, serve)
		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContainProp(response, fmt.Sprintf(`value="%s"`, templ.EscapeString(task.Description)), t)
	})

	t.Run("list tasks when open index", func(t *testing.T) {
		task, err, _ := PostNewTask(serve)
		err, response := GetTasksList(serve, "")
		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContain(response, task.Description, t)
	})

	t.Run("list completed tasks", func(t *testing.T) {
		task1, err, _ := PostNewTask(serve)
		if err != nil {
			t.Errorf("error in post task: %v", err)
		}
		task2, err, _ := PostNewTask(serve)
		if err != nil {
			t.Errorf("error in post task: %v", err)
		}
		_, _, err = ToggleTaskStatus(task2, serve)
		if err != nil {
			t.Errorf("error in post task: %v", err)
		}
		err, response := GetTasksList(serve, "?status=مكتمل")
		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContain(response, task2.Description, t)
		AssertResponseNotContain(response, task1.Description, t)
	})

	t.Run("submit put request to /tasks/id with description return task element", func(t *testing.T) {
		task, err, response := PostNewTask(serve)
		task, response, err = PutTask(task, serve)

		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContain(response, task.Description, t)
	})

	t.Run("send request to /tasks/toggle-status/id return task toggled", func(t *testing.T) {
		task, err, response := PostNewTask(serve)
		if err != nil {
			t.Errorf("error in post task: %v", err)
		}
		task, response, err = ToggleTaskStatus(task, serve)

		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContainProp(response, `checked="checked"`, t)
	})

	t.Run("submit delete request to /tasks/id should delete task", func(t *testing.T) {
		task, err, response := PostNewTask(serve)
		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		response, err = DeleteTask(task, serve)

		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}
	})
}
