package availability

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"clean-architecture/domain/models"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
	"clean-architecture/pkg/types"
	"clean-architecture/pkg/utils"
)

type Controller struct {
	service *Service
	logger  framework.Logger
}

func NewController(
	service *Service,
	logger framework.Logger,
) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

// CreateAvailability godoc
// @Summary Create a new availability slot
// @Description Create a new availability slot for a doctor
// @Tags availability
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body CreateAvailabilityRequest true "Availability details"
// @Success 201 {object} responses.DetailResponseType
// @Failure 400 {object} responses.DetailResponseType
// @Failure 401 {object} responses.DetailResponseType
// @Router /api/v1/availability [post]
func (c *Controller) CreateAvailability(ctx *gin.Context) {
	var req CreateAvailabilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	doctorID := ctx.MustGet(framework.UserIDKey).(types.BinaryUUID)

	availability, err := c.service.CreateAvailability(ctx.Request.Context(), doctorID, &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusCreated, responses.DetailResponseType[AvailabilityDTO]{
		Item:    utils.SafeDeref(c.toAvailabilityDTO(availability)),
		Message: "Availability slot created successfully",
	})
}

// UpdateAvailability godoc
// @Summary Update an availability slot
// @Description Update an existing availability slot for a doctor
// @Tags availability
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param uuid path string true "Availability UUID"
// @Param request body UpdateAvailabilityRequest true "Updated availability details"
// @Success 200 {object} responses.DetailResponseType
// @Failure 400 {object} responses.DetailResponseType
// @Failure 401 {object} responses.DetailResponseType
// @Failure 404 {object} responses.DetailResponseType
// @Router /api/v1/availability/{uuid} [put]
func (c *Controller) UpdateAvailability(ctx *gin.Context) {
	var req UpdateAvailabilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	uuid := types.ParseUUID(ctx.Param("uuid"))

	availability, err := c.service.UpdateAvailability(ctx.Request.Context(), uuid, &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[AvailabilityDTO]{
		Item:    utils.SafeDeref(c.toAvailabilityDTO(availability)),
		Message: "Availability slot updated successfully",
	})
}

// DeleteAvailability godoc
// @Summary Delete an availability slot
// @Description Delete an existing availability slot
// @Tags availability
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param uuid path string true "Availability UUID"
// @Success 200 {object} responses.DetailResponseType
// @Failure 401 {object} responses.DetailResponseType
// @Failure 404 {object} responses.DetailResponseType
// @Router /api/v1/availability/{uuid} [delete]
func (c *Controller) DeleteAvailability(ctx *gin.Context) {
	uuid := types.ParseUUID(ctx.Param("uuid"))

	err := c.service.DeleteAvailability(ctx.Request.Context(), uuid)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.MessageOnlyResponse(ctx, http.StatusOK, "Availability slot deleted successfully")
}

// ListAvailability godoc
// @Summary List available slots
// @Description List available slots for a doctor within a date range
// @Tags availability
// @Accept json
// @Produce json
// @Param doctor_id query string true "Doctor UUID"
// @Param from_date query string true "Start date (YYYY-MM-DD)"
// @Param to_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} responses.DetailResponseType
// @Failure 400 {object} responses.DetailResponseType
// @Router /api/v1/availability [get]
func (c *Controller) ListAvailability(ctx *gin.Context) {
	var query ListAvailabilityQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	availableSlots, err := c.service.ListDoctorAvailability(ctx.Request.Context(), &query)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Convert domain models to DTOs
	result := make([]AvailableSlotsDTO, len(availableSlots))
	for i, slots := range availableSlots {
		result[i] = c.toAvailableSlotsDTO(slots.Date, slots.Slots)
	}

	responses.ListResponse(ctx, http.StatusOK, responses.ListResponseType[AvailableSlotsDTO]{
		Items:   result,
		Message: "Availability slots retrieved successfully",
	})
}

// BookAppointment godoc
// @Summary Book an appointment
// @Description Book an appointment with a doctor
// @Tags appointments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body BookAppointmentRequest true "Appointment details"
// @Success 201 {object} responses.DetailResponseType
// @Failure 400 {object} responses.DetailResponseType
// @Failure 401 {object} responses.DetailResponseType
// @Router /api/v1/appointments [post]
func (c *Controller) BookAppointment(ctx *gin.Context) {
	var req BookAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	patientID := ctx.MustGet(framework.UserIDKey).(types.BinaryUUID)

	appointment, err := c.service.BookAppointment(ctx.Request.Context(), patientID, &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusCreated, responses.DetailResponseType[AppointmentResponse]{
		Item:    utils.SafeDeref(c.toAppointmentResponse(appointment)),
		Message: "Appointment booked successfully",
	})
}

