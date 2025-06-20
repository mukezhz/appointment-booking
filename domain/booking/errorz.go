package booking

import "clean-architecture/pkg/errorz"

var (
	ErrBookingNotFound      = errorz.NewNotFoundError("booking not found")
	ErrBookingAlreadyExists = errorz.NewConflictError("booking already exists")
	ErrInvalidBookingInput  = errorz.NewBadRequestError("invalid booking input")
)

var (
	ErrBookingMap = map[error]bool{
		ErrBookingNotFound:      true,
		ErrBookingAlreadyExists: true,
		ErrInvalidBookingInput:  true,
	}
)
