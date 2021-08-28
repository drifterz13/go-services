package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/drifterz13/go-services/internal/models"
	pb "github.com/drifterz13/go-services/proto/task"
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
	rows, err := r.conn.Query(ctx, "select * from tasks")
	if err != nil {
		return nil, err
	}

	var tasks []*pb.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}

		log.Printf("task title: %v", task.Title)
		log.Printf("task id: %v", task.ID)

		tasks = append(tasks, task.ToProto())
	}

	return tasks, nil
}

func (r *taskRepository) Create(ctx context.Context, title string) error {
	_, err := r.conn.Exec(ctx, "insert into tasks (title) values ($1)", title)

	return err
}

func registerDB() (*pgx.Conn, error) {
	var (
		dbhost   = os.Getenv("POSTGRES_HOST")
		dbport   = os.Getenv("POSTGRES_PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)

	log.Printf("host %v, port %v, user %v, password %v, db %v", dbhost, dbport, user, password, dbname)

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, dbhost, dbport, dbname))
	if err != nil {
		return nil, err
	}

	return conn, err
}
