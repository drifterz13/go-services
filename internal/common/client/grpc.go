package main

import (
	"log"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"google.golang.org/grpc"
)

const (
	addr = "localhost:50051"
)

func NewTaskClient() (pb.TaskServiceClient, func() error, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("error occur during connect to grpc: %v\n", err)

		return nil, func() error { return nil }, err
	}

	return pb.NewTaskServiceClient(conn), conn.Close, nil
}
