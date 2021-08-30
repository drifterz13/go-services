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
	taskClient pbTask.TaskServiceClient
}

func NewServer(taskClient pbTask.TaskServiceClient) *Server {
	return &Server{taskClient: taskClient}
}

func (s *Server) Serve() {
	r := gin.Default()
	s.registerTaskServer(r)
	r.Run(":8000")
}

func (s *Server) registerTaskServer(r *gin.Engine) {
	r.GET("/tasks", func(c *gin.Context) {
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

	r.POST("/tasks", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var req pbTask.CreateTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid create task payload."})

			return
		}

		_, err := s.taskClient.CreateTask(ctx, &req)
		if err != nil {
			log.Fatalf("error finding tasks: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
}

func convertToTasksResponse(pbTasks []*pbTask.Task) []models.Task {
	var tasks []models.Task

	for _, task := range pbTasks {
		var members []models.Member = []models.Member{}

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
