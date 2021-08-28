package main

import (
	"context"
	"log"
	"time"

	pb "github.com/drifterz13/go-services/proto/task"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	addr = "localhost:50051"
)

func NewTaskClient() (pb.TaskServiceClient, func() error, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return pb.NewTaskServiceClient(conn), conn.Close, nil
}

func main() {
	taskClient, close, err := NewTaskClient()
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := taskClient.CreateTask(ctx, &pb.CreateTaskRequest{Title: "Learn Docker"}); err != nil {
		log.Fatalf("cannot create a task: %v\n", err)
	}

	tasks, err := taskClient.FindTasks(ctx, &emptypb.Empty{})
	defer close()

	if err != nil {
		log.Fatalf("cannot find tasks: %v\n", err)
	}

	log.Printf("received tasks: %v\n", tasks)
}
