package task_manager

import (
	"github.com/boltdb/bolt"
	"os"
	"testing"
)

func TestNewTaskManager(t *testing.T) {
	t.Run("should add task", func(t *testing.T) {
		_, tm := setupTaskManager(t)
		_ = tm.Close()
	})

	t.Run("should delete task", func(t *testing.T) {
		task, tm := setupTaskManager(t)
		err := tm.Delete(task)
		if err != nil {
			t.Errorf("error deleting task: %v", err)
		}

		_ = tm.Close()
	})

	t.Run("should mark task as completed and verify before and after update", func(t *testing.T) {
		task, tm := setupTaskManager(t)

		completed, err := tm.Completed()
		if err != nil {
			t.Errorf("error fetching completed tasks: %v", err)
		}

		if len(completed) != 0 {
			t.Errorf("expected completed tasks to be 0, got %v", len(completed))
		}

		err = tm.Complete(task)
		if err != nil {
			t.Errorf("error marking task as completed: %v", err)
		}

		completed, err = tm.Completed()
		if err != nil {
			t.Errorf("error fetching completed tasks: %v", err)
		}

		if len(completed) != 1 {
			t.Errorf("expected completed tasks to be 1, got %v", len(completed))
		}

		_ = tm.Close()
	})

	t.Run("should return pending tasks", func(t *testing.T) {
		_, tm := setupTaskManager(t)
		pending, err := tm.Pending()
		if err != nil {
			t.Errorf("error fetching pending tasks: %v", err)
		}

		if len(pending) != 1 {
			t.Errorf("expected pending tasks to be 1, got %v", len(pending))
		}

		_ = tm.Close()
	})
}

func setupTaskManager(t *testing.T) (Task, *TaskManager) {
	db, closer := setupTempDatabase(t)

	tm := TaskManager{database: db, closer: closer}
	task := NewTask("Task 1", []string{"Description Line 1", "Description Line 2"})

	err := tm.Add(task)
	if err != nil {
		t.Errorf("error adding task: %v", err)
	}

	return task, &tm
}

func setupTempDatabase(t *testing.T) (*bolt.DB, func(db *bolt.DB) error) {
	db, err := bolt.Open(tempFile(), 0600, nil)
	if err != nil {
		t.Fatal(err)
	} else if db == nil {
		t.Fatal("expected db")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(defaultDatabaseBucket))
		if err != nil {
			t.Fatalf("create bucket failed: %s", err)
		}

		return nil
	})

	if err != nil {
		t.Fatal("could not create bucket")
	}

	return db, func(db *bolt.DB) error {
		return db.Close()
	}
}

func tempFile() string {
	f, err := os.CreateTemp("", "bolt-*")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
