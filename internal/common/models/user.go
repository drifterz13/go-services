package models

import (
	"time"

	pb "github.com/drifterz13/go-services/internal/common/genproto/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
