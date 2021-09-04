package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pbTask "github.com/drifterz13/go-services/internal/common/genproto/task"
	pbUser "github.com/drifterz13/go-services/internal/common/genproto/user"
	"github.com/drifterz13/go-services/internal/common/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	taskClient pbTask.TaskServiceClient
	userClient pbUser.UserServiceClient
}

func NewServer(taskClient pbTask.TaskServiceClient, userClient pbUser.UserServiceClient) *Server {
	return &Server{taskClient, userClient}
}

func (s *Server) Serve() {
	r := gin.Default()
	s.registerTaskServer(r)
	log.Println("task server has registerd.")

	s.registerUserServer(r)
	log.Println("user server has registerd.")

	r.Run(":8000")
}

func (s *Server) registerTaskServer(r *gin.Engine) {
	r.GET("/tasks", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := s.taskClient.FindTasks(ctx, &emptypb.Empty{})
		if err != nil {
			respondError(c, err)

			return
		}

		tasks := MarshalTask(resp.Tasks)
		c.JSON(http.StatusOK, gin.H{
			"tasks": tasks,
		})
	})

	r.POST("/tasks", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var req pbTask.CreateTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid create task payload.", "err": err.Error()})

			return
		}

		_, err := s.taskClient.CreateTask(ctx, &req)
		if err != nil {
			respondError(c, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
}

func (s *Server) registerUserServer(r *gin.Engine) {
	r.GET("/users", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := s.userClient.FindUsers(ctx, &emptypb.Empty{})
		if err != nil {
			respondError(c, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"users": MarshalUser(resp.Users)})
	})

	r.POST("/users", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var req pbUser.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid create user payload.", "err": err.Error()})

			return
		}

		_, err := s.userClient.CreateUser(ctx, &req)
		if err != nil {
			respondError(c, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
}

func MarshalUser(pbUsers []*pbUser.User) []models.User {
	var users []models.User = []models.User{}
	for _, user := range pbUsers {
		users = append(users, models.User{
			ID:        user.Id,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.AsTime(),
			UpdatedAt: user.UpdatedAt.AsTime(),
		})
	}

	return users
}

func MarshalTask(pbTasks []*pbTask.Task) []models.TaskResponse {
	log.Printf("pb tasks: %v\n", pbTasks)
	var tasks []models.TaskResponse = []models.TaskResponse{}

	for _, task := range pbTasks {
		var members []models.Member = []models.Member{}

		for _, member := range task.Members {
			members = append(members, models.Member{
				ID:   member.Id,
				Role: int(member.Role),
			})
		}

		tasks = append(tasks, models.TaskResponse{
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

func respondError(c *gin.Context, err error) {
	var msg string = err.Error()

	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.Unavailable:
			msg = "service is unavailable"
		default:
			msg = err.Error()
		}
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": msg,
	})
}
