package availability

import (
	"clean-architecture/domain/common"
	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"context"
)

type Repository struct {
	db     infrastructure.Database
	logger framework.Logger
}

func NewRepository(
	db infrastructure.Database,
	logger framework.Logger,
) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) Create(ctx context.Context, availability *models.Availability) error {
	if err := r.db.WithContext(ctx).Create(availability).Error; err != nil {
		r.logger.Error("failed to create availability: ", err)
		return common.HandleDBError(err, ErrAvailabilityMap)
	}
	return nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID uint) ([]*models.Availability, error) {
	var availabilities []*models.Availability
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&availabilities).Error
	if err != nil {
		r.logger.Error("failed to get user availability: ", err)
	}
	return availabilities, common.HandleDBError(err, ErrAvailabilityMap)
}

func (r *Repository) GetByWeekday(ctx context.Context, userID uint, weekday string) ([]*models.Availability, error) {
	var availabilities []*models.Availability
	err := r.db.WithContext(ctx).Where("user_id = ? AND weekday = ?", userID, weekday).Find(&availabilities).Error
	if err != nil {
		r.logger.Error("failed to get availability by weekday: ", err)
	}
	return availabilities, common.HandleDBError(err, ErrAvailabilityMap)
}
