package main

import (
	"context"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"github.com/drifterz13/go-services/internal/common/models"
	"github.com/jackc/pgx/v4"
)

type TaskRepository interface {
	Find(ctx context.Context) ([]*pb.Task, error)
	Create(ctx context.Context, title string) error
}

type taskRepository struct {
	conn *pgx.Conn
}

func NewTaskRepository(conn *pgx.Conn) TaskRepository {
	return &taskRepository{conn}
}

func (r *taskRepository) Find(ctx context.Context) ([]*pb.Task, error) {
	query := `
	SELECT 
		t.id, 
		t.title, 
		t.status, 
		COALESCE(json_agg(json_build_object('id', u.id, 'role', tm.role)) FILTER (WHERE u.id IS NOT NULL), '[]') AS members, 
		t.created_at, 
		t.updated_at
	FROM tasks t
	LEFT JOIN task_members tm ON tm.task_id = t.id
	LEFT JOIN users u ON tm.user_id = u.id
	GROUP BY t.id;
	`
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var tasks []*pb.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Status, &task.Members, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task.ToProto())
	}

	return tasks, nil
}

func (r *taskRepository) Create(ctx context.Context, title string) error {
	_, err := r.conn.Exec(ctx, "INSERT INTO tasks (title) VALUES ($1)", title)

	return err
}
