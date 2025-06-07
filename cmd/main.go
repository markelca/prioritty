package main

import (
	"log"

	"github.com/markelca/prioritty/internal/cli"
	"github.com/markelca/prioritty/internal/logger"
)

func main() {
	defer logger.ShutdownLogger()
	if err := cli.Execute(); err != nil {
		log.Fatalf("Command failed: %v", err)
	}
}
