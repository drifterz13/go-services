package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/drifterz13/go-services/internal/common/db"
	pb "github.com/drifterz13/go-services/internal/common/genproto/user"
	"google.golang.org/grpc"
)

var (
	port = ":50053"
)

func main() {
	log.Println("initialized user service.")

	conn, err := db.RegisterPostgresDB()
	repo := NewUserRepository(db.NewPostgresDBRepository(conn))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, NewUserService(repo))

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
