package appointment

import (
	"errors"

	"github.com/mukezhz/appointment-booking/pkg/errorz"
)

var (
	// ErrBookingNotFound is returned when a booking is not found
	ErrBookingNotFound = errors.New("booking not found")
	// ErrInvalidBookingStatus is returned when an invalid booking status is provided
	ErrInvalidBookingStatus = errors.New("invalid booking status")
	// ErrInvalidTimeRange is returned when an invalid time range is provided
	ErrInvalidTimeRange = errors.New("invalid time range")
	// ErrInvalidBookingTime is returned when a booking time is invalid
	ErrInvalidBookingTime = errors.New("invalid booking time")
	// ErrSlotUnavailable is returned when a booking slot is not available
	ErrSlotUnavailable = errors.New("slot is not available")
	// ErrBookingExists is returned when a booking already exists
	ErrBookingExists = errors.New("booking already exists")

	ErrAvailabilityExists   = errorz.ErrBadRequest.JoinError("Availability already exists for this time slot")
	ErrAvailabilityNotFound = errorz.ErrNotFound.JoinError("Availability not found")
)
