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
		select 
			t.id, 
			t.title, 
			t.status, 
			coalesce(json_agg(json_build_object('id', u.id, 'role', tm.role)) filter (where u.id is not null), '[]') as members, 
			t.created_at, 
			t.updated_at
		from tasks t
		left join task_members tm on tm.task_id = t.id
		left join users u on tm.user_id = u.id
		group by t.id;
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
	_, err := r.conn.Exec(ctx, "insert into tasks (title) values ($1)", title)

	return err
}
