package store

import (
	"alnoor/todo-go-htmx"
	"database/sql"
	"log"

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
    		description text 
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

	rows, err := s.DB.Query(`SELECT * FROM tasks WHERE id = ?`, id)

	defer rows.Close()
	rows.Next()

	if err = rows.Scan(&out.Id, &out.Description); err != nil {
		return out, err
	}

	return out, nil
}

func (s *SqliteStore) GetTasks() ([]todo.Task, error) {
	rows, err := s.DB.Query(`SELECT * FROM tasks`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []todo.Task

	for rows.Next() {
		task := todo.Task{}
		if err = rows.Scan(&task.Id, &task.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
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

	task := todo.Task{Id: id, Description: description}

	return task, nil
}

func (s *SqliteStore) GetTaskById(id int) (todo.Task, error) {
	single := s.DB.QueryRow("SELECT id, description FROM tasks WHERE id = ?", id)

	existing := &todo.Task{}

	err := single.Scan(&existing.Id, &existing.Description)

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
