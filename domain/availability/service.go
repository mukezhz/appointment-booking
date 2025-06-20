package availability

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"clean-architecture/domain/models"
	"clean-architecture/pkg/errorz"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/types"
)

type Service struct {
	repo   *Repository
	logger framework.Logger
}

func NewService(
	repo *Repository,
	logger framework.Logger,
) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Validate checks if the appointment data is valid
func ValidateAppointment(a *models.Appointment) error {
	// Check if the appointment date is in the future
	if a.AppointmentDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return ErrPastAppointmentDate
	}

	// Check time range
	if a.StartTime.After(a.EndTime) {
		return ErrInvalidTimeRange
	}

	// Check status validity
	switch a.Status {
	case models.AppointmentScheduled, models.AppointmentCancelled, models.AppointmentCompleted, models.AppointmentNoShow:
		// valid status
	default:
		return ErrInvalidAppointmentStatus
	}

	// If cancelled, must have a reason
	if a.Status == models.AppointmentCancelled && (a.CancellationReason == nil || *a.CancellationReason == "") {
		return ErrMissingCancellationReason
	}

	return nil
}

// Validate checks if the availability data is valid
func ValidateAvailability(a *models.Availability) error {
	// Check if either recurring (day_of_week) or specific dates are set
	if (a.DayOfWeek == nil && a.StartDate == nil) || (a.DayOfWeek != nil && a.StartDate != nil) {
		return ErrInvalidAvailabilityConfig
	}

	// For recurring availability
	if a.DayOfWeek != nil {
		if *a.DayOfWeek < 1 || *a.DayOfWeek > 7 {
			return ErrInvalidDayOfWeek
		}
	}

	// Check time range
	if a.StartTime.After(a.EndTime) {
		return ErrInvalidTimeRange
	}

	// Check slot duration
	if a.SlotDuration < 15 || a.SlotDuration > 120 {
		return ErrInvalidSlotDuration
	}

	return nil
}

// CreateAvailability creates a new availability slot for a doctor
func (s *Service) CreateAvailability(ctx context.Context, doctorID types.BinaryUUID, req *CreateAvailabilityRequest) (*models.Availability, error) {
	s.logger.Info("creating availability slot for doctor: ", doctorID.String())

	// Check if doctor exists first
	_, err := s.repo.GetDoctorByID(ctx, doctorID)
	if err != nil {
		return nil, ErrDoctorNotFound
	}

	// Create availability with new UUID
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errorz.NewInternalError("failed to generate uuid")
	}

	availability := &models.Availability{
		UUID:         types.BinaryUUID(id),
		DoctorID:     doctorID,
		DayOfWeek:    req.DayOfWeek,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		SlotDuration: req.SlotDuration,
		IsRecurring:  req.IsRecurring,
	}

	if err := ValidateAvailability(availability); err != nil {
		s.logger.Error("availability validation failed: ", err)
		return nil, err
	}

	// Check for overlapping slots
	existing, err := s.repo.ListDoctorAvailability(ctx, doctorID, time.Now(), time.Now().AddDate(1, 0, 0))
	if err != nil {
		s.logger.Error("failed to check overlapping slots: ", err)
		return nil, err
	}

	for _, slot := range existing {
		if availability.OverlapsWith(&slot) {
			return nil, errorz.NewBadRequestError("time slot overlaps with existing availability")
		}
	}

	if err := s.repo.CreateAvailability(ctx, availability); err != nil {
		s.logger.Error("failed to create availability: ", err)
		return nil, errorz.Wrap(err, "failed to create availability")
	}

	s.logger.Info("created availability slot: ", availability.UUID.String())
	return availability, nil
}

// UpdateAvailability updates an existing availability slot
func (s *Service) UpdateAvailability(ctx context.Context, uuid types.BinaryUUID, req *UpdateAvailabilityRequest) (*models.Availability, error) {
	availability, err := s.repo.GetAvailabilityByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if req.StartTime != nil {
		availability.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		availability.EndTime = *req.EndTime
	}
	if req.SlotDuration != nil {
		availability.SlotDuration = *req.SlotDuration
	}

	if err := ValidateAvailability(availability); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateAvailability(ctx, availability); err != nil {
		return nil, err
	}

	return availability, nil
}

