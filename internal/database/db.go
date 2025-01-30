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
	Date     string
	Duration string
	Task     string
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
			date TEXT,
			duration TEXT,
			task TEXT
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
		INSERT INTO sessions(date, duration, task) VALUES (?, ?, ?)
	`
	stmt, err := d.db.Prepare(insertSQL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec(session.Date, session.Duration, session.Task)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
