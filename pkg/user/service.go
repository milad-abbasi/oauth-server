package user

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	logger   *zap.Logger
	userRepo Repository
}

func NewService(logger *zap.Logger, userRepo Repository) *Service {
	return &Service{
		logger:   logger.Named("UserService"),
		userRepo: userRepo,
	}
}

func (s *Service) Register(ctx context.Context, user *TinyUser) (*User, error) {
	newUser := User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	// newUser.HashPassword()

	createdUser, err := s.userRepo.Create(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
