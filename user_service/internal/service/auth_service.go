package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"chat_app/internal/dto"
	"chat_app/internal/model"
	"chat_app/internal/repository"
	jwtutil "chat_app/pkg/jwt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) error {

	if req.Username == "" {
		return errors.New("username is required")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	existingUser, err :=
		s.userRepo.GetByEmail(req.Email)

	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err :=
		bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)

	if err != nil {
		return err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	return s.userRepo.Create(user)

}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest,
) (*dto.LoginResponse, error) {
	if req.Email == "" {
		return nil,
			errors.New("email is required")
	}

	if req.Password == "" {
		return nil,
			errors.New("password is required")
	}

	user, err :=
		s.userRepo.GetByEmail(req.Email)

	if err != nil {
		return nil,
			errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return nil,
			errors.New("invalid credentials")
	}

	accessToken, err :=
		jwtutil.GenerateAccessToken(
			user.ID.String(),
			user.Email,
			user.Role,
		)

	if err != nil {
		return nil, err
	}

	refreshToken, err :=
		jwtutil.GenerateRefreshToken(
			user.ID.String(),
		)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (string, error) {

	if refreshToken == "" {
		return "",
			errors.New("refresh token required")
	}

	userID, err :=
		jwtutil.ValidateRefreshToken(
			refreshToken,
		)

	if err != nil {
		return "",
			errors.New("invalid refresh token")
	}

	user, err :=
		s.userRepo.GetByID(userID)

	if err != nil {
		return "",
			errors.New("user not found")
	}

	newAccessToken, err :=
		jwtutil.GenerateAccessToken(
			user.ID.String(),
			user.Email,
			user.Role,
		)

	if err != nil {
		return "", err
	}

	return newAccessToken, nil

}

func (s *AuthService) Logout(
	ctx context.Context,
	userID string,
) error {

	// later:
	// delete refresh token from DB
	// blacklist token
	// revoke session

	return nil

}
