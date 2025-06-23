package appointment

import (
	"gorm.io/gorm"

	"github.com/mukezhz/appointment-booking/domain/models"
	"github.com/mukezhz/appointment-booking/pkg/types"
)

// Repository handles database operations for appointments
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new appointment repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// CreateAvailability creates a new availability
func (r *Repository) CreateAvailability(availability *models.Availability) error {
	return r.db.Create(availability).Error
}

// GetAvailabilityByUserID retrieves all availabilities for a user
func (r *Repository) GetAvailabilityByUserID(userID uint) ([]models.Availability, error) {
	var availabilities []models.Availability
	err := r.db.Where("user_id = ?", userID).Find(&availabilities).Error
	return availabilities, err
}

// DeleteAvailability deletes an availability by ID
func (r *Repository) DeleteAvailability(id uint) error {
	return r.db.Delete(&models.Availability{}, id).Error
}

// CreateBooking creates a new booking
func (r *Repository) CreateBooking(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

// GetBooking retrieves a booking by ID
func (r *Repository) GetBooking(id types.BinaryUUID) (*models.Booking, error) {
	var booking models.Booking
	result := r.db.Where("id = ?", id).First(&booking)
	if result.Error != nil {
		return nil, ErrBookingNotFound
	}
	return &booking, nil
}

// GetBookingByID retrieves a booking by ID
func (r *Repository) GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	result := r.db.First(&booking, id)
	if result.Error != nil {
		return nil, ErrBookingNotFound
	}
	return &booking, nil
}

// GetBookings retrieves all bookings for a user with pagination
func (r *Repository) GetBookings(userID uint, offset, limit int) ([]models.Booking, error) {
	var bookings []models.Booking
	result := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

// GetBookingsByUserID retrieves all bookings for a user
func (r *Repository) GetBookingsByUserID(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	result := r.db.Where("user_id = ?", userID).Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

// GetTotalBookings returns the total number of bookings for a user
func (r *Repository) GetTotalBookings(userID uint) (int, error) {
	var count int64
	result := r.db.Model(&models.Booking{}).
		Where("user_id = ?", userID).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

// UpdateBooking updates a booking
func (r *Repository) UpdateBooking(booking *models.Booking) error {
	result := r.db.Save(booking)
	return result.Error
}

// UpdateBookingStatus updates the status of a booking
func (r *Repository) UpdateBookingStatus(id uint, status models.BookingStatus) error {
	result := r.db.Model(&models.Booking{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrBookingNotFound
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
