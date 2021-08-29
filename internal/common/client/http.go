package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pbTask "github.com/drifterz13/go-services/internal/common/genproto/task"
	"github.com/drifterz13/go-services/internal/common/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	r          *gin.Engine
	taskClient pbTask.TaskServiceClient
}

func NewServer(r *gin.Engine, taskClient pbTask.TaskServiceClient) *Server {
	return &Server{r: r, taskClient: taskClient}
}

func (s *Server) Serve() {
	s.registerTaskServer()
	s.r.Run(":8000")
}

func (s *Server) registerTaskServer() {
	s.r.GET("/tasks", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := s.taskClient.FindTasks(ctx, &emptypb.Empty{})
		if err != nil {
			log.Fatalf("error finding tasks: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"tasks": convertToTasksResponse(resp.Tasks),
		})
	})
}

func convertToTasksResponse(pbTasks []*pbTask.Task) []models.Task {
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
