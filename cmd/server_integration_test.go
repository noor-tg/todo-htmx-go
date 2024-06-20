package main_test

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/server"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
	t.Run("list tasks when open index", func(t *testing.T) {
		task, err, _ := PostNewTask(serve)

		err, response := GetTasksList(serve)

		AssertResponseCode(response, http.StatusOK, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContain(response, task.Description, t)

	})

}

func GetTasksList(serve server.Server) (error, *httptest.ResponseRecorder) {
	get_request, err := http.NewRequest(http.MethodGet, "/", nil)

	get_response := httptest.NewRecorder()

	serve.Router.ServeHTTP(get_response, get_request)
	return err, get_response
}

func AssertResponseContain(response *httptest.ResponseRecorder, wanted string, t *testing.T) {
	if !strings.Contains(response.Body.String(), wanted) {
		t.Errorf("got %v wanted %v", response.Body.String(), wanted)
	}
}

func AssertResponseCode(response *httptest.ResponseRecorder, wanted int, t *testing.T) {
	if response.Code != wanted {
		t.Errorf("request /tasks failed: response not created, response code: %v", response.Code)
	}
}

func PostNewTask(serve server.Server) (todo.Task, error, *httptest.ResponseRecorder) {
	task := todo.Task{
		Description: "هذه مهمة من فحص",
	}
	input := strings.NewReader(fmt.Sprintf("description=%s", task.Description))

	request, err := http.NewRequest(http.MethodPost, "/tasks", input)

	// NOTE: setting the header is required
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)
	return task, err, response
}
