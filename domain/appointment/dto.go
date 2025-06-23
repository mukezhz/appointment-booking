package appointment

import (
	"time"

	"github.com/mukezhz/appointment-booking/pkg/responses"
)

// CreateAvailabilityRequest DTO for creating availability
type CreateAvailabilityRequest struct {
	Weekday   string `json:"weekday" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

// AvailabilityResponse DTO for availability response
type AvailabilityResponse struct {
	ID        uint      `json:"id"`
	Weekday   string    `json:"weekday"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	// BookingStatusPending represents a pending booking
	BookingStatusPending BookingStatus = "pending"
	// BookingStatusConfirmed represents a confirmed booking
	BookingStatusConfirmed BookingStatus = "confirmed"
	// BookingStatusCancelled represents a cancelled booking
	BookingStatusCancelled BookingStatus = "cancelled"
	// BookingStatusCompleted represents a completed booking
	BookingStatusCompleted BookingStatus = "completed"
)

// CreateBookingRequest DTO for creating a booking
type CreateBookingRequest struct {
	GuestName  string `json:"guest_name" binding:"required"`
	GuestEmail string `json:"guest_email" binding:"required,email"`
	Date       string `json:"date" binding:"required"`
	Time       string `json:"time" binding:"required"`
}

// BookingResponse DTO for booking response
type BookingResponse struct {
	ID         uint          `json:"id"`
	GuestName  string        `json:"guest_name"`
	GuestEmail string        `json:"guest_email"`
	Date       time.Time     `json:"date"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Status     BookingStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// UpdateBookingStatusRequest DTO for updating booking status
type UpdateBookingStatusRequest struct {
	Status BookingStatus `json:"status" binding:"required"`
}

// BookingListResponse DTO for paginated booking list
type BookingListResponse = responses.ListResponseType[BookingResponse]

// AvailabilityListResponse DTO for paginated availability list
type AvailabilityListResponse = responses.ListResponseType[AvailabilityResponse]
