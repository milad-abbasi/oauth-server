package user

import "errors"

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
