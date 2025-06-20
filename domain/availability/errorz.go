package availability

import (
	"clean-architecture/pkg/errorz"
)

var (
	// Model validation errors
	ErrPastAppointmentDate       = errorz.NewBadRequestError("appointment date cannot be in the past")
	ErrInvalidTimeRange          = errorz.NewBadRequestError("end time must be after start time")
	ErrInvalidAppointmentStatus  = errorz.NewBadRequestError("invalid appointment status")
	ErrInvalidSlotDuration       = errorz.NewBadRequestError("slot duration must be between 15 and 120 minutes")
	ErrInvalidDayOfWeek          = errorz.NewBadRequestError("day of week must be between 1 and 7")
	ErrInvalidAvailabilityConfig = errorz.NewBadRequestError("either day_of_week or start_date must be set, but not both")

	// Not found errors
	ErrDoctorNotFound       = errorz.NewNotFoundError("doctor not found")
	ErrAppointmentNotFound  = errorz.NewNotFoundError("appointment not found")
	ErrAvailabilityNotFound = errorz.NewNotFoundError("availability slot not found")

	// Business logic errors
	ErrSlotNotAvailable    = errorz.NewBadRequestError("slot is not available")
	ErrSlotAlreadyBooked   = errorz.NewBadRequestError("slot is already booked")
	ErrSlotInCooldown      = errorz.NewBadRequestError("slot is in cooldown period")
	ErrSlotLocked          = errorz.NewBadRequestError("slot is currently locked by another user")
	ErrCancellationTooLate = errorz.NewBadRequestError("appointments can only be cancelled at least 24 hours before the scheduled time")
	ErrCancellationReason  = errorz.NewBadRequestError("cancellation reason is required when cancelling an appointment")
	ErrInvalidSlotOverlap  = errorz.NewBadRequestError("slot overlaps with existing availability")
	ErrDoctorUnavailable   = errorz.NewBadRequestError("doctor is not available at the requested time")

	ErrMissingCancellationReason = errorz.NewBadRequestError("cancellation reason is required when cancelling an appointment")
)

// Error maps for database error handling
var (
	AvailabilityErrMap = map[error]bool{
		ErrAvailabilityNotFound: true,
		ErrDoctorNotFound:       true,
	}

	AppointmentErrMap = map[error]bool{
		ErrAppointmentNotFound: true,
		ErrDoctorNotFound:      true,
	}
)
