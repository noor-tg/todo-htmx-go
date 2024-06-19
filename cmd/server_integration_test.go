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
	t.Run("post task and return task element", func(t *testing.T) {
		task := todo.Task{
			Description: "هذه مهمة من فحص",
		}
		input := strings.NewReader(fmt.Sprintf("description=%s", task.Description))
		serve := server.NewTasksServer()

		request, err := http.NewRequest(http.MethodPost, "/tasks", input)
		// NOTE: setting the header is required
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response := httptest.NewRecorder()

		serve.Router.ServeHTTP(response, request)

		if response.Code != http.StatusCreated {
			t.Errorf("request /tasks failed: response not created, response code: %v", response.Code)
		}

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		if !strings.Contains(response.Body.String(), task.Description) {
			t.Errorf("got %v wanted %v", response.Body.String(), task.Description)
		}
	})
	t.Run("list tasks when open index", func(t *testing.T) {
		task := todo.Task{
			Description: "هذه مهمة من فحص",
		}
		input := strings.NewReader(fmt.Sprintf("description=%s", task.Description))
		serve := server.NewTasksServer()

		post_request, err := http.NewRequest(http.MethodPost, "/tasks", input)
		// NOTE: setting the header is required
		post_request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		post_response := httptest.NewRecorder()

		serve.Router.ServeHTTP(post_response, post_request)

		get_request, err := http.NewRequest(http.MethodGet, "/", nil)

		get_response := httptest.NewRecorder()

		serve.Router.ServeHTTP(get_response, get_request)

		if get_response.Code != http.StatusOK {
			t.Errorf("request / failed: response not created, response code: %v", get_response.Code)
		}

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		// if !strings.Contains(response.Body.String(), task.Description) {
		// 	t.Errorf("got %v wanted %v", response.Body.String(), task.Description)
		// }
	})

}
