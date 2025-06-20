package models

import (
	"clean-architecture/pkg/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Availability represents a weekly recurring availability slot
type Availability struct {
	gorm.Model
	UUID      types.BinaryUUID `json:"uuid" gorm:"index;notnull;unique"`
	UserID    uint             `json:"user_id" gorm:"index;notnull"`
	User      User             `json:"-" gorm:"foreignKey:UserID"`
	Weekday   string           `json:"weekday" gorm:"size:10;notnull"`   // Monday, Tuesday, etc.
	StartTime string           `json:"start_time" gorm:"size:5;notnull"` // Format: HH:MM
	EndTime   string           `json:"end_time" gorm:"size:5;notnull"`   // Format: HH:MM
}

func (a *Availability) BeforeCreate(tx *gorm.DB) error {
	if a.UUID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		a.UUID = types.BinaryUUID(id)
		return err
	}
	return nil
}

func (*Availability) TableName() string {
	return "availabilities"
}
