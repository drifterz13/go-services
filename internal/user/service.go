package main

import (
	"context"

	pb "github.com/drifterz13/go-services/internal/common/genproto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userService struct {
	pb.UnimplementedUserServiceServer
	repo UserRepository
}

func NewUserService(repo UserRepository) pb.UserServiceServer {
	return &userService{repo: repo}
}

func (s *userService) FindUsers(ctx context.Context, in *emptypb.Empty) (*pb.FindUsersResponse, error) {
	users, err := s.repo.Find(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, user.ToProto())
	}

	return &pb.FindUsersResponse{Users: pbUsers}, nil
}

func (s *userService) FindUser(ctx context.Context, in *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	return nil, nil
}

func (s *userService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*emptypb.Empty, error) {
	err := s.repo.Create(context.Background(), in.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
