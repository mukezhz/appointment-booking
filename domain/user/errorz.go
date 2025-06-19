package user

import (
	"clean-architecture/pkg/errorz"
)

var (
	ErrUserNotFound       = errorz.NewNotFoundError("user not found")
	ErrUserUnauthorized   = errorz.NewUnauthorizedError("user unauthorized")
	ErrUserInvalidUserID  = errorz.NewBadRequestError("invalid user id format")
	ErrInvalidCredentials = errorz.NewUnauthorizedError("invalid credentials")
	ErrEmailTaken         = errorz.NewBadRequestError("email already taken")
	ErrInvalidRole        = errorz.NewBadRequestError("invalid role")
)

var ErrUserMap = map[error]bool{
	ErrUserNotFound: true,
}
