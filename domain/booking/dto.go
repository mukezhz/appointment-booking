package booking

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/types"
	"time"
)

type CreateBookingRequest struct {
	GuestName  string `json:"guest_name" binding:"required"`
	GuestEmail string `json:"guest_email" binding:"required,email"`
	Date       string `json:"date" binding:"required"` // Format: YYYY-MM-DD
	StartTime  string `json:"start_time" binding:"required,len=5"`
	EndTime    string `json:"end_time" binding:"required,len=5"`
}

type BookingResponse struct {
	UUID       types.BinaryUUID `json:"uuid"`
	GuestName  string           `json:"guest_name"`
	GuestEmail string           `json:"guest_email"`
	Date       time.Time        `json:"date"`
	StartTime  string           `json:"start_time"`
	EndTime    string           `json:"end_time"`
	Status     string           `json:"status"`
}

func ToBookingResponse(booking *models.Booking) *BookingResponse {
	return &BookingResponse{
		UUID:       booking.UUID,
		GuestName:  booking.GuestName,
		GuestEmail: booking.GuestEmail,
		Date:       booking.Date,
		StartTime:  booking.StartTime,
		EndTime:    booking.EndTime,
		Status:     booking.Status,
	}
}

func ToBookingResponseList(bookings []*models.Booking) []*BookingResponse {
	var response []*BookingResponse
	for _, booking := range bookings {
		response = append(response, ToBookingResponse(booking))
	}
	return response
}
