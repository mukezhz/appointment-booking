package models

import (
	"time"

	"gorm.io/gorm"

	"clean-architecture/pkg/types"
)

type AppointmentStatus string

const (
	AppointmentScheduled AppointmentStatus = "scheduled"
	AppointmentCancelled AppointmentStatus = "cancelled"
	AppointmentCompleted AppointmentStatus = "completed"
	AppointmentNoShow    AppointmentStatus = "no_show"
)

type Appointment struct {
	gorm.Model
	UUID               types.BinaryUUID  `json:"uuid" gorm:"type:binary(16);uniqueIndex;not null"`
	DoctorID           types.BinaryUUID  `json:"doctor_id" gorm:"type:binary(16);index;not null"`
	PatientID          types.BinaryUUID  `json:"patient_id" gorm:"type:binary(16);index;not null"`
	AvailabilityID     uint              `json:"availability_id" gorm:"not null"`
	AppointmentDate    time.Time         `json:"appointment_date" gorm:"type:date;not null;index"`
	StartTime          time.Time         `json:"start_time" gorm:"type:time;not null"`
	EndTime            time.Time         `json:"end_time" gorm:"type:time;not null"`
	Status             AppointmentStatus `json:"status" gorm:"type:varchar(20);default:scheduled;index"`
	CancellationReason *string           `json:"cancellation_reason,omitempty" gorm:"type:text"`
	LockKey            *string           `json:"lock_key,omitempty" gorm:"type:varchar(255)"` // For concurrency control
	LockExpiry         *time.Time        `json:"lock_expiry,omitempty"`                       // For concurrency control
	CooldownUntil      *time.Time        `json:"cooldown_until,omitempty"`                    // For cancelled slot cooldown
}

func (a *Appointment) TableName() string {
	return "appointments"
}

// IsLocked checks if the appointment slot is currently locked
func (a *Appointment) IsLocked() bool {
	return a.LockKey != nil && a.LockExpiry != nil && time.Now().Before(*a.LockExpiry)
}

// IsInCooldown checks if the appointment slot is in cooldown period after cancellation
func (a *Appointment) IsInCooldown() bool {
	return a.CooldownUntil != nil && time.Now().Before(*a.CooldownUntil)
}

// CanBeCancelled checks if the appointment can be cancelled
func (a *Appointment) CanBeCancelled() bool {
	return a.Status == AppointmentScheduled &&
		time.Now().Add(24*time.Hour).Before(a.AppointmentDate) // 24-hour cancellation policy
}
