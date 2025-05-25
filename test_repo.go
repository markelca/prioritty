package main

import (
	"log"

	"github.com/markelca/prioritty/pkg/tasks"
)

func main() {
	repo, err := tasks.NewSQLiteRepository("/home/markel/.config/prioritty/prioritty.db")
	if err != nil {
		log.Fatal("Failed to create repository:", err)
	}

	err = repo.UpdateStatus(tasks.Task{Id: 1}, tasks.InProgress)
	if err != nil {
		log.Fatal("Failed to update task:", err)
	}
}
