package availability

import "clean-architecture/pkg/errorz"

var (
	ErrAvailabilityNotFound      = errorz.NewNotFoundError("availability not found")
	ErrAvailabilityAlreadyExists = errorz.NewConflictError("availability already exists")
	ErrInvalidAvailabilityInput  = errorz.NewBadRequestError("invalid availability input")
)

var (
	ErrAvailabilityMap = map[error]bool{
		ErrAvailabilityNotFound:      true,
		ErrAvailabilityAlreadyExists: true,
	}
)
