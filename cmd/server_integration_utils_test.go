package main_test

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/server"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/pioz/faker"
)

func GetTasksList(serve server.Server, querystring string) (error, *httptest.ResponseRecorder) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", querystring), nil)

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)
	return err, response
}

func AssertResponseContain(response *httptest.ResponseRecorder, wanted string, t testing.TB) {
	t.Helper()
	if !strings.Contains(response.Body.String(), templ.EscapeString(wanted)) {
		t.Errorf("got %v wanted %v", response.Body.String(), templ.EscapeString(wanted))
	}
}

func AssertResponseContainProp(response *httptest.ResponseRecorder, wanted string, t testing.TB) {
	t.Helper()
	if !strings.Contains(response.Body.String(), wanted) {
		t.Errorf("got %v wanted %v", response.Body.String(), wanted)
	}
}

func AssertResponseNotContain(response *httptest.ResponseRecorder, wanted string, t testing.TB) {
	t.Helper()
	if strings.Contains(response.Body.String(), templ.EscapeString(wanted)) {
		t.Errorf("got %v wanted %v", response.Body.String(), templ.EscapeString(wanted))
	}
}

func AssertResponseCode(response *httptest.ResponseRecorder, wanted int, t testing.TB) {
	t.Helper()

	if response.Code != wanted {
		t.Errorf("request /tasks failed: response not created, response code: %v", response.Code)
	}
}

func PostNewTask(serve server.Server) (todo.Task, error, *httptest.ResponseRecorder) {
	task := todo.Task{
		Description: faker.FullName(),
	}
	input := strings.NewReader(fmt.Sprintf("description=%s", task.Description))

	request, err := http.NewRequest(http.MethodPost, "/tasks", input)

	// NOTE: setting the header is required
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)

	task.Id = GetTaskId(response)

	if err != nil {
		log.Println(err)
	}

	return task, err, response
}

func PutTask(task todo.Task, serve server.Server) (todo.Task, *httptest.ResponseRecorder, error) {
	input := strings.NewReader(fmt.Sprintf("description=%s", task.Description))
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", task.Id), input)

	// NOTE: setting the header is required
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)

	task.Id = GetTaskId(response)

	return task, response, err
}

func ToggleTaskStatus(task todo.Task, serve server.Server) (todo.Task, *httptest.ResponseRecorder, error) {
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/toggle-status/%d", task.Id), nil)

	// NOTE: setting the header is required
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)

	task.Id = GetTaskId(response)

	return task, response, err
}

func DeleteTask(task todo.Task, serve server.Server) (*httptest.ResponseRecorder, error) {
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", task.Id), nil)

	// NOTE: setting the header is required
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)

	task.Id = GetTaskId(response)

	return response, err
}

func GetTaskId(response *httptest.ResponseRecorder) int {
	pattern := regexp.MustCompile(`hx-(put|post|get|delete)="/tasks/(\d+)"`)

	body := response.Body.String()

	matches := pattern.FindStringSubmatch(body)
	if len(matches) > 0 {
		id_num, _ := strconv.Atoi(matches[2])
		return id_num
	}

	return 0
}

func GetUpdateTaskForm(task todo.Task, serve server.Server) (*httptest.ResponseRecorder, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", task.Id), nil)
	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)
	return response, err
}
