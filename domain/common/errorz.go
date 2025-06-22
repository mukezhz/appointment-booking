package common

import "clean-architecture/pkg/errorz"

var (
	ErrInvalidUserID = errorz.NewBadRequestError("Invalid user ID provided. Please check the user ID and try again.")
)
