package models

import (
	"time"

	"gorm.io/gorm"
)

type Availability struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"not null"`
	Weekday   string    `json:"weekday" gorm:"not null"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
}

// IsOverlapping checks if this availability overlaps with another
func (a *Availability) IsOverlapping(other *Availability) bool {
	if a.Weekday != other.Weekday {
		return false
	}
	return !a.StartTime.After(other.EndTime) && !other.StartTime.After(a.EndTime)
}

// Validate checks if the availability settings are valid
func (a *Availability) Validate() error {
	if a.StartTime.After(a.EndTime) {
		return ErrInvalidTimeRange
	}
	return nil
}
