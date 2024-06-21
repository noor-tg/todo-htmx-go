package store_test

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/store"
	"testing"
)

func TestSqliteStore(t *testing.T) {
	store := store.SqliteStore{Path: "./todo.db"}
	t.Run("connect correctly to db", func(t *testing.T) {
		err := store.Open()

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
		input := todo.Task{
			Description: "مهمة 10",
		}

		tasks, err := store.GetTasks()
		if err != nil {
			t.Errorf("failed to get tasks: %v", err)
		}

		if tasks[len(tasks)-1].Description != input.Description {
			t.Errorf("failed to get tasks: %v", err)

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
			t.Errorf("failed update test: %v", err)
		}
		_, err = store.GetTaskById(task.Id)
		if err == nil {
			t.Errorf("failed update test: get task by id method should return error")
		}
	})

}
