package main_test

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/server"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
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
		task, err, response = PutTask(task, serve)

		AssertResponseCode(response, http.StatusCreated, t)

		if err != nil {
			t.Errorf("error in post task: %v", err)
		}

		AssertResponseContain(response, task.Description, t)
	})
}

func GetTasksList(serve server.Server) (error, *httptest.ResponseRecorder) {
	request, err := http.NewRequest(http.MethodGet, "/", nil)

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)
	return err, response
}

func AssertResponseContain(response *httptest.ResponseRecorder, wanted string, t testing.TB) {
	t.Helper()
	if !strings.Contains(response.Body.String(), wanted) {
		t.Errorf("got %v wanted %v", response.Body.String(), wanted)
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
		Description: "هذه مهمة من فحص",
	}
	input := strings.NewReader(fmt.Sprintf("description=%s", task.Description))

	request, err := http.NewRequest(http.MethodPost, "/tasks", input)

	// NOTE: setting the header is required
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)

	task.Id = GetTaskId(response)

	return task, err, response
}

func PutTask(task todo.Task, serve server.Server) (todo.Task, error, *httptest.ResponseRecorder) {
	input := strings.NewReader(fmt.Sprintf("description=%s", task.Description))
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", task.Id), input)

	// NOTE: setting the header is required
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	serve.Router.ServeHTTP(response, request)

	task.Id = GetTaskId(response)

	return task, err, response
}

func GetTaskId(response *httptest.ResponseRecorder) int {
	pattern := regexp.MustCompile(`task-id="(\d+)"`)

	body := response.Body.String()

	matches := pattern.FindStringSubmatch(body)
	if len(matches) > 0 {
		id_num, _ := strconv.Atoi(matches[1])
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
