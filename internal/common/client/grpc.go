package main

import (
	"errors"
	"time"

	pbTask "github.com/drifterz13/go-services/internal/common/genproto/task"
	pbUser "github.com/drifterz13/go-services/internal/common/genproto/user"
	"google.golang.org/grpc"
)

const (
	addr1 = "localhost:50051"
	addr2 = "localhost:50053"
)

type grpcConn struct {
	taskConn *grpc.ClientConn
	userConn *grpc.ClientConn
}

func NewGrpcConn() (*grpcConn, error) {
	taskConn, err := grpc.Dial(addr1, grpc.WithTimeout(5*time.Second), grpc.WithInsecure())
	if err != nil {
		return nil, errors.New("unable to connect to task gRPC")
	}

	userConn, err := grpc.Dial(addr2, grpc.WithTimeout(5*time.Second), grpc.WithInsecure())
	if err != nil {
		return nil, errors.New("unable to connect to user gRPC")
	}

	return &grpcConn{taskConn, userConn}, err
}

func (c *grpcConn) GetTaskClient() pbTask.TaskServiceClient {
	return pbTask.NewTaskServiceClient(c.taskConn)
}

func (c *grpcConn) GetUserClient() pbUser.UserServiceClient {
	return pbUser.NewUserServiceClient(c.userConn)
}

func (c *grpcConn) Close() {
	c.taskConn.Close()
	c.userConn.Close()
}
