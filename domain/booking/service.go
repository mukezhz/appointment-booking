package booking

import (
	"clean-architecture/domain/availability"
	"clean-architecture/domain/models"
	"clean-architecture/pkg/errorz"
	"context"
	"fmt"
	"time"
)

type Service struct {
	repository       *Repository
	availabilityRepo *availability.Repository
}

func NewService(repository *Repository, availabilityRepo *availability.Repository) *Service {
	return &Service{
		repository:       repository,
		availabilityRepo: availabilityRepo,
	}
}

func (s *Service) CreateBooking(ctx context.Context, userID uint, req *CreateBookingRequest) (*BookingResponse, error) {
	// Parse date and time
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errorz.NewBadRequestError("Invalid date format. Use YYYY-MM-DD")
	}

	// Validate time format
	_, err = time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, errorz.NewBadRequestError("Invalid start time format. Use HH:MM")
	}
	_, err = time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, errorz.NewBadRequestError("Invalid end time format. Use HH:MM")
	}

	// Check if the slot is available
	weekday := date.Weekday().String()
	availabilities, err := s.availabilityRepo.GetByWeekday(ctx, userID, weekday)
	if err != nil {
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}

	// Check if the requested time falls within any availability slot
	var isAvailable bool
	for _, slot := range availabilities {
		if req.StartTime >= slot.StartTime && req.EndTime <= slot.EndTime {
			isAvailable = true
			break
		}
	}

	if !isAvailable {
		return nil, errorz.NewBadRequestError("Selected time slot is not available")
	}

	// Check for overlapping bookings
	exists, err := s.repository.ExistsOverlappingBooking(ctx, userID, date, req.StartTime, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check overlapping bookings: %w", err)
	}

	if exists {
		return nil, errorz.NewBadRequestError("Selected time slot is already booked")
	}

	// Create booking
	booking := &models.Booking{
		UserID:     userID,
		GuestName:  req.GuestName,
		GuestEmail: req.GuestEmail,
		Date:       date,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Status:     "confirmed",
	}

	if err := s.repository.Create(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return ToBookingResponse(booking), nil
}

func (s *Service) GetUserBookings(ctx context.Context, userID uint) ([]*BookingResponse, error) {
	bookings, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}

	return ToBookingResponseList(bookings), nil
}
