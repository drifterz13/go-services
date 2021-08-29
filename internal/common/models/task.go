package models

import (
	"time"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Member struct {
	ID   string `json:"id"`
	Role int    `json:"role"`
}

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Members   []Member  `json:"members"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Task) ToProto() *pb.Task {
	var members []*pb.Member

	if len(t.Members) > 0 {
		for _, member := range t.Members {
			members = append(members, &pb.Member{Id: member.ID, Role: pb.MemberRole(member.Role)})
		}
	}

	return &pb.Task{
		Id:        t.ID,
		Title:     t.Title,
		Status:    pb.Status(t.Status),
		Members:   members,
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}
}
