package task_manager

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/xanish/gophercises/cli_task_manager/task"
)

func TestNewTaskManager(t *testing.T) {
	dbPath := tempFile()

	defer func(name string) {
		_ = os.Remove(name)
	}(dbPath)

	tm, err := NewTaskManager(dbPath)
	if err != nil {
		t.Errorf("failed to create task manager: %v", err)
	}

	if reflect.DeepEqual(tm, TaskManager{}) {
		t.Errorf("task manager should not be empty")
	}
}

func TestTaskManager_Functions(t *testing.T) {
	t.Run("should add task", func(t *testing.T) {
		tm := setupTestTaskManagerWithDummyTasks(t)
		_ = tm.Close()
	})

	t.Run("should delete task", func(t *testing.T) {
		tm := setupTestTaskManagerWithDummyTasks(t)

		taskId := 1
		assertNoError(t, tm.Delete(taskId))

		_ = tm.Close()
	})

	t.Run("should fetch all tasks", func(t *testing.T) {
		tm := setupTestTaskManagerWithDummyTasks(t)
		tasks, err := tm.List("")

		assertNoError(t, err)
		assertTasksLength(t, tasks, 5)
		_ = tm.Close()
	})

	t.Run("should fetch pending tasks", func(t *testing.T) {
		tm := setupTestTaskManagerWithDummyTasks(t)
		tasks, err := tm.List(task.StatusPending)

		assertNoError(t, err)
		assertTasksLength(t, tasks, 3)
		_ = tm.Close()
	})

	t.Run("should fetch completed tasks", func(t *testing.T) {
		tm := setupTestTaskManagerWithDummyTasks(t)
		tasks, err := tm.List(task.StatusCompleted)

		assertNoError(t, err)
		assertTasksLength(t, tasks, 2)
		_ = tm.Close()
	})

	t.Run("should mark task as completed and verify count before and after update", func(t *testing.T) {
		tm := setupTestTaskManagerWithDummyTasks(t)
		pending, err := tm.List(task.StatusPending)
		completed, err := tm.List(task.StatusCompleted)

		assertNoError(t, err)
		assertTasksLength(t, pending, 3)
		assertTasksLength(t, completed, 2)

		err = tm.Complete(pending[0].Id)
		assertNoError(t, err)

		pending, err = tm.List(task.StatusPending)
		completed, err = tm.List(task.StatusCompleted)
		assertNoError(t, err)
		assertTasksLength(t, pending, 2)
		assertTasksLength(t, completed, 3)

		_ = tm.Close()
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected err to be nil got: %v", err)
	}
}

func assertTasksLength(t *testing.T, tasks []task.Task, length int) {
	t.Helper()
	if len(tasks) != length {
		t.Errorf("expected %d task, got %v tasks", length, len(tasks))
	}
}

func setupTestTaskManagerWithDummyTasks(t *testing.T) *TaskManager {
	t.Helper()
	tm := setupTestTaskManager(t)

	for i := 1; i < 4; i++ {
		err := tm.Add(dummyTask(i, task.StatusPending))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}

	for i := 1; i < 3; i++ {
		err := tm.Add(dummyTask(i, task.StatusCompleted))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}

	return tm
}

func dummyTask(id int, status string) task.Task {
	return task.Task{
		Id:          id,
		Title:       "Dummy Task " + strconv.Itoa(id),
		Description: []string{"Description Line 1", "Description Line 2"},
		Status:      status,
	}
}

func setupTestTaskManager(t *testing.T) *TaskManager {
	db, closer := initTempDatabase(t)

	return &TaskManager{db: db, closer: closer}
}

func initTempDatabase(t *testing.T) (*bolt.DB, func(db *bolt.DB) error) {
	t.Helper()
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
