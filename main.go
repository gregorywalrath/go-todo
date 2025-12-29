package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"time"

	"github.com/gregorywalrath/go-todo/internal/storage"
	"github.com/gregorywalrath/go-todo/internal/todo"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd(os.Args[2:])
	case "list":
		handleList()
	case "complete":
		handleComplete(os.Args[2:])
	case "delete":
		handleDelete(os.Args[2:])
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go-todo add <task title>   Add a new todo")
	fmt.Println("  go-todo list               List all todos")
	fmt.Println("  go-todo complete <todo ID> Mark a todo as completed")
	fmt.Println("  go-todo delete <todo ID>   Delete a todo")
	fmt.Println("  go-todo help               Show this help message")
}

func handleAdd(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: missing task title")
		printUsage()
		return
	}

	title := strings.TrimSpace(strings.Join(args, " "))

	todos, err := storage.LoadTodos()
	if err != nil {
		log.Fatalf("Failed to load todos: %v", err)
	}

	maxID := 0
	for _, t := range todos {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	newID := maxID + 1

	newTodo := storage.NewTodo(newID, title)
	todos = append(todos, newTodo)

	if err := storage.SaveTodos(todos); err != nil {
		log.Fatalf("Error saving todos: %v", err)
	}

	fmt.Printf("Added todo: %s\n", title)
}

func handleList() {
	todos, err := storage.LoadTodos()
	if err != nil {
		log.Fatalf("Error loading todos: %v", err)
	}

	if len(todos) == 0 {
		fmt.Println("No todos found!")
		return
	}

	fmt.Println("Todos:")
	for _, t := range todos {
		status := " "
		if t.Completed {
			status = "X"
		}
		fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Title)
	}
}

func handleComplete(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: missing todo ID")
		printUsage()
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: todo ID must be a number")
		return
	}

	todos, err := storage.LoadTodos()
	if err != nil {
		log.Fatalf("Error loading todos: %v", err)
	}

	found := false
	for i, t := range todos {
		if t.ID == id {
			if t.Completed {
				fmt.Printf("Todo %d is already completed: %s\n", t.ID, t.Title)
				return
			}

			now := time.Now()
			todos[i].Completed = true
			todos[i].CompletedAt = &now
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Todo with ID %d not found\n", id)
		return
	}

	if err := storage.SaveTodos(todos); err != nil {
		log.Fatalf("Error saving todos: %v", err)
	}

	fmt.Printf("Todo %d marked as completed\n", id)
}

func handleDelete(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: missing todo ID")
		printUsage()
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: todo ID must be a number")
		return
	}

	todos, err := storage.LoadTodos()
	if err != nil {
		log.Fatalf("Error loading todos: %v", err)
	}

	found := false
	newTodos := make([]todo.Todo, 0, len(todos))
	for _, t := range todos {
		if t.ID == id {
			found = true
			continue
		}
		newTodos = append(newTodos, t)
	}

	if !found {
		fmt.Printf("Todo with ID %d not found\n", id)
		return
	}

	if err := storage.SaveTodos(newTodos); err != nil {
		log.Fatalf("Error saving todos: %v", err)
	}

	fmt.Printf("Todo %d deleted successfully\n", id)
}
