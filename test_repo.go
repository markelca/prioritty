package main

import (
	"fmt"
	"log"

	"github.com/markelca/prioritty/pkg/tasks"
)

func main() {
	repo, err := tasks.NewSQLiteRepository("data/test.db")
	if err != nil {
		log.Fatal("Failed to create repository:", err)
	}

	tasks := repo.FindAll()
	fmt.Printf("Found %d tasks:\n", len(tasks))

	for _, task := range tasks {
		fmt.Printf("ID: %d, Title: %s, Status: %d\n", task.Id, task.Title, task.Status)
	}
}

