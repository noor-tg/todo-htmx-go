package server_echo

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/views"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func (s *Server) IndexHandler(c echo.Context) error {
	activeStatus := c.QueryParam("status")
	description := c.QueryParam("description")

	var tasks []todo.Task
	var err error

	if activeStatus == "الكل" {
		activeStatus = ""
	}

	var search todo.Task

	search.Description = description
	search.Status = activeStatus

	tasks, err = s.Store.GetTasks(search)

	if err != nil {
		return Render(c, http.StatusUnprocessableEntity, views.ServerError())
	}

	total, completed, err := s.Store.GetTasksCounters()
	counts := todo.Counts{Total: total, Completed: completed}

	if err != nil {
		return Render(c, http.StatusUnprocessableEntity, views.ServerError())
	}

	if target, ok := c.Request().Header["Hx-Target"]; ok {
		if target[0] == "list" {
			Render(c, http.StatusOK, views.Tasks(tasks))
			return Render(c, http.StatusOK, views.Counters(counts, 0, true))
		} else {
			return Render(c, http.StatusOK, views.Index(activeStatus, tasks, counts))
		}
	} else {
		return Render(c, http.StatusOK, views.Index(activeStatus, tasks, counts))
	}
}

func (s *Server) PostTaskHandler(c echo.Context) error {
	description := templ.EscapeString(c.FormValue("description"))
	if description == "" {
		c.Response().Header().Set("HX-Retarget", "#new-task")
		c.Response().Header().Set("HX-Reswap", "outerHTML")
		return Render(c, http.StatusUnprocessableEntity, views.InputForm(false, "يجرى إدخال وصف المهمة"))
	}

	task, err := s.Store.InsertTask(description)
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusBadRequest, views.ServerError())
	}

	total, completed, err := s.Store.GetTasksCounters()
	counts := todo.Counts{Total: total, Completed: completed}
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusBadRequest, views.ServerError())
	}

	c.Response().Status = http.StatusCreated

	Render(c, http.StatusCreated, views.Task(task))
	Render(c, http.StatusCreated, views.InputForm(true, ""))
	return Render(c, http.StatusCreated, views.Counters(counts, counts.Completed-1, true))
}

func (s *Server) ToggleStatusOfTaskHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusNotFound, views.ServerError())
	}

	counts, task, oldCompleted, err := s.Store.ToggleAndAnimationData(id)

	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusBadRequest, views.ServerError())
	}

	Render(c, http.StatusOK, views.Task(task))
	return Render(c, http.StatusOK, views.Counters(counts, oldCompleted, true))
}

func (s *Server) GetTaskFormHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusNotFound, views.ServerError())
	}

	task, err := s.Store.GetTaskById(id)
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusBadRequest, views.ServerError())
	}

	return Render(c, http.StatusOK, views.UpdateForm(task))
}

func (s *Server) UpdateTaskHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusNotFound, views.ServerError())
	}

	description := templ.EscapeString(c.FormValue("description"))
	if description == "" {
		log.Printf("%v\n", err)
		return Render(c, http.StatusUnprocessableEntity, views.ServerError())
	}

	task, err := s.Store.UpdateTask(id, description)
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusBadRequest, views.ServerError())
	}

	return Render(c, http.StatusOK, views.Task(task))
}

func (s *Server) DeleteTaskHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusNotFound, views.ServerError())
	}

	err = s.Store.DeleteTask(id)
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusBadRequest, views.ServerError())
	}

	total, completed, err := s.Store.GetTasksCounters()
	counts := todo.Counts{Total: total, Completed: completed}
	if err != nil {
		log.Printf("%v\n", err)
		return Render(c, http.StatusBadRequest, views.ServerError())
	}

	return Render(c, http.StatusOK, views.Counters(counts, counts.Completed+1, true))
}
