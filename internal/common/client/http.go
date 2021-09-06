package main

import (
	"net/http"

	pbtask "github.com/drifterz13/go-services/internal/common/genproto/task"
	pbuser "github.com/drifterz13/go-services/internal/common/genproto/user"

	"github.com/drifterz13/go-services/internal/common/client/resources"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	taskClient pbtask.TaskServiceClient
	userClient pbuser.UserServiceClient
}

func NewServer(taskClient pbtask.TaskServiceClient, userClient pbuser.UserServiceClient) *Server {
	return &Server{taskClient, userClient}
}

func (s *Server) Serve() {
	addr := ":8000"
	r := chi.NewRouter()
	registerMiddleware(r)

	taskResources := resources.NewTaskResources(s.taskClient)
	r.Mount("/tasks", taskResources.Routes())

	userResources := resources.NewUserResources(s.userClient)
	r.Mount("/users", userResources.Routes())

	http.ListenAndServe(addr, r)
}

func registerMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
}
