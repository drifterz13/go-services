package resources

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	ce "github.com/drifterz13/go-services/internal/common/error"
	pbtask "github.com/drifterz13/go-services/internal/common/genproto/task"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"google.golang.org/protobuf/types/known/emptypb"
)

type taskResources struct {
	client pbtask.TaskServiceClient
}

func NewTaskResources(client pbtask.TaskServiceClient) *taskResources {
	return &taskResources{client}
}

func (rs taskResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Post("/", rs.Create)

	r.Route("/{id}", func(r chi.Router) {
		r.Patch("/", rs.Update)
		r.Delete("/", rs.Delete)
	})

	log.Println("registered task resources.")

	return r
}

func (rs *taskResources) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := rs.client.FindTasks(ctx, &emptypb.Empty{})
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

func (rs *taskResources) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data := &CreateTaskData{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ce.ErrBadRequest(err))
		return
	}

	_, err := rs.client.CreateTask(ctx, &pbtask.CreateTaskRequest{Title: data.Title})
	if err != nil {
		render.Render(w, r, ce.ErrInternalServer(err))
		return
	}

	render.NoContent(w, r)
}

func (rs *taskResources) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var taskId string
	if taskId = chi.URLParam(r, "id"); taskId == "" {
		render.Render(w, r, ce.ErrNotFound)
		return
	}

	data := &UpdateTaskData{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ce.ErrBadRequest(err))
		return
	}

	req := &pbtask.UpdateTaskRequest{TaskId: taskId}
	if data.Title != nil {
		req.Title = *data.Title
	}
	if data.Status != nil {
		req.Status = pbtask.Status(*data.Status)
	}

	_, err := rs.client.UpdateTask(ctx, req)
	if err != nil {
		render.Render(w, r, ce.ErrNotFound)
		return
	}

	render.NoContent(w, r)
}

func (rs *taskResources) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var taskId string
	if taskId = chi.URLParam(r, "id"); taskId == "" {
		render.Render(w, r, ce.ErrNotFound)
		return
	}

	_, err := rs.client.DeleteTask(ctx, &pbtask.DeleteTaskRequest{TaskId: taskId})
	if err != nil {
		render.Render(w, r, ce.ErrBadRequest(err))
		return
	}

	render.NoContent(w, r)
}

type Member struct {
	ID   string `json:"_id" bson:"_id"`
	Role int    `json:"role" bson:"role"`
}

type Members = []Member

type Task struct {
	ID        string    `json:"_id"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Members   Members   `json:"members,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskResponse struct {
	Task Task `json:"task"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func NewTaskFromPb(task *pbtask.Task) Task {
	var members Members

	for _, member := range task.Members {
		members = append(members, Member{
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

type CreateTaskData struct {
	Title string `json:"title"`
}

func (ct *CreateTaskData) Bind(r *http.Request) error {
	if ct.Title == "" {
		return errors.New("title is required.")
	}

	return nil
}

type UpdateTaskData struct {
	Title  *string `json:"title,omitempty"`
	Status *int    `json:"status,omitempty"`
}

func (ut *UpdateTaskData) Bind(r *http.Request) error {
	if ut.Title == nil && ut.Status == nil {
		return errors.New("update fields are not specified.")
	}

	return nil
}
