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
	Id           int
	Duration     int
	TaskID       int
	Start        string
	ScheduledEnd string
	EndedAt      string
	Completed    int
}

type DailySession struct {
	Date              string
	TotalDuration     string
	CompletedSessions int
	TotalSessions     int
}

type Task struct {
	Id   int
	Name string
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
		if err := db.createTables(); err != nil {
			db.Close()
			return nil, err
		}
	}
	return db, nil
}

func (d *Database) createTables() error {
	createSessionsTableSQL := `
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY,
			duration TEXT,
			task_id INTEGER,
			start TEXT,
			scheduled_end TEXT,
			ended_at TEXT,
			completed INTEGER,
			FOREIGN KEY(task_id) REFERENCES tasks(id)
		);
	`
	createTasksTableSQL := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			name TEXT
		);
	`

	tables := []string{
		createSessionsTableSQL,
		createTasksTableSQL,
	}

	for _, sql := range tables {
		stmt, err := d.db.Prepare(sql)
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AddSession(session Session) error {
	fmt.Println("INSERTING INTO DB")
	insertSQL := `
		INSERT INTO sessions(duration, task_id, start, scheduled_end, ended_at, completed) VALUES (?, ?, ?, ?, ?, ?)
	`
	stmt, err := d.db.Prepare(insertSQL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		session.Duration,
		session.TaskID,
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

func (d *Database) CreateTask(task Task) error {
	insertSQL := `
		INSERT INTO tasks(name) VALUES (?)
	`

	stmt, err := d.db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Name)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteTask(name string) error {
	deleteSQL := `
		DELETE FROM tasks
		WHERE name = ?
	`
	stmt, err := d.db.Prepare(deleteSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) RetrieveTasks() ([]Task, error) {
	selectSQL := `
		SELECT name FROM tasks
	`

	row, err := d.db.Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var output []Task
	for row.Next() {
		var name string
		row.Scan(&name)
		output = append(output, Task{Name: name})
	}
	return output, nil
}

func (d *Database) RetrieveTaskByName(name string) (*Task, error) {
	selectSQL := `
		SELECT id, name FROM tasks WHERE name = ?
	`

	stmt, err := d.db.Prepare(selectSQL)
	if err != nil {
		return nil, err
	}

	var task Task
	row := stmt.QueryRow(name)
	err = row.Scan(&task.Id, &task.Name)
	if err != nil {
		return nil, err
	}
	return &task, err
}

func (d *Database) RetrieveDailySummary() ([]DailySession, error) {
	selectSQL := `
		SELECT
			DATE(ended_at) as date,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		GROUP BY DATE(ended_at)
		ORDER BY DATE(ended_at) DESC
	`

	stmt, err := d.db.Prepare(selectSQL)
	if err != nil {
		return nil, err
	}

	var output []DailySession
	row, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var dailySession DailySession
		err = row.Scan(
			&dailySession.Date,
			&dailySession.TotalDuration,
			&dailySession.CompletedSessions,
			&dailySession.TotalSessions,
		)
		if err != nil {
			return nil, err
		}

		output = append(output, dailySession)
	}

	return output, nil
}
