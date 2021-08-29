package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/drifterz13/go-services/internal/common/genproto/task"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

var (
	port = ":50051"
)

func main() {
	conn, err := registerDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v\n", err)
	}

	repo := NewTaskRepository(conn)

	// Setup gRPC
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, NewTaskService(repo))
	log.Printf("server listening at: %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)

	<-sigint
	s.GracefulStop()
	log.Println("shutdown gracefully.")
}

func registerDB() (*pgx.Conn, error) {
	var (
		dbhost   = os.Getenv("POSTGRES_HOST")
		dbport   = os.Getenv("POSTGRES_PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, dbhost, dbport, dbname))
	if err != nil {
		return nil, err
	}

	return conn, err
}