// CancelAppointment godoc
// @Summary Cancel an appointment
// @Description Cancel an existing appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param uuid path string true "Appointment UUID"
// @Param request body CancelAppointmentRequest true "Cancellation details"
// @Success 200 {object} responses.DetailResponseType
// @Failure 400 {object} responses.DetailResponseType
// @Failure 401 {object} responses.DetailResponseType
// @Failure 404 {object} responses.DetailResponseType
// @Router /api/v1/appointments/{uuid} [delete]
func (c *Controller) CancelAppointment(ctx *gin.Context) {
	var req CancelAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	uuid := types.ParseUUID(ctx.Param("uuid"))

	result, err := c.service.CancelAppointment(ctx.Request.Context(), uuid, &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[AppointmentResponse]{
		Item:    utils.SafeDeref(result),
		Message: "Appointment cancelled successfully",
	})
}

// ListAppointments godoc
// @Summary List appointments
// @Description List appointments based on filters
// @Tags appointments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param doctor_id query string false "Doctor UUID"
// @Param patient_id query string false "Patient UUID"
// @Param status query []string false "Appointment status" collectionFormat "multi" Enums(scheduled,cancelled,completed,no_show)
// @Param from_date query string false "Start date (YYYY-MM-DD)"
// @Param to_date query string false "End date (YYYY-MM-DD)"
// @Param page query int true "Page number"
// @Param per_page query int true "Items per page"
// @Success 200 {object} responses.DetailResponseType
// @Failure 400 {object} responses.DetailResponseType
// @Failure 401 {object} responses.DetailResponseType
// @Router /api/v1/appointments [get]
func (c *Controller) ListAppointments(ctx *gin.Context) {
	var query ListAppointmentsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	result, total, err := c.service.ListAppointments(ctx.Request.Context(), &query)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.ListResponse(ctx, http.StatusOK, responses.ListResponseType[AppointmentResponse]{
		Items:   result,
		Message: "Appointments retrieved successfully",
		Pagination: responses.PaginationResponseType{
			Total:   total,
			HasNext: query.Page*query.PerPage < int(total),
		},
	})
}

// toAvailabilityDTO converts a domain model to DTO
func (c *Controller) toAvailabilityDTO(m *models.Availability) *AvailabilityDTO {
	if m == nil {
		return nil
	}
	return &AvailabilityDTO{
		UUID:         m.UUID.String(),
		DoctorID:     m.DoctorID.String(),
		DayOfWeek:    m.DayOfWeek,
		StartDate:    m.StartDate,
		EndDate:      m.EndDate,
		StartTime:    m.StartTime,
		EndTime:      m.EndTime,
		SlotDuration: m.SlotDuration,
		IsRecurring:  m.IsRecurring,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// toAvailableSlotsDTO converts domain slot data to DTO
func (c *Controller) toAvailableSlotsDTO(date time.Time, slots []models.Slot) AvailableSlotsDTO {
	dtoSlots := make([]SlotDTO, len(slots))
	for i, slot := range slots {
		dtoSlots[i] = SlotDTO{
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
			Available: slot.Available,
			Cooldown:  slot.Cooldown,
		}
	}
	return AvailableSlotsDTO{
		Date:  date,
		Slots: dtoSlots,
	}
}

// toSlotDTOs converts a slice of domain slots to DTOs
func (c *Controller) toSlotDTOs(slots []models.Slot) []SlotDTO {
	dtoSlots := make([]SlotDTO, len(slots))
	for i, slot := range slots {
		dtoSlots[i] = SlotDTO{
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
			Available: slot.Available,
			Cooldown:  slot.Cooldown,
		}
	}
	return dtoSlots
}

// toAppointmentResponse converts a domain appointment to response DTO
func (c *Controller) toAppointmentResponse(m *models.Appointment) *AppointmentResponse {
	if m == nil {
		return nil
	}
	return &AppointmentResponse{
		UUID:               m.UUID.String(),
		DoctorID:           m.DoctorID.String(),
		PatientID:          m.PatientID.String(),
		AppointmentDate:    m.AppointmentDate,
		StartTime:          m.StartTime,
		EndTime:            m.EndTime,
		Status:             string(m.Status),
		CancellationReason: m.CancellationReason,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}
