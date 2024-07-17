package task_manager

import (
	"github.com/boltdb/bolt"
	"github.com/xanish/gophercises/task_manager/task"
	"os"
	"reflect"
	"testing"
)

func TestNewTaskManager(t *testing.T) {
	t.Run("should create new task manager", func(t *testing.T) {
		dbPath := tempFile()

		defer func(name string) {
			_ = os.Remove(name)
		}(dbPath)

		tm, err := NewTaskManager(dbPath)
		if err != nil {
			t.Errorf("error creating task manager: %v", err)
		}

		if reflect.DeepEqual(tm, TaskManager{}) {
			t.Errorf("task manager should not be empty")
		}
	})

	t.Run("should add task", func(t *testing.T) {
		_, tm := setupTaskManager(t)
		_ = tm.Close()
	})

	t.Run("should delete task", func(t *testing.T) {
		tsk, tm := setupTaskManager(t)
		err := tm.Delete(tsk)
		if err != nil {
			t.Errorf("error deleting task: %v", err)
		}

		_ = tm.Close()
	})

	t.Run("should mark task as completed and verify before and after update", func(t *testing.T) {
		tsk, tm := setupTaskManager(t)

		completed, err := tm.Completed()
		if err != nil {
			t.Errorf("error fetching completed tasks: %v", err)
		}

		if len(completed) != 0 {
			t.Errorf("expected completed tasks to be 0, got %v", len(completed))
		}

		err = tm.Complete(tsk)
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

		tsk := task.NewTask("Title 2", []string{"Description"})
		tsk.Status = task.StatusCompleted
		err := tm.Add(tsk)
		if err != nil {
			t.Errorf("error adding task: %v", err)
		}

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

func setupTaskManager(t *testing.T) (task.Task, *TaskManager) {
	db, closer := setupTempDatabase(t)

	tm := TaskManager{database: db, closer: closer}
	tsk := task.NewTask("Task 1", []string{"Description Line 1", "Description Line 2"})

	err := tm.Add(tsk)
	if err != nil {
		t.Errorf("error adding task: %v", err)
	}

	return tsk, &tm
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
