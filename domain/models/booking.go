package models

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCanceled  BookingStatus = "canceled"
)

type Booking struct {
	gorm.Model
	UserID     uint          `json:"user_id" gorm:"not null"`
	GuestName  string        `json:"guest_name" gorm:"not null"`
	GuestEmail string        `json:"guest_email" gorm:"not null"`
	Date       time.Time     `json:"date" gorm:"not null"`
	StartTime  time.Time     `json:"start_time" gorm:"not null"`
	EndTime    time.Time     `json:"end_time" gorm:"not null"`
	Status     BookingStatus `json:"status" gorm:"not null;default:'pending'"`
}

// Validate checks if the booking is valid
func (b *Booking) Validate() error {
	if b.StartTime.After(b.EndTime) {
		return ErrInvalidTimeRange
	}
	if b.StartTime.Before(time.Now()) {
		return ErrInvalidBookingTime
	}
	return nil
}

// IsOverlapping checks if this booking overlaps with another
func (b *Booking) IsOverlapping(other *Booking) bool {
	return !b.StartTime.After(other.EndTime) && !other.StartTime.After(b.EndTime)
}
