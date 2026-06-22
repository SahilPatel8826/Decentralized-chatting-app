package service

import (
	"chat_app/internal/dto"
	"chat_app/internal/model"
	"chat_app/internal/repository"
	"context"
	"errors"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
) *UserService {

	return &UserService{
		userRepo: userRepo,
	}
}
func (s *UserService) GetByID(
	ctx context.Context,
	id string,
) (*model.User, error) {

	if id == "" {
		return nil,
			errors.New("user id required")
	}

	return s.userRepo.GetByID(id)
}
func (s *UserService) UpdateProfile(
	ctx context.Context,
	userID string,
	req dto.UpdateProfileRequest,
) error {

	user, err :=
		s.userRepo.GetByID(userID)

	if err != nil {
		return err
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	return s.userRepo.Update(user)
}
func (s *UserService) Delete(
	ctx context.Context,
	userID string,
) error {

	user, err :=
		s.userRepo.GetByID(userID)

	if err != nil {
		return err
	}

	return s.userRepo.Delete(
		user.ID.String(),
	)
}
func (s *UserService) Search(
	ctx context.Context,
	query string,
) ([]model.User, error) {

	if query == "" {
		return []model.User{}, nil
	}

	return s.userRepo.Search(query)
}
