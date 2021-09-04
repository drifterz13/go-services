package models

import (
	"time"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Member struct {
	ID   string `bson:"id"`
	Role int    `bson:"role"`
}

type Task struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	Status    int                `json:"status" bson:"status"`
	Members   []Member           `json:"members" bson:"members"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type TaskResponse struct {
	ID        string    `json:"_id"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Members   []Member  `json:"members"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Task) ToProto() *pb.Task {
	var members []*pb.Member = []*pb.Member{}

	for _, member := range t.Members {
		members = append(members, &pb.Member{Id: member.ID, Role: pb.MemberRole(member.Role)})
	}

	return &pb.Task{
		Id:        t.ID.Hex(),
		Title:     t.Title,
		Status:    pb.Status(t.Status),
		Members:   members,
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}
}
