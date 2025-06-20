package availability

import (
	"context"
	"time"

	"clean-architecture/domain/common"
	"clean-architecture/domain/models"
	"clean-architecture/pkg/errorz"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/types"
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

// GetDoctorByID retrieves a doctor by their UUID
func (r *Repository) GetDoctorByID(ctx context.Context, uuid types.BinaryUUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("uuid = ? AND role = 'doctor'", uuid).First(&user).Error; err != nil {
		r.logger.Error("failed to find doctor: ", err)
		return nil, common.HandleDBError(err, AvailabilityErrMap)
	}
	return &user, nil
}

// CreateAvailability creates a new availability slot
func (r *Repository) CreateAvailability(ctx context.Context, availability *models.Availability) error {
	if err := r.db.WithContext(ctx).Create(availability).Error; err != nil {
		r.logger.Error("failed to create availability: ", err)
		return common.HandleDBError(err, AvailabilityErrMap)
	}
	return nil
}

// UpdateAvailability updates an existing availability slot
func (r *Repository) UpdateAvailability(ctx context.Context, availability *models.Availability) error {
	if err := r.db.WithContext(ctx).Save(availability).Error; err != nil {
		r.logger.Error("failed to update availability: ", err)
		return common.HandleDBError(err, AvailabilityErrMap)
	}
	return nil
}

// DeleteAvailability soft-deletes an availability slot
func (r *Repository) DeleteAvailability(ctx context.Context, uuid types.BinaryUUID) error {
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Availability{}).Error; err != nil {
		r.logger.Error("failed to delete availability: ", err)
		return common.HandleDBError(err, AvailabilityErrMap)
	}
	return nil
}

// GetAvailabilityByUUID retrieves an availability slot by UUID
func (r *Repository) GetAvailabilityByUUID(ctx context.Context, uuid types.BinaryUUID) (*models.Availability, error) {
	var availability models.Availability
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&availability).Error; err != nil {
		r.logger.Error("failed to get availability: ", err)
		return nil, common.HandleDBError(err, AvailabilityErrMap)
	}
	return &availability, nil
}

// ListDoctorAvailability lists all availability slots for a doctor within a date range
func (r *Repository) ListDoctorAvailability(ctx context.Context, doctorID types.BinaryUUID, fromDate, toDate time.Time) ([]models.Availability, error) {
	var availabilities []models.Availability
	query := r.db.WithContext(ctx).Where("doctor_id = ?", doctorID)

	// Handle recurring and specific date availabilities
	query = query.Where(`
		(is_recurring = true AND (day_of_week IS NOT NULL OR (start_date <= ? AND (end_date IS NULL OR end_date >= ?))))
		OR
		(is_recurring = false AND start_date <= ? AND (end_date IS NULL OR end_date >= ?))
	`, toDate, fromDate, toDate, fromDate)

	if err := query.Find(&availabilities).Error; err != nil {
		r.logger.Error("failed to list availabilities: ", err)
		return nil, common.HandleDBError(err, AvailabilityErrMap)
	}
	return availabilities, nil
}

// CreateAppointment creates a new appointment
func (r *Repository) CreateAppointment(ctx context.Context, appointment *models.Appointment) error {
	if err := r.db.WithContext(ctx).Create(appointment).Error; err != nil {
		r.logger.Error("failed to create appointment: ", err)
		return common.HandleDBError(err, AppointmentErrMap)
	}
	return nil
}

// UpdateAppointment updates an existing appointment
func (r *Repository) UpdateAppointment(ctx context.Context, appointment *models.Appointment) error {
	if err := r.db.WithContext(ctx).Save(appointment).Error; err != nil {
		r.logger.Error("failed to update appointment: ", err)
		return common.HandleDBError(err, AppointmentErrMap)
	}
	return nil
}

// GetAppointmentByUUID retrieves an appointment by UUID
func (r *Repository) GetAppointmentByUUID(ctx context.Context, uuid types.BinaryUUID) (*models.Appointment, error) {
	var appointment models.Appointment
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&appointment).Error; err != nil {
		r.logger.Error("failed to get appointment: ", err)
		return nil, common.HandleDBError(err, AppointmentErrMap)
	}
	return &appointment, nil
}

