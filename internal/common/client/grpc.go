package main

import (
	"context"
	"time"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"github.com/drifterz13/go-services/internal/common/models"
	"go.uber.org/zap"
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
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	sugar := logger.Sugar()

	taskClient, close, err := NewTaskClient()
	if err != nil {
		sugar.Fatalf("cannot connect to task client: %v\n", err)
	}
	defer close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// if _, err := taskClient.CreateTask(ctx, &pb.CreateTaskRequest{Title: "Learn Docker"}); err != nil {
	// 	log.Fatalf("cannot create a task: %v\n", err)
	// }

	tasks, err := taskClient.FindTasks(ctx, &emptypb.Empty{})
	if err != nil {
		sugar.Fatalf("cannot find tasks: %v\n", err)
	}

	sugar.Info("tasks: ", tasks)

	// tasks_marshaled, err := protojson.Marshal(tasks)
	// if err != nil {
	// 	sugar.Fatalf("cannot marshaled tasks: %v\n", err)
	// }
	// tasks_marshaled, err := convertToTasksResponse(tasks.Tasks)
	// if err != nil {
	// 	sugar.Fatalf("cannot marshal tasks: %v\n", err)
	// }

	sugar.Info("received tasks: ", convertToTasksResponse(tasks.Tasks))
}

func convertToTasksResponse(pbTasks []*pb.Task) []models.Task {
	var tasks []models.Task

	for _, task := range pbTasks {
		var members []models.Member
		for _, member := range task.Members {
			members = append(members, models.Member{
				ID:   member.Id,
				Role: int(member.Role),
			})
		}

		tasks = append(tasks, models.Task{
			ID:        task.Id,
			Title:     task.Title,
			Status:    int(task.Status),
			Members:   members,
			CreatedAt: task.CreatedAt.AsTime(),
			UpdatedAt: task.UpdatedAt.AsTime(),
		})
	}

	return tasks
}
