package user

import (
	"context"
)

type Repository interface {
	Exists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *User) (*User, error)
}
