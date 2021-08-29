package main

import (
	"context"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type taskService struct {
	pb.UnimplementedTaskServiceServer
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) pb.TaskServiceServer {
	return &taskService{repo: repo}
}

func (s *taskService) FindTasks(ctx context.Context, in *emptypb.Empty) (*pb.FindTasksResponse, error) {
	tasks, err := s.repo.Find(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.FindTasksResponse{Tasks: tasks}, nil
}

func (s *taskService) FindTask(ctx context.Context, in *pb.FindTaskRequest) (*pb.FindTaskResponse, error) {
	return nil, nil
}

func (s *taskService) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	err := s.repo.Create(context.Background(), in.Title)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateTaskResponse{}, nil
}

func (s *taskService) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	return nil, nil
}
