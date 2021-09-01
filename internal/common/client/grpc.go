package main

import (
	"log"

	pbTask "github.com/drifterz13/go-services/internal/common/genproto/task"
	pbUser "github.com/drifterz13/go-services/internal/common/genproto/user"
	"google.golang.org/grpc"
)

const (
	addr = "localhost:50051"
)

func NewTaskClient() (pbTask.TaskServiceClient, func() error, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("error occur during connect to task gRPC: %v\n", err)

		return nil, func() error { return nil }, err
	}

	return pbTask.NewTaskServiceClient(conn), conn.Close, nil
}

func NewUserClient() (pbUser.UserServiceClient, func() error, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("error occur during connect to user gRPC: %v\n", err)

		return nil, func() error { return nil }, err
	}

	return pbUser.NewUserServiceClient(conn), conn.Close, nil
}
