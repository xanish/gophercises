// Package task_manager provides a simple task management data structure
// and helper functions for creating and working with tasks.
package task_manager

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
)

// Task represents a single task with its title, description, status,
// and creation timestamp.
type Task struct {
	Title       string    `json:"Title"`
	Description []string  `json:"Description"`
	Status      string    `json:"Status"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

// NewTask creates a new Task instance with the provided title and description.
// The status is automatically set to "Pending" and the CreatedAt timestamp
// is set to the current time.
func NewTask(title string, description []string) Task {
	return Task{
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}
}

// NewTaskFromJSON attempts to unmarshal a Task object from a JSON byte slice.
// It returns the un-marshaled Task and any error encountered during the process.
func NewTaskFromJSON(bytes []byte) (Task, error) {
	task := Task{}

	err := json.Unmarshal(bytes, &task)
	if err != nil {
		return Task{}, err
	}

	return task, err
}

// Id generates a unique identifier for the Task based on its creation timestamp
// and title (formatted with underscores instead of spaces for consistency).
func (t Task) Id() string {
	return t.CreatedAt.Format(strings.ReplaceAll(time.DateTime, " ", "_")) + "_" + strings.ReplaceAll(t.Title, " ", "_")
}

// JSON marshals the Task object into a JSON byte slice.
// It returns the byte slice and any error encountered during marshalling.
func (t Task) JSON() ([]byte, error) {
	return json.Marshal(t)
}
