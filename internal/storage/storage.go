// Package storage implements a JSON-file-based storage for tasks.
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Status string

const (
	StatusActive Status = "active"
	StatusDone   Status = "done"
)

type Task struct {
	ID       int        `json:"id"`
	Text     string     `json:"text"`
	Priority Priority   `json:"priority"`
	Due      *time.Time `json:"due_at,omitempty"`
	Status   Status     `json:"status"`
	Created  time.Time  `json:"created_at"`
	Updated  time.Time  `json:"updated_at"`
}

type Store struct {
	LastID int    `json:"last_id"`
	Tasks  []Task `json:"tasks"`
}

func (s *Store) nextID() int {
	s.LastID++
	return s.LastID
}

func LoadStore(path string) (s *Store, err error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Store{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load store %s: %w", path, err)
	}
	defer (func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	})()

	decoder := json.NewDecoder(file)
	s = &Store{}
	err = decoder.Decode(s)
	if errors.Is(err, io.EOF) {
		return &Store{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load store %s: %w", path, err)
	}

	return s, nil
}

func (s *Store) SaveStore(path string) error {
	tmp := path + ".tmp"
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	if err := writeJSON(tmp, s); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

func writeJSON(path string, s *Store) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer (func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	})()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(s)
}

func (s *Store) AddTask(task Task) int {
	id := s.nextID()
	s.Tasks = append(s.Tasks, Task{
		ID:       id,
		Text:     task.Text,
		Priority: task.Priority,
		Due:      task.Due,
		Status:   StatusActive,
		Created:  time.Now(),
		Updated:  time.Now(),
	})

	return id
}

func (s *Store) DoneTask(id int) bool {
	for i := range s.Tasks {
		task := &s.Tasks[i]
		if task.ID == id {
			task.Status = StatusDone
			task.Updated = time.Now()
			return true
		}
	}

	return false
}

func (s *Store) RmTask(id int) bool {
	startLen := len(s.Tasks)
	s.Tasks = slices.DeleteFunc(s.Tasks, func(t Task) bool {
		return t.ID == id
	})

	return startLen != len(s.Tasks)
}

func (s *Store) RmTasks(filterFn FilterByFn) []Task {
	filtered, removed := FilterSlice(s.Tasks, filterFn)
	s.Tasks = filtered
	return removed
}
