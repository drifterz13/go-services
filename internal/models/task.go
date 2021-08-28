package models

import (
	"time"

	pb "github.com/drifterz13/go-services/proto/task"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Task) ToProto() *pb.Task {
	return &pb.Task{
		Id:        t.ID,
		Title:     t.Title,
		Status:    pb.Status(t.Status),
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}
}