// DeleteAvailability deletes an availability slot
func (s *Service) DeleteAvailability(ctx context.Context, uuid types.BinaryUUID) error {
	return s.repo.DeleteAvailability(ctx, uuid)
}

// ListDoctorAvailability lists available slots for a doctor
func (s *Service) ListDoctorAvailability(ctx context.Context, query *ListAvailabilityQuery) ([]models.AvailableSlots, error) {
	s.logger.Info("listing availability slots for doctor: ", query.DoctorID)

	doctorID := types.ParseUUID(query.DoctorID)
	if _, err := s.repo.GetDoctorByID(ctx, doctorID); err != nil {
		return nil, ErrDoctorNotFound
	}

	availabilities, err := s.repo.ListDoctorAvailability(ctx, doctorID, query.FromDate, query.ToDate)
	if err != nil {
		s.logger.Error("failed to list availabilities: ", err)
		return nil, errorz.NewInternalError("failed to retrieve availability slots")
	}

	// Group slots by date
	slotsByDate := make(map[time.Time][]models.Slot)
	for _, avail := range availabilities {
		slots := s.generateSlots(ctx, &avail, query.FromDate, query.ToDate)
		for date, dateSlots := range slots {
			slotsByDate[date] = append(slotsByDate[date], dateSlots...)
		}
	}

	s.logger.Info(fmt.Sprintf("found %d dates with available slots", len(slotsByDate)))

	// Convert to response format
	var result []models.AvailableSlots
	for date, slots := range slotsByDate {
		result = append(result, models.AvailableSlots{
			Date:  date,
			Slots: slots,
		})
	}

	return result, nil
}

// BookAppointment books an appointment with a doctor
func (s *Service) BookAppointment(ctx context.Context, patientID types.BinaryUUID, req *BookAppointmentRequest) (*models.Appointment, error) {
	doctorID := types.ParseUUID(req.DoctorID)

	// Generate new appointment UUID
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errorz.NewInternalError("failed to generate uuid")
	}

	appointment := &models.Appointment{
		UUID:            types.BinaryUUID(id),
		DoctorID:        doctorID,
		PatientID:       patientID,
		AppointmentDate: req.AppointmentDate,
		StartTime:       req.StartTime,
		Status:          models.AppointmentScheduled,
	}

	// Check for doctor's availability and book the slot
	if err := s.repo.CreateAppointment(ctx, appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

// CancelAppointment cancels an existing appointment
func (s *Service) CancelAppointment(ctx context.Context, uuid types.BinaryUUID, req *CancelAppointmentRequest) (*AppointmentResponse, error) {
	appointment, err := s.repo.GetAppointmentByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if !appointment.CanBeCancelled() {
		return nil, ErrCancellationTooLate
	}

	err = s.repo.UpdateAppointmentStatus(ctx, uuid, models.AppointmentCancelled, &req.Reason)
	if err != nil {
		return nil, err
	}

	// Set cooldown period for the cancelled slot (e.g., 15 minutes)
	cooldownUntil := time.Now().Add(15 * time.Minute)
	err = s.repo.SetSlotCooldown(ctx, uuid, cooldownUntil)
	if err != nil {
		return nil, err
	}

	appointment.Status = models.AppointmentCancelled
	appointment.CancellationReason = &req.Reason
	appointment.CooldownUntil = &cooldownUntil

	return s.toAppointmentResponse(appointment), nil
}

// ListAppointments lists appointments based on the provided filters
func (s *Service) ListAppointments(ctx context.Context, query *ListAppointmentsQuery) ([]AppointmentResponse, int64, error) {
	s.logger.Info("listing appointments with filters: ", query)

	var doctorID, patientID *types.BinaryUUID

	if query.DoctorID != nil {
		id := types.ParseUUID(*query.DoctorID)
		doctorID = &id
		// Verify doctor exists
		if _, err := s.repo.GetDoctorByID(ctx, id); err != nil {
			return nil, 0, ErrDoctorNotFound
		}
	}

	if query.PatientID != nil {
		id := types.ParseUUID(*query.PatientID)
		patientID = &id
		// Note: We could add a similar check for patient existence here if needed
	}

	appointments, total, err := s.repo.ListAppointments(ctx, *query, doctorID, patientID)
	if err != nil {
		s.logger.Error("failed to list appointments: ", err)
		return nil, 0, errorz.NewInternalError("failed to retrieve appointments")
	}

	var response []AppointmentResponse
	for _, apt := range appointments {
		response = append(response, *s.toAppointmentResponse(&apt))
	}

	s.logger.Info(fmt.Sprintf("found %d appointments out of %d total", len(response), total))
	return response, total, nil
}

// Helper functions

func (s *Service) findMatchingAvailability(ctx context.Context, doctorID types.BinaryUUID, date time.Time, startTime time.Time) (*models.Availability, error) {
	availabilities, err := s.repo.ListDoctorAvailability(ctx, doctorID, date, date)
	if err != nil {
		return nil, err
	}

	dayOfWeek := int8(date.Weekday())
	if dayOfWeek == 0 { // Convert Sunday from 0 to 7
		dayOfWeek = 7
	}

	for _, avail := range availabilities {
		if avail.IsRecurring && avail.DayOfWeek != nil && *avail.DayOfWeek == dayOfWeek {
			if s.isTimeInRange(startTime, avail.StartTime, avail.EndTime) {
				return &avail, nil
			}
		} else if !avail.IsRecurring && avail.StartDate != nil && !date.Before(*avail.StartDate) && (avail.EndDate == nil || !date.After(*avail.EndDate)) {
			if s.isTimeInRange(startTime, avail.StartTime, avail.EndTime) {
				return &avail, nil
			}
		}
	}

	return nil, ErrSlotNotAvailable
}

func (s *Service) isTimeInRange(t, start, end time.Time) bool {
	tMinutes := t.Hour()*60 + t.Minute()
	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()
	return tMinutes >= startMinutes && tMinutes < endMinutes
}

// generateSlots generates all available slots for an availability within a date range
func (s *Service) generateSlots(ctx context.Context, avail *models.Availability, fromDate, toDate time.Time) map[time.Time][]models.Slot {
	slots := make(map[time.Time][]models.Slot)

	// For recurring availability
	if avail.DayOfWeek != nil {
		// Find all matching days in the range
		for d := fromDate; !d.After(toDate); d = d.AddDate(0, 0, 1) {
			if int8(d.Weekday()) == *avail.DayOfWeek {
				slots[d] = s.generateDaySlots(ctx, avail, d)
			}
		}
		return slots
	}

	// For non-recurring availability
	if avail.StartDate != nil && avail.EndDate != nil {
		start := fromDate
		if avail.StartDate.After(fromDate) {
			start = *avail.StartDate
		}

		end := toDate
		if avail.EndDate.Before(toDate) {
			end = *avail.EndDate
		}

		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			slots[d] = s.generateDaySlots(ctx, avail, d)
		}
	}

	return slots
}

