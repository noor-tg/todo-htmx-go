package store

import (
	"alnoor/todo-go-htmx"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	Path string
	DB   *sql.DB
}

func New(path string) SqliteStore {
	return SqliteStore{Path: path}
}

func (s *SqliteStore) Open(cleanup bool) error {
	if cleanup {
		os.Remove(s.Path)
	}
	db, _ := sql.Open("sqlite3", s.Path)
	s.DB = db

	return nil
}

func (s *SqliteStore) Migrate() error {
	prepared, _ := s.DB.Prepare(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
    		description text,
    		status TEXT CHECK(status IN ('مكتمل', 'مجدول')) NOT NULL DEFAULT 'مجدول'
		)
	`)

	prepared.Exec()

	return nil
}

func (s *SqliteStore) InsertTask(description string) (todo.Task, error) {
	prepared, _ := s.DB.Prepare(`
		INSERT INTO tasks (description) VALUES (?)
	`)

	result, _ := prepared.Exec(description)
	id, _ := result.LastInsertId()

	out, _ := s.GetTaskById(int(id))
	return out, nil
}

func (s *SqliteStore) GetTasks(filter todo.Task) ([]todo.Task, error) {
	query := `SELECT * FROM tasks`
	var queryArgs []interface{}

	FilterBy("status", "=", filter.Status, &query, &queryArgs)
	FilterBy("description", "LIKE", filter.Description, &query, &queryArgs)

	query += " ORDER BY id DESC"
	rows, err := s.DB.Query(query, queryArgs...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []todo.Task

	for rows.Next() {
		task := &todo.Task{}
		if err = rows.Scan(&task.Id, &task.Description, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (s *SqliteStore) UpdateTask(id int, description string) (todo.Task, error) {
	out, err := s.GetTaskById(id)
	if err != nil {
		log.Printf("could not update task: %v", err)
		return out, err
	}

	prepared, _ := s.DB.Prepare(`
		UPDATE tasks 
		SET description = ?
		WHERE id = ? 
	`)

	prepared.Exec(description, id)
	task, _ := s.GetTaskById(id)

	return task, nil
}

func (s *SqliteStore) GetTaskById(id int) (todo.Task, error) {
	single := s.DB.QueryRow("SELECT id, description, status FROM tasks WHERE id = ?", id)

	existing := &todo.Task{}

	err := single.Scan(&existing.Id, &existing.Description, &existing.Status)
	if err != nil {
		log.Printf("%v\n", err)
		return *existing, err
	}

	return *existing, nil
}

func (s *SqliteStore) DeleteTask(id int) error {
	_, err := s.GetTaskById(id)
	if err != nil {
		log.Printf("could not update task: %v", err)
		return err
	}

	prepared, _ := s.DB.Prepare(`
		DELETE FROM tasks 
		WHERE id = ? 
	`)

	prepared.Exec(id)

	return nil
}

func (s *SqliteStore) ToggleTaskStatus(id int) (todo.Task, error) {
	out, err := s.GetTaskById(id)
	if err != nil {
		log.Printf("could not update task: %v", err)
		return todo.Task{}, err
	}

	var status string
	if out.Status == "مكتمل" {
		status = "مجدول"
	} else {
		status = "مكتمل"
	}

	prepared, _ := s.DB.Prepare(`
		UPDATE tasks 
		SET status = ?
		WHERE id = ? 
	`)

	_, err = prepared.Exec(status, id)
	if err != nil {
		log.Printf("could not insert task: %v", err)
		return todo.Task{}, err
	}

	task, _ := s.GetTaskById(id)

	return task, nil
}

func (s *SqliteStore) GetTasksByStatus(status string) ([]todo.Task, error) {
	rows, _ := s.DB.Query(`SELECT * FROM tasks WHERE status = ?`, status)
	defer rows.Close()

	var tasks []todo.Task

	for rows.Next() {
		task := todo.Task{}
		if err := rows.Scan(&task.Id, &task.Description, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *SqliteStore) GetTasksCount() (int, error) {
	single := s.DB.QueryRow("SELECT count(*) as count FROM tasks")
	count := 0
	single.Scan(&count)
	return count, nil
}

func (s *SqliteStore) GetCompletedTasksCount() (int, error) {
	single := s.DB.QueryRow("SELECT count(*) as count FROM tasks where status = ?", "مكتمل")
	count := 0
	single.Scan(&count)
	return count, nil
}

func (s *SqliteStore) GetTasksCounters() (int, int, error) {
	total, _ := s.GetTasksCount()
	completed, _ := s.GetCompletedTasksCount()
	return total, completed, nil
}

func (s *SqliteStore) ToggleAndAnimationData(id int) (todo.Counts, todo.Task, int, error) {
	// get old task info
	old, err := s.GetTaskById(id)
	if err != nil {
		return todo.Counts{}, todo.Task{}, 0, err
	}

	// toggle task and get task info
	task, _ := s.ToggleTaskStatus(id)

	// calc complete increase or decrease
	encreaseCount := false
	if old.Status == "مكتمل" && task.Status == "مجدول" {
		encreaseCount = false
	}
	if old.Status == "مجدول" && task.Status == "مكتمل" {
		encreaseCount = true
	}

	// get tasks counts
	total, completed, err := s.GetTasksCounters()
	counts := todo.Counts{Total: total, Completed: completed}
	oldCompleted := 0
	if encreaseCount {
		oldCompleted = counts.Completed - 1
	} else {
		oldCompleted = counts.Completed + 1
	}

	// return data
	return counts, task, oldCompleted, nil
}
