package availability

import (
	"time"
)

// CreateAvailabilityRequest represents the request to create a new availability slot
type CreateAvailabilityRequest struct {
	DayOfWeek    *int8      `json:"day_of_week,omitempty" binding:"omitempty,min=1,max=7"`
	StartDate    *time.Time `json:"start_date,omitempty" binding:"omitempty,required_without=DayOfWeek"`
	EndDate      *time.Time `json:"end_date,omitempty" binding:"omitempty,gtfield=StartDate"`
	StartTime    time.Time  `json:"start_time" binding:"required"`
	EndTime      time.Time  `json:"end_time" binding:"required,gtfield=StartTime"`
	SlotDuration int        `json:"slot_duration" binding:"required,min=15,max=120"`
	IsRecurring  bool       `json:"is_recurring" binding:"required"`
}

// UpdateAvailabilityRequest represents the request to update an existing availability slot
type UpdateAvailabilityRequest struct {
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty" binding:"omitempty,gtfield=StartTime"`
	SlotDuration *int       `json:"slot_duration,omitempty" binding:"omitempty,min=15,max=120"`
}

// ListAvailabilityQuery represents the query parameters for listing availability slots
type ListAvailabilityQuery struct {
	DoctorID string    `form:"doctor_id" binding:"required,uuid4"`
	FromDate time.Time `form:"from_date" binding:"required"`
	ToDate   time.Time `form:"to_date" binding:"required,gtfield=FromDate"`
}

// BookAppointmentRequest represents the request to book an appointment
type BookAppointmentRequest struct {
	DoctorID        string    `json:"doctor_id" binding:"required,uuid4"`
	AppointmentDate time.Time `json:"appointment_date" binding:"required,min_datetime"`
	StartTime       time.Time `json:"start_time" binding:"required"`
}

// CancelAppointmentRequest represents the request to cancel an appointment
type CancelAppointmentRequest struct {
	Reason string `json:"reason" binding:"required,min=10,max=500"`
}

// AppointmentResponse represents the response for appointment operations
type AppointmentResponse struct {
	UUID               string    `json:"uuid"`
	DoctorID           string    `json:"doctor_id"`
	PatientID          string    `json:"patient_id"`
	AppointmentDate    time.Time `json:"appointment_date"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	Status             string    `json:"status"`
	CancellationReason *string   `json:"cancellation_reason,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ListAppointmentsQuery represents the query parameters for listing appointments
type ListAppointmentsQuery struct {
	DoctorID  *string    `form:"doctor_id" binding:"omitempty,uuid4"`
	PatientID *string    `form:"patient_id" binding:"omitempty,uuid4"`
	Status    []string   `form:"status" binding:"omitempty,dive,oneof=scheduled cancelled completed no_show"`
	FromDate  *time.Time `form:"from_date" binding:"omitempty"`
	ToDate    *time.Time `form:"to_date" binding:"omitempty,gtfield=FromDate"`
	Page      int        `form:"page" binding:"required,min=1"`
	PerPage   int        `form:"per_page" binding:"required,min=1,max=100"`
}

// AvailabilityDTO represents the response for availability operations
type AvailabilityDTO struct {
	UUID         string     `json:"uuid"`
	DoctorID     string     `json:"doctor_id"`
	DayOfWeek    *int8      `json:"day_of_week,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      time.Time  `json:"end_time"`
	SlotDuration int        `json:"slot_duration"`
	IsRecurring  bool       `json:"is_recurring"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AvailableSlotsDTO represents the available slots for a doctor
type AvailableSlotsDTO struct {
	Date  time.Time `json:"date"`
	Slots []SlotDTO `json:"slots"`
}

// SlotDTO represents a single available time slot
type SlotDTO struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Available bool      `json:"available"`
	Cooldown  bool      `json:"cooldown"`
}
