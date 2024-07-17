package task

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"reflect"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	want := Task{
		Title:       "Task 1",
		Description: []string{"Description Line 1", "Description Line 2"},
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}
	got := NewTask("Task 1", []string{"Description Line 1", "Description Line 2"})

	t.Run("returns proper task", func(t *testing.T) {
		if !cmp.Equal(want, got, cmpopts.IgnoreFields(Task{}, "CreatedAt")) {
			t.Error(cmp.Diff(want, got))
		}
	})
}

func TestNewTaskFromJSON(t *testing.T) {
	want := Task{
		Title:       "Task 1",
		Description: []string{"Description Line 1", "Description Line 2"},
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}

	wantJSON, err := want.JSON()
	if err != nil {
		t.Fatalf("expected err to not be nil")
	}

	got, err := NewTaskFromJSON(wantJSON)
	if err != nil {
		t.Fatalf("expected err to not be nil")
	}

	t.Run("returns proper task from JSON", func(t *testing.T) {
		if !cmp.Equal(want, got, cmpopts.IgnoreFields(Task{}, "CreatedAt")) {
			t.Error(cmp.Diff(want, got))
		}
	})

	t.Run("returns proper task status", func(t *testing.T) {
		if got.Status != want.Status {
			t.Errorf("want status %s, got %s", want.Status, got.Status)
		}
	})

	t.Run("returns proper task id", func(t *testing.T) {
		if got.Id() != want.Id() {
			t.Errorf("want id %s, got %s", want.Id(), got.Id())
		}
	})

	t.Run("returns proper task json", func(t *testing.T) {
		gotJSON, err := got.JSON()
		if err != nil {
			t.Fatalf("expected err to not be nil")
		}

		wantJSON, err := want.JSON()
		if err != nil {
			t.Fatalf("expected err to not be nil")
		}

		if string(gotJSON) != string(wantJSON) {
			t.Errorf("want json %s, got %s", wantJSON, gotJSON)
		}
	})

	t.Run("should return error if json unmarshal fails", func(t *testing.T) {
		got, err = NewTaskFromJSON([]byte("faulty json"))
		if err == nil {
			t.Errorf("expected err to not be nil")
		}

		if !reflect.DeepEqual(Task{}, got) {
			t.Errorf("want task %+v, got %+v", Task{}, got)
		}
	})
}
