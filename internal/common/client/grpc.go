package main

import (
	"errors"

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
	taskConn, err := grpc.Dial(addr1, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.New("unable to connect to task gRPC")
	}

	userConn, err := grpc.Dial(addr2, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.New("unable to connect to user gRPC")
	}

	return &grpcConn{taskConn, userConn}, nil
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
