package main_test

import (
	"alnoor/todo-go-htmx/server"
	"fmt"
	"net/http"
	"testing"
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

		AssertResponseContain(response, fmt.Sprintf(`value="%s"`, task.Description), t)
	})

	t.Run("list tasks when open index", func(t *testing.T) {
		task, err, _ := PostNewTask(serve)
		err, response := GetTasksList(serve)
		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContain(response, task.Description, t)
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
		task, response, err = ToggleTaskStatus(task, serve)

		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContain(response, `checked="checked"`, t)
	})

	t.Run("submit delete request to /tasks/id should delete task", func(t *testing.T) {
		task, err, response := PostNewTask(serve)
		response, err = DeleteTask(task, serve)

		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}
	})
}
