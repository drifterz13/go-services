package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	pbtask "github.com/drifterz13/go-services/internal/common/genproto/task"
	pbuser "github.com/drifterz13/go-services/internal/common/genproto/user"
	"google.golang.org/grpc"
)

var (
	taskaddr = fmt.Sprintf("task:%s", os.Getenv("TASK_GRPC_PORT"))
	useraddr = fmt.Sprintf("user:%s", os.Getenv("USER_GRPC_PORT"))
)

type grpcConn struct {
	taskConn *grpc.ClientConn
	userConn *grpc.ClientConn
}

func NewGrpcConn() (*grpcConn, error) {
	taskConn, err := grpc.Dial(taskaddr, grpc.WithTimeout(5*time.Second), grpc.WithInsecure())
	if err != nil {
		return nil, errors.New("unable to connect to task gRPC")
	}

	userConn, err := grpc.Dial(useraddr, grpc.WithTimeout(5*time.Second), grpc.WithInsecure())
	if err != nil {
		return nil, errors.New("unable to connect to user gRPC")
	}

	return &grpcConn{taskConn, userConn}, err
}

func (c *grpcConn) GetTaskClient() pbtask.TaskServiceClient {
	return pbtask.NewTaskServiceClient(c.taskConn)
}

func (c *grpcConn) GetUserClient() pbuser.UserServiceClient {
	return pbuser.NewUserServiceClient(c.userConn)
}

func (c *grpcConn) Close() {
	c.taskConn.Close()
	c.userConn.Close()
}
