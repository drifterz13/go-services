package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	pbTask "github.com/drifterz13/go-services/internal/common/genproto/task"
	pbUser "github.com/drifterz13/go-services/internal/common/genproto/user"
	"github.com/drifterz13/go-services/internal/common/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	taskClient pbTask.TaskServiceClient
	userClient pbUser.UserServiceClient
}

func NewServer(taskClient pbTask.TaskServiceClient, userClient pbUser.UserServiceClient) *Server {
	return &Server{taskClient, userClient}
}

func (s *Server) Serve() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	s.registerTaskServer(r)
	log.Println("task server has registerd.")

	s.registerUserServer(r)
	log.Println("user server has registerd.")

	http.ListenAndServe(":8000", r)
}

func (s *Server) registerTaskServer(r chi.Router) {
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := s.taskClient.FindTasks(ctx, &emptypb.Empty{})
		if err != nil {
			respondError(w, err, http.StatusInternalServerError)
			return
		}

		var tasks TasksResponse
		for _, t := range resp.Tasks {
			tasks = append(tasks, NewTaskResponseFromPb(t))
		}

		b, err := json.Marshal(&tasks)
		if err != nil {
			respondError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	r.Post("/tasks", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var req pbTask.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, err, http.StatusBadRequest)
			return
		}

		if req.Title == "" {
			respondError(w, errors.New("task title is not provided."), http.StatusBadRequest)
			return
		}

		_, err := s.taskClient.CreateTask(ctx, &req)
		if err != nil {
			respondError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func (s *Server) registerUserServer(r chi.Router) {
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := s.userClient.FindUsers(ctx, &emptypb.Empty{})
		if err != nil {
			respondError(w, err, http.StatusInternalServerError)
			return
		}

		var users UsersResponse
		for _, u := range resp.Users {
			users = append(users, NewUserResponseFromPb(u))
		}

		b, err := json.Marshal(&users)
		if err != nil {
			respondError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var req pbUser.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, err, http.StatusBadRequest)
			return
		}

		_, err := s.userClient.CreateUser(ctx, &req)
		if err != nil {
			respondError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

type TaskResponse struct {
	ID        string         `json:"_id"`
	Title     string         `json:"title"`
	Status    int            `json:"status"`
	Members   models.Members `json:"members,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type TasksResponse = []TaskResponse

func NewTaskResponseFromPb(task *pbTask.Task) TaskResponse {
	var members models.Members

	for _, member := range task.Members {
		members = append(members, models.Member{
			ID:   member.Id,
			Role: int(member.Role),
		})
	}

	return TaskResponse{
		ID:        task.Id,
		Title:     task.Title,
		Status:    int(task.Status),
		Members:   members,
		CreatedAt: task.CreatedAt.AsTime(),
		UpdatedAt: task.UpdatedAt.AsTime(),
	}
}

func (tr *TaskResponse) ToJSON() ([]byte, error) {
	b, err := json.Marshal(tr)
	if err != nil {
		return nil, errors.New("cannot marshal task.")
	}

	return b, nil
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersResponse = []UserResponse

func NewUserResponseFromPb(pbUser *pbUser.User) UserResponse {
	return UserResponse{
		ID:        pbUser.Id,
		Email:     pbUser.Email,
		CreatedAt: pbUser.CreatedAt.AsTime(),
		UpdatedAt: pbUser.UpdatedAt.AsTime(),
	}
}

func respondError(w http.ResponseWriter, err error, code int) {
	var msg string = err.Error()

	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.Unavailable:
			msg = "service is unavailable"
		default:
			msg = err.Error()
		}
	}

	http.Error(w, msg, code)
}
