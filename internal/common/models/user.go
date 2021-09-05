package models

import (
	"time"

	pb "github.com/drifterz13/go-services/internal/common/genproto/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"id"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:        u.ID.Hex(),
		Email:     u.Email,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
