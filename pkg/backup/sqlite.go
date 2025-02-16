package backup

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteConfig struct {
	Path string
	Key  string
}

type SQLiteChecker struct {
	db         *sql.DB
	repository RemoteFilesRepository
	cfg        SqliteConfig
	mu         sync.Mutex
	ignored    int
	uploaded   int
}

func NewSQLiteChecker(cfg SqliteConfig, repository RemoteFilesRepository) *SQLiteChecker {
	return &SQLiteChecker{
		repository: repository,
		cfg:        cfg,
	}
}

func (c *SQLiteChecker) Open(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.repository.Download(ctx, c.cfg.Key, c.cfg.Path)
	if err != nil {
		fmt.Printf("No existing database found or error downloading: %v\n", err)
	}

	db, err := sql.Open("sqlite3", c.cfg.Path)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	c.db = db

	_, err = c.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS files (
			path TEXT PRIMARY KEY,
			uploaded_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return nil
}

func (c *SQLiteChecker) Add(path string, uploadedAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.db.Exec(
		"INSERT or REPLACE INTO files (`path`, uploaded_at) VALUES (?, ?)",
		path,
		uploadedAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		fmt.Printf("Error adding file to database: %v\n", err)
		return
	}
	c.uploaded++
}

func (c *SQLiteChecker) Remove(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.db.Exec("DELETE FROM files WHERE path = ?", path)
	if err != nil {
		fmt.Printf("Error removing file from database: %v\n", err)
	}
}

func (c *SQLiteChecker) Exists(path string, lastUpdated time.Time) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	var storedTime time.Time
	var timeStr string
	err := c.db.QueryRow(
		"SELECT uploaded_at FROM files WHERE path = ?",
		path,
	).Scan(&timeStr)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		fmt.Printf("Error checking file existence: %v\n", err)
		return false
	}

	storedTime, err = time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		fmt.Printf("Error parsing stored time: %v\n", err)
		return false
	}

	if lastUpdated.After(storedTime) {
		return false
	}

	c.ignored++
	return true
}

func (c *SQLiteChecker) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.db == nil {
		return nil
	}

	err := c.repository.PutEditable(ctx, c.cfg.Path, c.cfg.Key)
	if err != nil {
		return fmt.Errorf("error uploading database: %w", err)
	}

	err = c.db.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}

	return nil
}

func (c *SQLiteChecker) Ignored() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ignored
}

func (c *SQLiteChecker) Uploaded() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.uploaded
}

func (c *SQLiteChecker) GetFiles() map[string]time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	files := make(map[string]time.Time)
	rows, err := c.db.Query("SELECT path, uploaded_at FROM files")
	if err != nil {
		fmt.Printf("Error getting files: %v\n", err)
		return files
	}
	defer rows.Close()

	for rows.Next() {
		var path, timeStr string
		err := rows.Scan(&path, &timeStr)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		uploadedAt, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			fmt.Printf("Error parsing time: %v\n", err)
			continue
		}

		files[path] = uploadedAt
	}

	return files
}
