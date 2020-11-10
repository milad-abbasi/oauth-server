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

func (s *Service) NewUser(ctx context.Context, user *TinyUser) (*User, error) {
	hashedPassword, err := GenerateHash(user.Password)
	if err != nil {
		return nil, err
	}

	newUser := User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	createdUser, err := s.userRepo.Create(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
