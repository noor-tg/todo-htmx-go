package store

import (
	"alnoor/todo-go-htmx"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	Open() error
	Close() error
}

type SqliteStore struct {
	Path string
	DB   *sql.DB
}

func (s *SqliteStore) Open() error {
	db, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		log.Fatalf("could not connect to sqlite db: %v", err)
		return err
	}
	s.DB = db

	return nil
}

func (s *SqliteStore) Migrate() error {
	prepared, err := s.DB.Prepare(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
    		description text,
    		status TEXT CHECK(status IN ('مكتمل', 'مجدول')) NOT NULL DEFAULT 'مجدول'
		)
	`)
	if err != nil {
		log.Printf("could not migrate: %v", err)
		return err
	}

	_, err = prepared.Exec()
	if err != nil {
		log.Printf("could not migrate: %v", err)
		return err
	}

	return nil
}

func (s *SqliteStore) InsertTask(description string) (todo.Task, error) {
	out := todo.Task{}

	prepared, err := s.DB.Prepare(`
		INSERT INTO tasks (description) VALUES (?)
	`)
	if err != nil {
		log.Printf("could not insert task: %v", err)
		return out, err
	}

	result, err := prepared.Exec(description)
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("could not insert task: %v", err)
		return out, err
	}

	out, err = s.GetTaskById(int(id))
	if err != nil {
		log.Printf("could not insert task: %v", err)
		return out, err
	}

	return out, nil
}

func (s *SqliteStore) GetTasks(filters map[string]string) ([]todo.Task, error) {
	query := `SELECT * FROM tasks`
	var queryArgs []interface{}

	if filters != nil {
		query, queryArgs = FilterBy(
			filters, query, queryArgs, "status", "=",
		)
		query, queryArgs = FilterBy(
			filters, query, queryArgs, "description", "LIKE",
		)
	}

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

func FilterBy(filters map[string]string, query string, queryArgs []interface{}, key string, operator string) (string, []interface{}) {
	filter, filterExists := filters[key]
	if filterExists {
		if !ContainWhere(query) {
			query += " WHERE"
		} else {
			query += " OR"
		}
		if operator == "LIKE" {
			queryArgs = append(queryArgs, fmt.Sprintf("%%%s%%", filter))
		} else {
			queryArgs = append(queryArgs, filter)
		}
		query += fmt.Sprintf(" %s %s ?", key, operator)
	}
	return query, queryArgs
}

func ContainWhere(query string) bool {
	if strings.Contains(query, "WHERE") {
		return true
	}
	return false
}

func (s *SqliteStore) UpdateTask(id int, description string) (todo.Task, error) {
	out, err := s.GetTaskById(id)
	if err != nil {
		log.Printf("could not update task: %v", err)
		return out, err
	}

	prepared, err := s.DB.Prepare(`
		UPDATE tasks 
		SET description = ?
		WHERE id = ? 
	`)
	if err != nil {
		log.Printf("could not update task: %v", err)
		return out, err
	}

	_, err = prepared.Exec(description, id)

	if err != nil {
		log.Printf("could not insert task: %v", err)
		return out, err
	}

	task, err := s.GetTaskById(id)

	if err != nil {
		log.Printf("could not insert task: %v", err)
		return out, err
	}

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

	prepared, err := s.DB.Prepare(`
		DELETE FROM tasks 
		WHERE id = ? 
	`)
	if err != nil {
		log.Printf("could not update task: %v", err)
		return err
	}

	_, err = prepared.Exec(id)

	if err != nil {
		log.Printf("could not insert task: %v", err)
		return err
	}

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

	prepared, err := s.DB.Prepare(`
		UPDATE tasks 
		SET status = ?
		WHERE id = ? 
	`)
	if err != nil {
		log.Printf("could not update task: %v", err)
		return todo.Task{}, err
	}

	_, err = prepared.Exec(status, id)

	if err != nil {
		log.Printf("could not insert task: %v", err)
		return todo.Task{}, err
	}

	task, err := s.GetTaskById(id)

	if err != nil {
		log.Printf("could not insert task: %v", err)
		return todo.Task{}, err
	}

	return task, nil
}

func (s *SqliteStore) GetTasksByStatus(status string) ([]todo.Task, error) {
	rows, err := s.DB.Query(`SELECT * FROM tasks WHERE status = ?`, status)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []todo.Task

	for rows.Next() {
		task := todo.Task{}
		if err = rows.Scan(&task.Id, &task.Description, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
