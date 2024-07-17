// Package task_manager provides a persistent task management system using BoltDB.
package task_manager

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/xanish/gophercises/cli_task_manager/task"
)

// Default database bucket name used for storing tasks.
const defaultDatabaseBucket = "taskManager"

// TaskManager manages tasks by persisting them in a BoltDB database.
type TaskManager struct {
	db     *bolt.DB
	closer func(db *bolt.DB) error
}

// NewTaskManager creates a new TaskManager instance by initializing
// and opening a BoltDB database connection.
func NewTaskManager(dbPath string) (TaskManager, error) {
	db, dbCloser, err := initDatabase(dbPath)
	if err != nil {
		return TaskManager{}, err
	}

	return TaskManager{db: db, closer: dbCloser}, nil
}

// Add persists a new Task to the database.
func (tm *TaskManager) Add(t task.Task) error {
	return tm.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultDatabaseBucket))

		id64, _ := b.NextSequence()
		t.Id = int(id64)

		json, e := t.JSON()
		if e != nil {
			return fmt.Errorf("failed to encode task: %w", e)
		}

		e = b.Put(itob(t.Id), json)
		if e != nil {
			return fmt.Errorf("failed to create task: %w", e)
		}

		return nil
	})
}

// Complete marks a Task as completed in the database.
// It fetches the Task and updates its status to "StatusCompleted" before persisting it back in database.
func (tm *TaskManager) Complete(id int) error {
	return tm.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultDatabaseBucket))
		t, err := task.NewTaskFromJSON(b.Get(itob(id)))
		if err != nil {
			return fmt.Errorf("failed to fetch task with id %d: %w", id, err)
		}

		t.Status = task.StatusCompleted

		json, e := t.JSON()
		if e != nil {
			return fmt.Errorf("failed to encode task: %w", e)
		}

		e = b.Put(itob(t.Id), json)
		if e != nil {
			return fmt.Errorf("failed to update task: %w", e)
		}

		return nil
	})
}

// Delete removes a Task from the database.
// It attempts to delete the key corresponding to the Task's ID from the default bucket.
func (tm *TaskManager) Delete(id int) error {
	return tm.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultDatabaseBucket))

		e := b.Delete(itob(id))
		if e != nil {
			return fmt.Errorf("failed to delete task with id %d: %w", id, e)
		}

		return nil
	})
}

// List retrieves all tasks that match the given status from database.
func (tm *TaskManager) List(status string) ([]task.Task, error) {
	tasks := make([]task.Task, 0)

	err := tm.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(defaultDatabaseBucket)).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			t, e := task.NewTaskFromJSON(v)
			if e != nil {
				return fmt.Errorf("failed to fetch tasks with status %s: %w", status, e)
			}

			if status != "" && t.Status != strings.ToLower(status) {
				continue
			}

			tasks = append(tasks, t)
		}

		return nil
	})

	return tasks, err
}

// Close closes the database handler for cleanup.
func (tm *TaskManager) Close() error {
	return tm.closer(tm.db)
}

// initDatabase initializes and opens a BoltDB database connection.
// It creates the default bucket if it doesn't exist. Returns the opened database,
// a closure to close the database, and any encountered error.
func initDatabase(dbPath string) (*bolt.DB, func(db *bolt.DB) error, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultDatabaseBucket))
		return err
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialise data bucket: %w", err)
	}

	return db, func(db *bolt.DB) error {
		return db.Close()
	}, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
