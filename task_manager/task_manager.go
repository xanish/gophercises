// Package task_manager provides a persistent task management system using BoltDB.
package task_manager

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"path/filepath"
	"time"
)

// Default database bucket name used for storing tasks.
const defaultDatabaseBucket = "taskManager"

// TaskManager manages tasks by persisting them in a BoltDB database.
type TaskManager struct {
	database *bolt.DB
	closer   func(db *bolt.DB) error
}

// NewTaskManager creates a new TaskManager instance by initializing
// and opening a BoltDB database connection.
func NewTaskManager() (TaskManager, error) {
	db, dbCloser, err := setupDatabase()
	if err != nil {
		return TaskManager{}, err
	}

	return TaskManager{database: db, closer: dbCloser}, nil
}

// Add persists a new Task to the database.
func (tm *TaskManager) Add(t Task) error {
	err := tm.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultDatabaseBucket))

		json, e := t.JSON()
		if e != nil {
			return fmt.Errorf("failed to encode json in writeable format: %w", e)
		}

		e = b.Put([]byte(t.Id()), json)
		if e != nil {
			return fmt.Errorf("failed to put key in bucket: %w", e)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create task %w", err)
	}

	return nil
}

// Complete marks a Task as completed in the database.
// It updates the Task's status to "StatusCompleted" and then calls the Add function
// to persist the updated Task in the database.
func (tm *TaskManager) Complete(t Task) error {
	t.Status = StatusCompleted

	return tm.Add(t)
}

// Delete removes a Task from the database.
// It attempts to delete the key corresponding to the Task's ID from the default bucket.
func (tm *TaskManager) Delete(t Task) error {
	err := tm.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultDatabaseBucket))

		e := b.Delete([]byte(t.Id()))
		if e != nil {
			return fmt.Errorf("failed to delete key in bucket: %w", e)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete task %w", err)
	}

	return nil
}

// Pending retrieves a list of tasks with a "StatusPending" status from the database.
func (tm *TaskManager) Pending() ([]Task, error) {
	tasks := make([]Task, 0)
	err := tm.database.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(defaultDatabaseBucket)).Cursor()
		prefix := []byte(time.Now().Format(time.DateOnly))

		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			data, e := NewTaskFromJSON(v)
			if e != nil {
				return fmt.Errorf("failed to fetch tasks %w", e)
			}
			if data.Status != StatusPending {
				continue
			}

			tasks = append(tasks, data)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks %w", err)
	}

	return tasks, nil
}

// Completed retrieves a list of tasks with a "StatusCompleted" status from the database.
func (tm *TaskManager) Completed() ([]Task, error) {
	tasks := make([]Task, 0)
	err := tm.database.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(defaultDatabaseBucket)).Cursor()
		prefix := []byte(time.Now().Format(time.DateOnly))

		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			data, e := NewTaskFromJSON(v)
			if e != nil {
				return fmt.Errorf("failed to fetch tasks %w", e)
			}
			if data.Status != StatusCompleted {
				continue
			}

			tasks = append(tasks, data)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks %w", err)
	}

	return tasks, nil
}

// Close closes the database handler for cleanup.
func (tm *TaskManager) Close() error {
	return tm.closer(tm.database)
}

// setupDatabase initializes and opens a BoltDB database connection.
// It creates the default bucket if it doesn't exist. Returns the opened database,
// a closure to close the database, and any encountered error.
func setupDatabase() (*bolt.DB, func(db *bolt.DB) error, error) {
	dbPath := filepath.Join(os.Getenv("DB_PATH"), "tasks.db")

	// Check if database file already exists
	newlyCreated, err := createDatabase(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open / create database: %w", err)
	}

	// Open the database connection
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create the default bucket if it doesn't exist
	if newlyCreated {
		err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(defaultDatabaseBucket))
			return err
		})

		if err != nil {
			return nil, nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Return the database connection, closer function, and any error
	return db, func(db *bolt.DB) error {
		return db.Close()
	}, nil
}

// createDatabase checks if the database file exists at the specified path (dbPath).
// If it doesn't exist, it attempts to create the file.
// It returns a boolean indicating whether the database file already existed (true)
// or was created successfully (false), and any encountered error.
func createDatabase(dbPath string) (bool, error) {
	_, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create the database file if it doesn't exist
			dbFile, err := os.Create(dbPath)
			if err != nil {
				return false, fmt.Errorf("failed to create database file: %w", err)
			}

			defer func(dbFile *os.File) {
				_ = dbFile.Close()
			}(dbFile)

			return true, nil
		} else {
			return false, fmt.Errorf("failed to stat database file: %w", err)
		}
	}

	return false, nil
}
