package task

import (
	"context"
	"errors"
	"net/http"
	"time"

	ce "github.com/drifterz13/go-services/internal/common/error"
	pbtask "github.com/drifterz13/go-services/internal/common/genproto/task"
	"github.com/drifterz13/go-services/internal/common/models"
	"github.com/go-chi/render"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskHandler struct {
	client pbtask.TaskServiceClient
}

func NewTaskHander(client pbtask.TaskServiceClient) *TaskHandler {
	return &TaskHandler{client: client}
}

func (th *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := th.client.FindTasks(ctx, &emptypb.Empty{})
	if err != nil {
		render.Render(w, r, ce.ErrInternalServer(err))
		return
	}

	var tasks []Task
	for _, t := range resp.Tasks {
		tasks = append(tasks, NewTaskFromPb(t))
	}

	if err := render.Render(w, r, &TasksResponse{Tasks: tasks}); err != nil {
		render.Render(w, r, ce.ErrRender(err))
		return
	}
}

func (th *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data := &CreateTaskRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ce.ErrBadRequest(err))
		return
	}

	_, err := th.client.CreateTask(ctx, &pbtask.CreateTaskRequest{Title: data.Title})
	if err != nil {
		render.Render(w, r, ce.ErrInternalServer(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type Task struct {
	ID        string         `json:"_id"`
	Title     string         `json:"title"`
	Status    int            `json:"status"`
	Members   models.Members `json:"members,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type TaskResponse struct {
	Task Task `json:"task"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func NewTaskFromPb(task *pbtask.Task) Task {
	var members models.Members

	for _, member := range task.Members {
		members = append(members, models.Member{
			ID:   member.Id,
			Role: int(member.Role),
		})
	}

	return Task{
		ID:        task.Id,
		Title:     task.Title,
		Status:    int(task.Status),
		Members:   members,
		CreatedAt: task.CreatedAt.AsTime(),
		UpdatedAt: task.UpdatedAt.AsTime(),
	}
}

func (t *TasksResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (t *TaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func (ct *CreateTaskRequest) Bind(r *http.Request) error {
	if ct.Title == "" {
		return errors.New("title is required.")
	}

	return nil
}
