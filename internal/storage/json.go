package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gregorywalrath/go-todo/internal/todo"
)

// DefaultFileName is the JSON file where todos will be stored
const DefaultFileName = ".go-todo.json"

// LoadTodos reads todos from a JSON file in the user's home directory.
// If the file does not exist, it returns an empty slice.
func LoadTodos() ([]todo.Todo, error) {
	path, err := getFilePath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// Return empty slice if file does not exist
		return []todo.Todo{}, nil
	}

	// Read file contents
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var todos []todo.Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return todos, nil
}

// SaveTodos writes the slice of todos to the JSON file in the user's home directory
func SaveTodos(todos []todo.Todo) error {
	path, err := getFilePath()
	if err != nil {
		return err
	}

	// Marshal JSON with indentation for readability
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	// Write file
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	os.Rename(tmp, path)

	return nil
}

// getFilePath returns the full path to the JSON file in the user's home directory
func getFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, DefaultFileName), nil
}

// Example function to create a new Todo
func NewTodo(id int, title string) todo.Todo {
	return todo.Todo{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
		// CompletedAt remains nil until completed
	}
}
