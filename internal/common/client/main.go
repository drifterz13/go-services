package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	taskClient, close, err := NewTaskClient()
	if err != nil {
		log.Fatalf("cannot connect to task client: %v\n", err)
	}
	defer close()

	srv := NewServer(r, taskClient)
	srv.Serve()
}
