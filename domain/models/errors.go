package models

import "errors"

var (
	ErrInvalidTimeRange   = errors.New("start time must be before end time")
	ErrSlotUnavailable    = errors.New("slot is not available")
	ErrInvalidBookingTime = errors.New("booking time is invalid")
	ErrBookingExists      = errors.New("booking already exists for this slot")
	ErrUnauthorized       = errors.New("unauthorized access")
)
