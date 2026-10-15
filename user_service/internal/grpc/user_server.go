package grpc

import (
	"context"
	"errors"
	"log"
	"strconv"

	pb "chat_app/pb"

	"chat_app/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer

	userRepo *repository.UserRepository
}

func NewUserServer(
	userRepo *repository.UserRepository,
) *UserServer {

	return &UserServer{
		userRepo: userRepo,
	}
}
func (s *UserServer) UserExists(
	ctx context.Context,
	req *pb.UserExistsRequest,
) (
	*pb.UserExistsResponse,
	error,
) {

	_, err := s.userRepo.GetByID(req.UserId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &pb.UserExistsResponse{Exists: false}, nil
	}

	if err != nil {
		log.Println("UserExists lookup failed:", err)
		return nil, status.Error(codes.Internal, "user lookup failed")
	}

	return &pb.UserExistsResponse{Exists: true}, nil
}

func (s *UserServer) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (
	*pb.UserResponse,
	error,
) {

	user, err :=
		s.userRepo.GetByID(

			req.UserId,
		)

	if err != nil {

		return nil, err
	}

	if user == nil {

		return &pb.UserResponse{}, nil
	}

	return &pb.UserResponse{
		Id:       strconv.FormatUint(uint64(user.ID), 10),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
func (s *UserServer) SearchUsers(
	ctx context.Context,
	req *pb.SearchUsersRequest,
) (
	*pb.SearchUsersResponse,
	error,
) {

	users, err := s.userRepo.SearchUsers(

		req.Query,
	)

	if err != nil {
		return nil, err
	}

	response := &pb.SearchUsersResponse{}

	for _, user := range users {

		response.Users = append(
			response.Users,
			&pb.User{
				Id:       strconv.FormatUint(uint64(user.ID), 10),
				Username: user.Username,
				Email:    user.Email,
			},
		)
	}

	return response, nil
}
