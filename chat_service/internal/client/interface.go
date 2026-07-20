package client

import (
	"context"

	pb "chat_service/pb"
)

type UserService interface {
	UserExists(
		ctx context.Context,
		userID string,
	) (bool, error)

	GetUser(
		ctx context.Context,
		userID string,
	) (*pb.UserResponse, error)

	SearchUsers(
		ctx context.Context,
		query string,
	) (
		*pb.SearchUsersResponse,
		error,
	)
}
