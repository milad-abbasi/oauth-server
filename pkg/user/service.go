package user

import (
	"context"

	"github.com/milad-abbasi/oauth-server/pkg/hash"
	"go.uber.org/zap"
)

type Service struct {
	logger   *zap.Logger
	UserRepo Repository
}

func NewService(logger *zap.Logger, userRepo Repository) *Service {
	return &Service{
		logger:   logger.Named("UserService"),
		UserRepo: userRepo,
	}
}

func (s *Service) NewUser(ctx context.Context, user *TinyUser) (*User, error) {
	hashedPassword, err := hash.GenerateArgon2Hash(user.Password)
	if err != nil {
		return nil, err
	}

	newUser := User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	createdUser, err := s.UserRepo.Create(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
