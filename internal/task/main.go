package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/drifterz13/go-services/proto/task"
	"google.golang.org/grpc"

	"github.com/jackc/pgx/v4"
)

const (
	port = ":50051"
)

func main() {
	setupDB()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &taskService{})
	log.Printf("server listening at: %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)

	<-sigint
	s.GracefulStop()
	log.Panicln("shutdown gracefully.")
}

var (
	dbport   = 5432
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

func setupDB() {
	databaseURL := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", user, password, dbport, dbname)
	log.Printf("database url: %v\n", databaseURL)

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	log.Println("connected to postgresql database.")
}
