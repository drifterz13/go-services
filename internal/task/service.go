package main

import (
	"context"
	"fmt"

	pb "github.com/drifterz13/go-services/proto/task"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type taskService struct {
	pb.UnimplementedTaskServiceServer
}

func (s *taskService) FindTasks(ctx context.Context, in *emptypb.Empty) (*pb.FindTasksResponse, error) {
	var tasks []*pb.Task
	for _, title := range []string{"task1", "task2"} {
		task := &pb.Task{
			Id:        uuid.NewString(),
			Title:     title,
			Status:    pb.Status_ACTIVE,
			Members:   []*pb.Member{},
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		}

		tasks = append(tasks, task)
	}

	return &pb.FindTasksResponse{Tasks: tasks}, nil
}

func (s *taskService) FindTask(ctx context.Context, in *pb.FindTaskRequest) (*pb.FindTaskResponse, error) {
	return nil, nil
}

func (s *taskService) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	fmt.Printf("debug: task title is: %v\n", in.GetTitle())

	task := &pb.Task{
		Id:        uuid.NewString(),
		Title:     in.GetTitle(),
		Status:    pb.Status_ACTIVE,
		Members:   []*pb.Member{},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	return &pb.CreateTaskResponse{Task: task}, nil
}

func (s *taskService) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	return nil, nil
}
