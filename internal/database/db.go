package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

type Session struct {
	Duration     string
	Task         string
	Start        string
	ScheduledEnd string
	EndedAt      string
	Completed    int
}

func New() (*Database, error) {
	var isInitialized bool = true

	path := "./sqlite-database.db"
	if _, err := os.Stat(path); err != nil {
		isInitialized = false
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db := &Database{db: sqlDB}

	if !isInitialized {
		if err := db.createTable(); err != nil {
			db.Close()
			return nil, err
		}
	}
	return db, nil
}

func (d *Database) createTable() error {
	createSessionTableSQL := `
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY,
			duration TEXT,
			task TEXT,
			start TEXT,
			scheduled_end TEXT,
			ended_at TEXT,
			completed INTEGER
		);
	`
	statement, err := d.db.Prepare(createSessionTableSQL)
	if err != nil {
		return err
	}
	if _, err := statement.Exec(); err != nil {
		return err
	}
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AddSession(session Session) error {
	fmt.Println("INSERTING INTO DB")
	insertSQL := `
		INSERT INTO sessions(duration, task, start, scheduled_end, ended_at, completed) VALUES (?, ?, ?, ?, ?, ?)
	`
	stmt, err := d.db.Prepare(insertSQL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec(
		session.Duration,
		session.Task,
		session.Start,
		session.ScheduledEnd,
		session.EndedAt,
		session.Completed,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
