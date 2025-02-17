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

type SessionSummary struct {
	Date              string
	DateRange         string
	TotalDuration     int
	CompletedSessions int
	TotalSessions     int
}

type TaskSummary struct {
	TaskName          string
	TotalDuration     int
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

	insertTasksSQL := `
		INSERT INTO tasks(name) VALUES
			("study"),
			("read"),
			("focus"),
			("work");
	`

	tables := []string{
		createSessionsTableSQL,
		createTasksTableSQL,
		insertTasksSQL,
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

func (d *Database) RetrieveDailySummary() ([]SessionSummary, error) {
	query := `
		SELECT
			DATE(ended_at) as date,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		GROUP BY DATE(ended_at)
		ORDER BY DATE(ended_at) DESC
	`

	return d.querySession(query, func(rows *sql.Rows) (SessionSummary, error) {
		var s SessionSummary
		err := rows.Scan(
			&s.Date,
			&s.TotalDuration,
			&s.CompletedSessions,
			&s.TotalSessions,
		)
		return s, err
	})
}

func (d *Database) RetrieveWeeklySummary() ([]SessionSummary, error) {
	query := `
		SELECT
			STRFTIME('%Y-%W', ended_at) as date,
			STRFTIME('%Y-%m-%d', DATE(ended_at, 'weekday 0', '-6 days')) || ' to ' ||
			STRFTIME('%Y-%m-%d', DATE(ended_at, 'weekday 0')) as date_range,
			SUM(duration) as total_duration,	
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		GROUP BY date
		ORDER BY date DESC
	`

	return d.querySession(query, func(rows *sql.Rows) (SessionSummary, error) {
		var s SessionSummary
		err := rows.Scan(
			&s.Date,
			&s.DateRange,
			&s.TotalDuration,
			&s.CompletedSessions,
			&s.TotalSessions,
		)
		return s, err
	})
}

func (d *Database) RetrieveMonthlySummary() ([]SessionSummary, error) {
	query := `
		SELECT
			STRFTIME('%Y-%m', ended_at) as date,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		GROUP BY date
		ORDER BY date DESC
	`

	return d.querySession(query, func(rows *sql.Rows) (SessionSummary, error) {
		var s SessionSummary
		err := rows.Scan(
			&s.Date,
			&s.TotalDuration,
			&s.CompletedSessions,
			&s.TotalSessions,
		)
		return s, err
	})
}

func (d *Database) RetrieveYearlySummary() ([]SessionSummary, error) {
	query := `
		SELECT
			STRFTIME('%Y', ended_at) as date,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		GROUP BY date
		ORDER BY date DESC
	`

	return d.querySession(query, func(rows *sql.Rows) (SessionSummary, error) {
		var s SessionSummary
		err := rows.Scan(
			&s.Date,
			&s.TotalDuration,
			&s.CompletedSessions,
			&s.TotalSessions,
		)
		return s, err
	})
}

func (d *Database) RetrieveAlltimeSummary() ([]SessionSummary, error) {
	query := `
		SELECT
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
	`

	return d.querySession(query, func(rows *sql.Rows) (SessionSummary, error) {
		var s SessionSummary
		err := rows.Scan(
			&s.TotalDuration,
			&s.CompletedSessions,
			&s.TotalSessions,
		)
		return s, err
	})
}

func (d *Database) querySession(sql string, scanner func(*sql.Rows) (SessionSummary, error), args ...any) ([]SessionSummary, error) {
	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	var output []SessionSummary
	for rows.Next() {
		sessionSummary, err := scanner(rows)
		if err != nil {
			return nil, err
		}
		output = append(output, sessionSummary)
	}
	return output, nil
}

func (d *Database) RetrieveDailyTaskSummary(date string) ([]TaskSummary, error) {
	query := `
		SELECT
			tasks.name as task_name,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		LEFT JOIN tasks ON sessions.task_id = tasks.id
		WHERE DATE(ended_at) = ?
		GROUP BY task_id, task_name
	`

	return d.queryTask(query, date)
}

func (d *Database) RetrieveWeeklyTaskSummary(date string) ([]TaskSummary, error) {
	query := `
		SELECT
			tasks.name as task_name,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		LEFT JOIN tasks ON sessions.task_id = tasks.id
		WHERE STRFTIME('%Y-%W', ended_at) = ?
		GROUP BY task_id, task_name
	`

	return d.queryTask(query, date)
}

func (d *Database) RetrieveMonthlyTaskSummary(date string) ([]TaskSummary, error) {
	query := `
		SELECT
			tasks.name as task_name,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		LEFT JOIN tasks ON sessions.task_id = tasks.id
		WHERE STRFTIME('%Y-%m', ended_at) = ?
		GROUP BY task_id, task_name
	`

	return d.queryTask(query, date)
}

func (d *Database) RetrieveYearlyTaskSummary(date string) ([]TaskSummary, error) {
	query := `
		SELECT
			tasks.name as task_name,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		LEFT JOIN tasks ON sessions.task_id = tasks.id
		WHERE STRFTIME('%Y', ended_at) = ?
		GROUP BY task_id, task_name
	`

	return d.queryTask(query, date)
}

func (d *Database) RetrieveAlltimeTaskSummary() ([]TaskSummary, error) {
	query := `
		SELECT
			tasks.name as task_name,
			SUM(duration) as total_duration,
			SUM(completed) as completed_sessions,
			COUNT(*) as total_sessions
		FROM sessions
		LEFT JOIN tasks ON sessions.task_id = tasks.id
		GROUP BY task_id, task_name
	`

	return d.queryTask(query)
}

func (d *Database) queryTask(query string, args ...any) ([]TaskSummary, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	var output []TaskSummary
	for rows.Next() {
		var taskSummary TaskSummary
		err := rows.Scan(
			&taskSummary.TaskName,
			&taskSummary.TotalDuration,
			&taskSummary.CompletedSessions,
			&taskSummary.TotalSessions,
		)
		if err != nil {
			return nil, err
		}
		output = append(output, taskSummary)
	}
	return output, nil
}
