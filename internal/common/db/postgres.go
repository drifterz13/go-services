package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var (
	dbhost   = os.Getenv("POSTGRES_HOST")
	dbport   = os.Getenv("POSTGRES_PORT")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

type PostgresRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type PostgresDBRepository struct {
	conn *pgx.Conn
}

func NewPostgresDBRepository(conn *pgx.Conn) *PostgresDBRepository {
	return &PostgresDBRepository{conn}
}

func (r *PostgresDBRepository) Query(ctx context.Context, query string) (PostgresRows, error) {
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (r *PostgresDBRepository) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := r.conn.Exec(ctx, query, args)

	return err
}

func RegisterPostgresDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, dbhost, dbport, dbname))
	if err != nil {
		return nil, err
	}

	return conn, err
}
