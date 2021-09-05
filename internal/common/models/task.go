package models

import (
	"time"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Member struct {
	ID   string `json:"_id" bson:"_id"`
	Role int    `json:"role" bson:"role"`
}

type Members = []Member

type Task struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	Status    int                `json:"status" bson:"status"`
	Members   Members            `json:"members" bson:"members"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UpdateTaskRequest struct {
	ID     string  `json:"_id"`
	Title  *string `json:"title,omitempty"`
	Status *int    `json:"status,omitempty"`
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
