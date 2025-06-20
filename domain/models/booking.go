package models

import (
	"clean-architecture/pkg/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Booking represents an appointment booking by a guest
type Booking struct {
	gorm.Model
	UUID       types.BinaryUUID `json:"uuid" gorm:"index;notnull;unique"`
	UserID     uint             `json:"user_id" gorm:"index;notnull"`
	User       User             `json:"-" gorm:"foreignKey:UserID"`
	GuestName  string           `json:"guest_name" gorm:"size:255;notnull"`
	GuestEmail string           `json:"guest_email" gorm:"size:255;notnull"`
	Date       time.Time        `json:"date" gorm:"notnull"`
	StartTime  string           `json:"start_time" gorm:"size:5;notnull"` // Format: HH:MM
	EndTime    string           `json:"end_time" gorm:"size:5;notnull"`   // Format: HH:MM
	Status     string           `json:"status" gorm:"size:20;default:'confirmed'"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.UUID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		b.UUID = types.BinaryUUID(id)
		return err
	}
	return nil
}

func (*Booking) TableName() string {
	return "bookings"
}
