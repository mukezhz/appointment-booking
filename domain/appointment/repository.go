package appointment

import (
	"github.com/mukezhz/appointment-booking/domain/common"
	"github.com/mukezhz/appointment-booking/domain/models"
	"github.com/mukezhz/appointment-booking/pkg/framework"
	"github.com/mukezhz/appointment-booking/pkg/infrastructure"
	"github.com/mukezhz/appointment-booking/pkg/types"
)

// Repository handles database operations for appointments
type Repository struct {
	db     infrastructure.Database
	logger framework.Logger
}

// NewRepository creates a new appointment repository
func NewRepository(
	db infrastructure.Database,
	logger framework.Logger,
) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CreateAvailability creates a new availability
func (r *Repository) CreateAvailability(availability *models.Availability) error {
	err := r.db.Create(availability).Error
	if err != nil {
		r.logger.Error("[Repository...CreateAvailability] Error creating availability:", err)
		return common.HandleDBError(err, ErrAppointmentMap)
	}
	return nil
}

// GetAvailabilityByUserID retrieves all availabilities for a user
func (r *Repository) GetAvailabilityByUserID(userID uint) ([]models.Availability, error) {
	var availabilities []models.Availability
	err := r.db.Where("user_id = ?", userID).Find(&availabilities).Error
	if err != nil {
		r.logger.Error("[Repository...GetAvailabilityByUserID] Error retrieving availabilities:", err)
		return nil, common.HandleDBError(err, ErrAppointmentMap)
	}
	return availabilities, nil
}

// DeleteAvailability deletes an availability by ID
func (r *Repository) DeleteAvailability(id uint) error {
	err := r.db.Delete(&models.Availability{}, id).Error
	if err != nil {
		r.logger.Error("[Repository...DeleteAvailability] Error deleting availability:", err)
		return common.HandleDBError(err, ErrAppointmentMap)
	}
	return nil
}

// CreateBooking creates a new booking
func (r *Repository) CreateBooking(booking *models.Booking) error {
	err := r.db.Create(booking).Error
	if err != nil {
		r.logger.Error("[Repository...CreateBooking] Error creating booking:", err)
		return common.HandleDBError(err, ErrAppointmentMap)
	}
	return nil
}

// GetBooking retrieves a booking by ID
func (r *Repository) GetBooking(id types.BinaryUUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Where("id = ?", id).First(&booking).Error
	if err != nil {
		r.logger.Error("[Repository...GetBooking] Error retrieving booking:", err)
		return nil, common.HandleDBError(err, ErrAppointmentMap)
	}
	return &booking, nil
}

// GetBookingByID retrieves a booking by ID
func (r *Repository) GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.First(&booking, id).Error
	if err != nil {
		r.logger.Error("[Repository...GetBookingByID] Error retrieving booking:", err)
		return nil, common.HandleDBError(err, ErrAppointmentMap)
	}
	return &booking, nil
}

// GetBookings retrieves all bookings for a user with pagination
func (r *Repository) GetBookings(userID uint, offset, limit int) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error
	if err != nil {
		r.logger.Error("[Repository...GetBookings] Error retrieving bookings:", err)
		return nil, common.HandleDBError(err, ErrAppointmentMap)
	}
	return bookings, nil
}

// GetBookingsByUserID retrieves all bookings for a user
func (r *Repository) GetBookingsByUserID(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("user_id = ?", userID).Find(&bookings).Error
	if err != nil {
		r.logger.Error("[Repository...GetBookingsByUserID] Error retrieving bookings:", err)
		return nil, common.HandleDBError(err, ErrAppointmentMap)
	}
	return bookings, nil
}

// GetTotalBookings returns the total number of bookings for a user
func (r *Repository) GetTotalBookings(userID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).
		Where("user_id = ?", userID).
		Count(&count).
		Error
	if err != nil {
		r.logger.Error("[Repository...GetTotalBookings] Error counting bookings:", err)
		return 0, common.HandleDBError(err, ErrAppointmentMap)
	}
	return int(count), nil
}

// UpdateBooking updates a booking
func (r *Repository) UpdateBooking(booking *models.Booking) error {
	err := r.db.Save(booking).Error
	if err != nil {
		r.logger.Error("[Repository...UpdateBooking] Error updating booking:", err)
		return common.HandleDBError(err, ErrAppointmentMap)
	}
	return nil
}

// UpdateBookingStatus updates the status of a booking
func (r *Repository) UpdateBookingStatus(id uint, status models.BookingStatus) error {
	err := r.db.Model(&models.Booking{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		r.logger.Error("[Repository...UpdateBookingStatus] Error updating booking status:", err)
		return common.HandleDBError(err, ErrAppointmentMap)
	}
	return nil
}

// IsValidBookingStatus checks if a booking status is valid
func IsValidBookingStatus(status models.BookingStatus) bool {
	switch status {
	case models.BookingStatusPending,
		models.BookingStatusConfirmed,
		models.BookingStatusCanceled:
		return true
	default:
		return false
	}
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		r.logger.Error("[Repository...GetUserByEmail] Error retrieving user by email:", err)
		return nil, common.HandleDBError(err, ErrAppointmentMap)
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *models.User) (*models.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		r.logger.Error("[Repository...CreateUser] Error creating user:", err)
		return nil, common.HandleDBError(err, ErrAppointmentMap)
	}
	return user, nil
}
