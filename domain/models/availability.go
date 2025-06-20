package models

import (
	"time"

	"gorm.io/gorm"

	"clean-architecture/pkg/types"
)

type Availability struct {
	gorm.Model
	UUID         types.BinaryUUID `json:"uuid" gorm:"type:binary(16);uniqueIndex;not null"`
	DoctorID     types.BinaryUUID `json:"doctor_id" gorm:"type:binary(16);index;not null"`
	DayOfWeek    *int8            `json:"day_of_week,omitempty" gorm:"type:tinyint;index"` // 1 = Monday, 7 = Sunday
	StartDate    *time.Time       `json:"start_date,omitempty" gorm:"type:date;index"`
	EndDate      *time.Time       `json:"end_date,omitempty" gorm:"type:date"`
	StartTime    time.Time        `json:"start_time" gorm:"type:time;not null"`
	EndTime      time.Time        `json:"end_time" gorm:"type:time;not null"`
	SlotDuration int              `json:"slot_duration" gorm:"not null"` // Duration in minutes
	IsRecurring  bool             `json:"is_recurring" gorm:"default:false"`
}

func (a *Availability) TableName() string {
	return "availabilities"
}

// OverlapsWith checks if this availability overlaps with another
func (a *Availability) OverlapsWith(other *Availability) bool {
	// If they're for different days of the week, no overlap
	if a.DayOfWeek != nil && other.DayOfWeek != nil {
		return *a.DayOfWeek == *other.DayOfWeek && a.timeOverlaps(other)
	}

	// If they're for specific dates
	if a.StartDate != nil && other.StartDate != nil {
		// If dates don't overlap, no overlap
		if a.StartDate.After(*other.EndDate) || other.StartDate.After(*a.EndDate) {
			return false
		}
		return a.timeOverlaps(other)
	}

	return false
}

// timeOverlaps checks if time ranges overlap
func (a *Availability) timeOverlaps(other *Availability) bool {
	// Convert to comparable format (minutes since midnight)
	aStart := a.StartTime.Hour()*60 + a.StartTime.Minute()
	aEnd := a.EndTime.Hour()*60 + a.EndTime.Minute()
	bStart := other.StartTime.Hour()*60 + other.StartTime.Minute()
	bEnd := other.EndTime.Hour()*60 + other.EndTime.Minute()

	return aStart < bEnd && bStart < aEnd
}

// Slot represents an available time slot for appointments
type Slot struct {
	StartTime time.Time
	EndTime   time.Time
	Available bool // Whether the slot is available for booking
	Cooldown  bool // Whether the slot is in cooldown period
}

// AvailableSlots represents available slots for a specific date
type AvailableSlots struct {
	Date  time.Time
	Slots []Slot
}
