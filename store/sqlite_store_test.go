package store_test

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/store"
	"testing"

	"github.com/pioz/faker"
)

func TestSqliteStore(t *testing.T) {
	store := store.New("./todo.db")
	t.Run("connect correctly to db", func(t *testing.T) {
		err := store.Open(true)

		if err != nil {
			t.Errorf("error in test connect: %v", err)
		}
	})

	t.Run("migrate db", func(t *testing.T) {
		err := store.Migrate()
		if err != nil {
			t.Errorf("migration test failed: %v", err)
		}
	})

	t.Run("insert task", func(t *testing.T) {
		input := todo.Task{
			Description: "مهمة 10",
		}

		task, err := store.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed insert test: %v", err)
		}

		if input.Description != task.Description {
			t.Errorf("failed insert test: wanted %v got %v", input.Description, task.Description)
		}

		if task.Id == 0 {
			t.Errorf("failed insert test: got id %v", task.Id)
		}
	})

	t.Run("get all tasks", func(t *testing.T) {

		task, err := store.InsertTask(faker.ColorName())
		tasks, err := store.GetTasks(nil)
		if err != nil {
			t.Errorf("failed to get tasks: %v", err)
		}

		if len(tasks) > 0 && tasks[0].Description != task.Description {
			t.Errorf("failed to get tasks: %v", task.Description)
		}
	})
	t.Run("filter tasks by params", func(t *testing.T) {
		task1, err := store.InsertTask(faker.ColorName())
		if err != nil {
			t.Errorf("failed insert task 1: %v", err)
		}
		_, err = store.InsertTask(faker.ColorName())
		if err != nil {
			t.Errorf("failed insert task 2: %v", err)
		}
		task3, err := store.InsertTask(faker.ColorName())
		if err != nil {
			t.Errorf("failed insert task 3: %v", err)
		}

		tasks, err := store.GetTasks(map[string]string{"description": task3.Description})
		if err != nil {
			t.Errorf("failed to get tasks: %v", err)
		}

		exists := false
		notexist := true
		for _, task := range tasks {
			if task.Description == task3.Description {
				exists = true
				break
			}
			if task.Description == task1.Description {
				exists = false
				break
			}
		}
		if !exists {
			t.Errorf("task should exist: %v", task3.Description)
		}
		if !notexist {
			t.Errorf("task should not exist: %v", task1.Description)
		}
	})
	t.Run("update task", func(t *testing.T) {
		input := todo.Task{
			Description: "مهمة 10",
		}

		task, err := store.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed update test: %v", err)
		}

		description := "updated task"
		updatedTask, err := store.UpdateTask(task.Id, description)
		if description != updatedTask.Description {
			t.Errorf("failed update test: wanted %v got %v", description, task.Description)
		}

		if task.Id != updatedTask.Id {
			t.Errorf("failed update test: got id %v", task.Id)
		}
	})

	t.Run("get single task", func(t *testing.T) {
		input := todo.Task{
			Description: "مهمة 10",
		}

		task, err := store.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed insert test: %v", err)
		}
		existing, err := store.GetTaskById(task.Id)

		if err != nil {
			t.Errorf("failed insert test: %v", err)
		}
		if existing.Description != task.Description {
			t.Errorf("failed insert test: wanted %v got %v", existing.Description, task.Description)
		}
	})

	t.Run("delete task", func(t *testing.T) {
		input := todo.Task{
			Description: "مهمة 10",
		}

		task, err := store.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed update test: %v", err)
		}

		err = store.DeleteTask(task.Id)
		if err != nil {
			t.Errorf("failed delete test: %v", err)
		}
		_, err = store.GetTaskById(task.Id)
		if err == nil {
			t.Errorf("failed delete test: get task by id method should return error")
		}
	})

	t.Run("toggle task status", func(t *testing.T) {
		input := todo.Task{
			Description: "مهمة 10",
		}

		task, err := store.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed update test: %v", err)
		}

		toggled, err := store.ToggleTaskStatus(task.Id)
		if err != nil {
			t.Errorf("failed update test: get task by id method should return error")
		}

		if toggled.Status != "مكتمل" {
			t.Errorf("failed toggle test: status is not toggled")
		}
	})
	t.Run("get tasks by status", func(t *testing.T) {
		input := todo.Task{
			Description: "مهمة 10",
		}

		task, err := store.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed update test: %v", err)
		}

		_, err = store.ToggleTaskStatus(task.Id)
		if err != nil {
			t.Errorf("failed update test: get task by id method should return error")
		}

		tasks, err := store.GetTasksByStatus("مكتمل")
		if err != nil {
			t.Errorf("failed to get tasks: %v", err)
		}

		if tasks[len(tasks)-1].Description != input.Description {
			t.Errorf("failed to get tasks: %v", err)

		}
		if tasks[len(tasks)-1].Status != "مكتمل" {
			t.Errorf("failed toggle test: status is not toggled")
		}
	})

}
