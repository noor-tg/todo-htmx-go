package store

import (
	"alnoor/todo-go-htmx"
	"database/sql"
	"log"
	"os"

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
	queries, err := os.ReadFile("../db/sqlite/migrate/up.sql")
	if err != nil {
		log.Fatalf("could not migrate: %v", err)
		return err
	}

	prepared, err := s.DB.Prepare(string(queries))
	if err != nil {
		log.Fatalf("could not migrate: %v", err)
		return err
	}

	_, err = prepared.Exec()
	if err != nil {
		log.Fatalf("could not migrate: %v", err)
		return err
	}

	return nil
}

func (s *SqliteStore) InsertTask(task todo.Task) (todo.Task, error) {
	out := todo.Task{}

	prepared, err := s.DB.Prepare(`
		INSERT INTO tasks (description) VALUES (?)
	`)
	if err != nil {
		log.Fatalf("could not migrate: %v", err)
		return out, err
	}

	result, err := prepared.Exec(task.Description)
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("could not migrate: %v", err)
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
