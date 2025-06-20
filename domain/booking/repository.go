package booking

import (
	"clean-architecture/domain/common"
	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"context"
	"time"
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

func (r *Repository) Create(ctx context.Context, booking *models.Booking) error {
	if err := r.db.WithContext(ctx).Create(booking).Error; err != nil {
		r.logger.Error("failed to create booking: ", err)
		return common.HandleDBError(err, ErrBookingMap)
	}
	return nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID uint) ([]*models.Booking, error) {
	var bookings []*models.Booking
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&bookings).Error
	if err != nil {
		r.logger.Error("failed to get user bookings: ", err)
	}
	return bookings, common.HandleDBError(err, ErrBookingMap)
}

func (r *Repository) GetByDateRange(ctx context.Context, userID uint, start, end time.Time) ([]*models.Booking, error) {
	var bookings []*models.Booking
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, start, end).
		Find(&bookings).Error
	if err != nil {
		r.logger.Error("failed to get bookings by date range: ", err)
		return nil, common.HandleDBError(err, ErrBookingMap)
	}
	return bookings, nil
}

func (r *Repository) ExistsOverlappingBooking(ctx context.Context, userID uint, date time.Time, startTime, endTime string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("user_id = ? AND date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
			userID, date, endTime, startTime, endTime, startTime).
		Count(&count).Error
	if err != nil {
		r.logger.Error("failed to check overlapping bookings: ", err)
		return false, common.HandleDBError(err, ErrBookingMap)
	}
	r.logger.Debug("Overlapping bookings count: ", count)
	return count > 0, nil
}
