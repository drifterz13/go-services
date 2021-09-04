package main

import (
	"context"
	"time"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"github.com/drifterz13/go-services/internal/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	Find(ctx context.Context) ([]*pb.Task, error)
	Create(ctx context.Context, title string) error
}

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(c *mongo.Collection) TaskRepository {
	return &taskRepository{c}
}

func (r *taskRepository) Find(ctx context.Context) ([]*pb.Task, error) {
	var tasks []*pb.Task
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var task *models.Task
		err := cur.Decode(&task)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task.ToProto())
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) Create(ctx context.Context, title string) error {
	task := models.Task{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Members:   []models.Member{},
		Status:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := r.collection.InsertOne(ctx, &task)

	return err
}
