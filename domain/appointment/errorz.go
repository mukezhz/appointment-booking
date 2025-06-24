package appointment

import (
	"github.com/mukezhz/appointment-booking/pkg/errorz"
)

var (
	// ErrBookingNotFound is returned when a booking is not found
	ErrBookingNotFound = errorz.NewNotFoundError("booking not found")
	// ErrInvalidBookingStatus is returned when an invalid booking status is provided
	ErrInvalidBookingStatus = errorz.NewBadRequestError("invalid booking status")
	// ErrInvalidTimeRange is returned when an invalid time range is provided
	ErrInvalidTimeRange = errorz.NewBadRequestError("invalid time range")
	// ErrInvalidBookingTime is returned when a booking time is invalid
	ErrInvalidBookingTime = errorz.NewBadRequestError("invalid booking time")
	// ErrSlotUnavailable is returned when a booking slot is not available
	ErrSlotUnavailable = errorz.NewConflictError("slot is not available")
	// ErrBookingExists is returned when a booking already exists
	ErrBookingExists = errorz.NewConflictError("booking already exists")

	ErrAvailabilityExists   = errorz.NewBadRequestError("Availability already exists for this time slot")
	ErrAvailabilityNotFound = errorz.NewNotFoundError("Availability not found")
)

var ErrAppointmentMap = map[error]bool{
	ErrBookingNotFound: true,
	ErrBookingExists:   true,
}
