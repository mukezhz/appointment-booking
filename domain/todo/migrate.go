package todo

import (
	"github.com/mukezhz/appointment-booking/domain/models"
	"github.com/mukezhz/appointment-booking/pkg/infrastructure"
)

// Migrate automigrates the todo model
func Migrate(db infrastructure.Database) {
	db.AutoMigrate(&models.Todo{})
}