// generateDaySlots generates slots for a specific day
func (s *Service) generateDaySlots(ctx context.Context, avail *models.Availability, date time.Time) []models.Slot {
	var slots []models.Slot

	// Calculate slots between start and end time
	start := avail.StartTime
	for start.Before(avail.EndTime) {
		end := start.Add(time.Duration(avail.SlotDuration) * time.Minute)
		if end.After(avail.EndTime) {
			break
		}

		slot := models.Slot{
			StartTime: start,
			EndTime:   end,
			Available: true,
			Cooldown:  false,
		}

		slots = append(slots, slot)
		start = end
	}

	// Check appointments and mark slots as unavailable
	appointments, err := s.repo.ListAppointmentsForDay(ctx, avail.DoctorID, date)
	if err != nil {
		s.logger.Error("failed to check appointments: ", err)
		return slots
	}

	for _, apt := range appointments {
		for i := range slots {
			if slots[i].StartTime == apt.StartTime && slots[i].EndTime == apt.EndTime {
				slots[i].Available = false
				if apt.IsInCooldown() {
					slots[i].Cooldown = true
				}
			}
		}
	}

	return slots
}

func (s *Service) toAppointmentResponse(apt *models.Appointment) *AppointmentResponse {
	return &AppointmentResponse{
		UUID:               apt.UUID.String(),
		DoctorID:           apt.DoctorID.String(),
		PatientID:          apt.PatientID.String(),
		AppointmentDate:    apt.AppointmentDate,
		StartTime:          apt.StartTime,
		EndTime:            apt.EndTime,
		Status:             string(apt.Status),
		CancellationReason: apt.CancellationReason,
		CreatedAt:          apt.CreatedAt,
		UpdatedAt:          apt.UpdatedAt,
	}
}
