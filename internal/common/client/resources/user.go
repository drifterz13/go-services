package resources

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	ce "github.com/drifterz13/go-services/internal/common/error"
	pbuser "github.com/drifterz13/go-services/internal/common/genproto/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userResources struct {
	client pbuser.UserServiceClient
}

func NewUserResources(client pbuser.UserServiceClient) *userResources {
	return &userResources{client: client}
}

func (rs userResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Post("/", rs.Create)

	log.Println("registered user resources.")

	return r
}

func (rs *userResources) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := rs.client.FindUsers(ctx, &emptypb.Empty{})
	if err != nil {
		render.Render(w, r, ce.ErrInternalServer(err))
		return
	}

	var users []User
	for _, u := range resp.Users {
		users = append(users, NewUserFromPb(u))
	}

	if err := render.Render(w, r, &UsersResponse{Users: users}); err != nil {
		render.Render(w, r, ce.ErrRender(err))
		return
	}
}

func (rs *userResources) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data := &CreateUserData{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ce.ErrBadRequest(err))
		return
	}

	_, err := rs.client.CreateUser(ctx, &pbuser.CreateUserRequest{Email: data.Email})
	if err != nil {
		render.Render(w, r, ce.ErrInternalServer(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

func NewUserFromPb(pbuser *pbuser.User) User {
	return User{
		ID:        pbuser.Id,
		Email:     pbuser.Email,
		CreatedAt: pbuser.CreatedAt.AsTime(),
		UpdatedAt: pbuser.UpdatedAt.AsTime(),
	}
}

func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *UsersResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type CreateUserData struct {
	Email string `json:"email"`
}

func (cr *CreateUserData) Bind(r *http.Request) error {
	if cr.Email == "" {
		return errors.New("email is required.")
	}

	return nil
}
