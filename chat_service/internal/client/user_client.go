package client

import (
	"context"

	pb "chat_service/pb"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(
	client pb.UserServiceClient,
) *UserClient {

	return &UserClient{
		client: client,
	}
}
func (u *UserClient) UserExists(
	ctx context.Context,
	userID string,
) (bool, error) {

	res, err :=
		u.client.UserExists(
			ctx,
			&pb.UserExistsRequest{
				UserId: userID,
			},
		)

	if err != nil {
		return false, err
	}

	return res.Exists, nil
}
func (u *UserClient) GetUser(
	ctx context.Context,
	userID string,
) (*pb.UserResponse, error) {

	return u.client.GetUser(
		ctx,
		&pb.GetUserRequest{
			UserId: userID,
		},
	)
}
func (u *UserClient) SearchUsers(
	ctx context.Context,
	query string,
) (
	*pb.SearchUsersResponse,
	error,
) {

	return u.client.SearchUsers(
		ctx,
		&pb.SearchUsersRequest{
			Query: query,
		},
	)
}
