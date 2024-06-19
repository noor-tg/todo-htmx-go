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
		task, err := store.InsertTask(input)
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

}
