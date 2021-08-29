package main

import (
	"log"
)

func main() {
	taskClient, close, err := NewTaskClient()
	if err != nil {
		log.Fatalf("cannot connect to task client: %v\n", err)
	}
	defer close()

	srv := NewServer(taskClient)
	srv.Serve()
}
