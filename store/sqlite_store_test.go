package store_test

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/store"
	"testing"

	"github.com/google/uuid"
	"github.com/pioz/faker"
)

func TestSqliteStore(t *testing.T) {
	sqlStore, _ := store.New("./todo.db", true)
	t.Run("connect correctly to db", func(t *testing.T) {
		_, err := store.New("./todo.db", true)

		if err != nil {
			t.Errorf("error in test connect: %v", err)
		}
	})

	t.Run("migrate db", func(t *testing.T) {
		err := sqlStore.Migrate()
		if err != nil {
			t.Errorf("migration test failed: %v", err)
		}
	})

	t.Run("insert task", func(t *testing.T) {
		input := todo.Task{
			Description: "مهمة 10",
		}

		task, err := sqlStore.InsertTask(input.Description)
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
	t.Run("insert task fail", func(t *testing.T) {
		input := todo.Task{
			Description: "",
		}

		task, err := sqlStore.InsertTask(input.Description)
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

		task, err := sqlStore.InsertTask(faker.ColorName())
		tasks, err := sqlStore.GetTasks(todo.Task{})
		if err != nil {
			t.Errorf("failed to get tasks: %v", err)
		}

		if len(tasks) > 0 && tasks[0].Description != task.Description {
			t.Errorf("failed to get tasks: %v", task.Description)
		}
	})
	t.Run("filter tasks by params", func(t *testing.T) {
		task1, err := sqlStore.InsertTask(faker.ColorName())
		if err != nil {
			t.Errorf("failed insert task 1: %v", err)
		}
		_, err = sqlStore.InsertTask(faker.ColorName())
		if err != nil {
			t.Errorf("failed insert task 2: %v", err)
		}
		task3, err := sqlStore.InsertTask(faker.ColorName())
		if err != nil {
			t.Errorf("failed insert task 3: %v", err)
		}

		tasks, err := sqlStore.GetTasks(task3)
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

		task, err := sqlStore.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed update test: %v", err)
		}

		description := "updated task"
		updatedTask, err := sqlStore.UpdateTask(task.Id, description)
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

		task, err := sqlStore.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed insert test: %v", err)
		}
		existing, err := sqlStore.GetTaskById(task.Id)

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

		task, err := sqlStore.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed update test: %v", err)
		}

		err = sqlStore.DeleteTask(task.Id)
		if err != nil {
			t.Errorf("failed delete test: %v", err)
		}
		_, err = sqlStore.GetTaskById(task.Id)
		if err == nil {
			t.Errorf("failed delete test: get task by id method should return error")
		}
	})

	t.Run("toggle task status", func(t *testing.T) {
		input := todo.Task{
			Description: "مهمة 10",
		}

		task, err := sqlStore.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed update test: %v", err)
		}

		toggled, err := sqlStore.ToggleTaskStatus(task.Id)
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

		task, err := sqlStore.InsertTask(input.Description)
		if err != nil {
			t.Errorf("failed update test: %v", err)
		}

		_, err = sqlStore.ToggleTaskStatus(task.Id)
		if err != nil {
			t.Errorf("failed update test: get task by id method should return error")
		}

		tasks, err := sqlStore.GetTasksByStatus("مكتمل")
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

	t.Run("get tasks count", func(t *testing.T) {
		db := uuid.NewString() + ".db"
		sqlStore, _ = store.New(db, true)
		sqlStore.Migrate()

		sqlStore.InsertTask(faker.ColorName())
		sqlStore.InsertTask(faker.ColorName())
		sqlStore.InsertTask(faker.ColorName())

		count, err := sqlStore.GetTasksCount()

		if err != nil {
			t.Errorf("failed to get tasks: %v", err)
		}

		if count != 3 {
			t.Errorf("count is not correct: got %v, wanted %v", count, 3)
		}
	})

	t.Run("get completed tasks count", func(t *testing.T) {
		db := uuid.NewString() + ".db"
		sqlStore, _ = store.New(db, true)
		sqlStore.Migrate()

		sqlStore.InsertTask(faker.ColorName())
		task, _ := sqlStore.InsertTask(faker.ColorName())
		sqlStore.ToggleTaskStatus(task.Id)

		count, err := sqlStore.GetCompletedTasksCount()

		if err != nil {
			t.Errorf("failed to get tasks: %v", err)
		}

		if count != 1 {
			t.Errorf("count is not correct: got %v, wanted %v", count, 1)
		}
	})
	t.Run("get tasks counters: total , completed", func(t *testing.T) {
		db := uuid.NewString() + ".db"
		sqlStore, _ = store.New(db, true)
		sqlStore.Migrate()

		sqlStore.InsertTask(faker.ColorName())
		task, _ := sqlStore.InsertTask(faker.ColorName())
		sqlStore.ToggleTaskStatus(task.Id)

		total, completed, err := sqlStore.GetTasksCounters()

		if err != nil {
			t.Errorf("failed to get tasks: %v", err)
		}

		if total != 2 {
			t.Errorf("total is not correct: got %v, wanted %v", total, 2)
		}
		if completed != 1 {
			t.Errorf("count is not correct: got %v, wanted %v", completed, 1)
		}
	})
	t.Run("get animation data for counters", func(t *testing.T) {
		db := uuid.NewString() + ".db"
		sqlStore, _ = store.New(db, true)
		sqlStore.Migrate()

		description := faker.ColorName()
		task, _ := sqlStore.InsertTask(description)
		counts, task, oldCompleted, err := sqlStore.ToggleAndAnimationData(task.Id)
		if err != nil {
			t.Errorf("error toggle and animate %v", err)
		}
		if counts.Completed != 1 {
			t.Errorf("completed != 1")
		}
		if oldCompleted != 0 {
			t.Errorf("oldCompleted should be 0 is %v", oldCompleted)
		}
		if counts.Total != 1 {
			t.Errorf("Total != 1")
		}
		if description != task.Description {
			t.Errorf("description is not correct sent %v, registered %v", description, task.Description)
		}
		counts, task, oldCompleted, err = sqlStore.ToggleAndAnimationData(task.Id)
		if oldCompleted != 1 {
			t.Errorf("oldCompleted should be 0 is %v", oldCompleted)
		}
		if counts.Completed != 0 {
			t.Errorf("completed != 0")
		}
	})

	t.Run("get animation data for counters", func(t *testing.T) {
		db := uuid.NewString() + ".db"
		sqlStore, _ = store.New(db, true)
		sqlStore.Migrate()

		_, _, _, err := sqlStore.ToggleAndAnimationData(1000)
		if err == nil {
			t.Errorf("should throw error %v", err)
		}
	})
}
