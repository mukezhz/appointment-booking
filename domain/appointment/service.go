package appointment

import (
	"github.com/mukezhz/appointment-booking/domain/constants"
	"github.com/mukezhz/appointment-booking/domain/models"
	"github.com/mukezhz/appointment-booking/pkg/framework"
)

// Service handles business logic for appointments
type Service struct {
	logger framework.Logger
	repo   *Repository
}

// NewService creates a new appointment service
func NewService(
	logger framework.Logger,
	repo *Repository,
) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

// CreateAvailability creates a new availability slot
func (s *Service) CreateAvailability(availability *models.Availability) error {
	s.logger.Info("[AppointmentService...CreateAvailability]")
	return s.repo.CreateAvailability(availability)
}

// GetAvailabilityByUserID gets all availability slots for a user
func (s *Service) GetAvailabilityByUserID(userID uint) ([]models.Availability, error) {
	s.logger.Info("[AppointmentService...GetAvailabilityByUserID]")
	return s.repo.GetAvailabilityByUserID(userID)
}

// DeleteAvailability deletes an availability slot
func (s *Service) DeleteAvailability(id uint) error {
	s.logger.Info("[AppointmentService...DeleteAvailability]")
	return s.repo.DeleteAvailability(id)
}

// CreateBooking creates a new booking
func (s *Service) CreateBooking(booking *models.Booking) error {
	s.logger.Info("[AppointmentService...CreateBooking]")

	// Try to find existing user by email
	existingUser, err := s.repo.GetUserByEmail(booking.GuestEmail)
	if err != nil {
		// User doesn't exist, create a new one
		s.logger.Info("[AppointmentService...CreateBooking] Creating new user for:", booking.GuestEmail)

		existingUser, err = s.repo.CreateUser(
			&models.User{
				Email:    booking.GuestEmail,
				FullName: booking.GuestName,
				Role:     constants.UserRoleGuest,
			},
		)
		if err != nil {
			s.logger.Error("[AppointmentService...CreateBooking] Failed to create user:", err)
			return err
		}
	}

	// Set the guest user ID from the found or newly created user
	booking.UserID = existingUser.ID

	// Validate the booking
	if err := booking.Validate(); err != nil {
		return err
	}

	// Create the booking
	return s.repo.CreateBooking(booking)
}

// GetBookingByID gets a booking by ID
func (s *Service) GetBookingByID(id uint) (*models.Booking, error) {
	s.logger.Info("[AppointmentService...GetBookingByID]")
	return s.repo.GetBookingByID(id)
}

// GetBookingsByUserID gets all bookings for a user
func (s *Service) GetBookingsByUserID(userID uint) ([]models.Booking, error) {
	s.logger.Info("[AppointmentService...GetBookingsByUserID]")
	return s.repo.GetBookingsByUserID(userID)
}

// GetBookings retrieves all bookings for a user with pagination
func (s *Service) GetBookings(userID uint, page, limit int) ([]models.Booking, int, error) {
	s.logger.Info("[AppointmentService...GetBookings]")
	offset := (page - 1) * limit
	bookings, err := s.repo.GetBookings(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.GetTotalBookings(userID)
	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

// UpdateBookingStatus updates the status of a booking
func (s *Service) UpdateBookingStatus(id uint, status models.BookingStatus) error {
	s.logger.Info("[AppointmentService...UpdateBookingStatus]")

	if !IsValidBookingStatus(status) {
		return ErrInvalidBookingStatus
	}

	return s.repo.UpdateBookingStatus(id, status)
}