// ListAppointments lists appointments based on the provided filters
func (r *Repository) ListAppointments(ctx context.Context, query ListAppointmentsQuery, doctorID, patientID *types.BinaryUUID) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	db := r.db.WithContext(ctx)

	if doctorID != nil {
		db = db.Where("doctor_id = ?", doctorID)
	}
	if patientID != nil {
		db = db.Where("patient_id = ?", patientID)
	}
	if len(query.Status) > 0 {
		db = db.Where("status IN ?", query.Status)
	}
	if query.FromDate != nil {
		db = db.Where("appointment_date >= ?", query.FromDate)
	}
	if query.ToDate != nil {
		db = db.Where("appointment_date <= ?", query.ToDate)
	}

	// Get total count
	if err := db.Model(&models.Appointment{}).Count(&total).Error; err != nil {
		r.logger.Error("failed to get appointments count: ", err)
		return nil, 0, common.HandleDBError(err, AppointmentErrMap)
	}

	// Get paginated results
	offset := (query.Page - 1) * query.PerPage
	if err := db.Offset(offset).Limit(query.PerPage).Find(&appointments).Error; err != nil {
		r.logger.Error("failed to list appointments: ", err)
		return nil, 0, common.HandleDBError(err, AppointmentErrMap)
	}

	return appointments, total, nil
}

// ListAppointmentsForDay gets all appointments for a doctor on a specific day
func (r *Repository) ListAppointmentsForDay(ctx context.Context, doctorID types.BinaryUUID, date time.Time) ([]models.Appointment, error) {
	var appointments []models.Appointment

	err := r.db.WithContext(ctx).
		Where("doctor_id = ? AND appointment_date = ?", doctorID, date).
		Find(&appointments).
		Error

	if err != nil {
		r.logger.Error("failed to list appointments: ", err)
		return nil, errorz.Wrap(err, "failed to list appointments")
	}

	return appointments, nil
}

// CheckSlotAvailability checks if a slot is available for booking
func (r *Repository) CheckSlotAvailability(ctx context.Context, doctorID types.BinaryUUID, date time.Time, startTime, endTime time.Time) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("doctor_id = ? AND appointment_date = ? AND status = ?", doctorID, date, models.AppointmentScheduled).
		Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?)",
			endTime, startTime, endTime, startTime, startTime, endTime).
		Count(&count).Error; err != nil {
		r.logger.Error("failed to check slot availability: ", err)
		return false, common.HandleDBError(err, AppointmentErrMap)
	}
	return count == 0, nil
}

// GetOverlappingAppointments returns appointments that overlap with the given time slot
func (r *Repository) GetOverlappingAppointments(ctx context.Context, doctorID types.BinaryUUID, date time.Time, startTime, endTime time.Time) ([]models.Appointment, error) {
	var appointments []models.Appointment
	if err := r.db.WithContext(ctx).
		Where("doctor_id = ? AND appointment_date = ? AND status = ?", doctorID, date, models.AppointmentScheduled).
		Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?)",
			endTime, startTime, endTime, startTime, startTime, endTime).
		Find(&appointments).Error; err != nil {
		r.logger.Error("failed to get overlapping appointments: ", err)
		return nil, common.HandleDBError(err, AppointmentErrMap)
	}
	return appointments, nil
}

// UpdateAppointmentStatus updates the status of an appointment
func (r *Repository) UpdateAppointmentStatus(ctx context.Context, uuid types.BinaryUUID, status models.AppointmentStatus, reason *string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if reason != nil {
		updates["cancellation_reason"] = reason
	}
	if err := r.db.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("uuid = ?", uuid).
		Updates(updates).Error; err != nil {
		r.logger.Error("failed to update appointment status: ", err)
		return common.HandleDBError(err, AppointmentErrMap)
	}
	return nil
}

// AcquireLock attempts to acquire a lock on an appointment slot
func (r *Repository) AcquireLock(ctx context.Context, uuid types.BinaryUUID, lockKey string, expiry time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("uuid = ? AND (lock_key IS NULL OR lock_expiry < ?)", uuid, time.Now()).
		Updates(map[string]interface{}{
			"lock_key":    lockKey,
			"lock_expiry": expiry,
		})
	if result.Error != nil {
		r.logger.Error("failed to acquire lock: ", result.Error)
		return common.HandleDBError(result.Error, AppointmentErrMap)
	}
	if result.RowsAffected == 0 {
		return ErrSlotLocked
	}
	return nil
}

// ReleaseLock releases a lock on an appointment slot
func (r *Repository) ReleaseLock(ctx context.Context, uuid types.BinaryUUID, lockKey string) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("uuid = ? AND lock_key = ?", uuid, lockKey).
		Updates(map[string]interface{}{
			"lock_key":    nil,
			"lock_expiry": nil,
		}).Error; err != nil {
		r.logger.Error("failed to release lock: ", err)
		return common.HandleDBError(err, AppointmentErrMap)
	}
	return nil
}

// SetSlotCooldown sets a cooldown period for a cancelled slot
func (r *Repository) SetSlotCooldown(ctx context.Context, uuid types.BinaryUUID, cooldownUntil time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("uuid = ?", uuid).
		Update("cooldown_until", cooldownUntil).Error; err != nil {
		r.logger.Error("failed to set slot cooldown: ", err)
		return common.HandleDBError(err, AppointmentErrMap)
	}
	return nil
}
