package user

import (
	"context"
)

type Repository interface {
	Exists(ctx context.Context, email string) (bool, error)
	FindOne(ctx context.Context, user *User) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
}
