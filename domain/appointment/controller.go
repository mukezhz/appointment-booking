package appointment

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mukezhz/appointment-booking/domain/models"
	"github.com/mukezhz/appointment-booking/pkg/framework"
	"github.com/mukezhz/appointment-booking/pkg/responses"
)

type PaginationMeta struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Total     int    `json:"total"`
	LastPage  int    `json:"lastPage"`
	From      int    `json:"from"`
	To        int    `json:"to"`
	Path      string `json:"path"`
	FirstPage int    `json:"firstPage"`
}

// Controller handles appointment-related HTTP requests
type Controller struct {
	logger  framework.Logger
	service *Service
}

// NewController creates a new appointment controller
func NewController(
	logger framework.Logger,
	service *Service,
) *Controller {
	return &Controller{
		logger:  logger,
		service: service,
	}
}

// CreateAvailability creates a new availability slot
func (c *Controller) CreateAvailability(ctx *gin.Context) {
	var req CreateAvailabilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	availability := &models.Availability{
		UserID:    ctx.GetUint("user_id"),
		Weekday:   req.Weekday,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := c.service.CreateAvailability(availability); err != nil {
		switch err {
		case ErrInvalidTimeRange, ErrAvailabilityExists:
			responses.HandleValidationError(ctx, c.logger, err)
		default:
			responses.HandleError(ctx, c.logger, errors.New("Failed to create availability"))
		}
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusCreated,
		responses.DetailResponseType[AvailabilityResponse]{
			Message: "Availability created successfully",
			Item: AvailabilityResponse{
				ID:        availability.ID,
				Weekday:   availability.Weekday,
				StartTime: availability.StartTime,
				EndTime:   availability.EndTime,
				CreatedAt: availability.CreatedAt,
				UpdatedAt: availability.UpdatedAt,
			},
		},
	)
}

// GetAvailabilities gets all availability slots for a user
func (c *Controller) GetAvailabilities(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	availabilities, err := c.service.GetAvailabilityByUserID(userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, errors.New("Failed to fetch availabilities"))
		return
	}

	var response []AvailabilityResponse
	for _, a := range availabilities {
		response = append(response, AvailabilityResponse{
			ID:        a.ID,
			Weekday:   a.Weekday,
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		})
	}

	responses.ListResponse(
		ctx,
		http.StatusOK,
		responses.ListResponseType[AvailabilityResponse]{
			Message: "Availabilities fetched successfully",
			Items:   response,
		},
	)
}

// CreateBooking creates a new booking
func (c *Controller) CreateBooking(ctx *gin.Context) {
	var req CreateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04", req.Date+" "+req.Time)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	// Calculate end time (assuming 1 hour duration)
	endTime := startTime.Add(time.Hour)

	booking := &models.Booking{
		UserID:     ctx.GetUint("user_id"),
		GuestName:  req.GuestName,
		GuestEmail: req.GuestEmail,
		Date:       date,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     models.BookingStatusPending,
	}

	if err := c.service.CreateBooking(booking); err != nil {
		switch {
		case errors.Is(err, ErrSlotUnavailable):
			responses.HandleValidationError(ctx, c.logger, err)
		case errors.Is(err, ErrBookingExists):
			responses.HandleValidationError(ctx, c.logger, err)
		default:
			responses.HandleError(ctx, c.logger, errors.New("Failed to create booking"))
		}
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusCreated,
		responses.DetailResponseType[BookingResponse]{
			Message: "Booking created successfully",
			Item: BookingResponse{
				GuestName:  booking.GuestName,
				GuestEmail: booking.GuestEmail,
				Date:       booking.Date,
				StartTime:  booking.StartTime,
				EndTime:    booking.EndTime,
				Status:     BookingStatus(booking.Status), // Convert models.BookingStatus to dto.BookingStatus
				CreatedAt:  booking.CreatedAt,
				UpdatedAt:  booking.UpdatedAt,
			},
		},
	)
}

// GetBooking gets a booking by ID
func (c *Controller) GetBooking(ctx *gin.Context) {
	bookingID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	booking, err := c.service.GetBookingByID(uint(bookingID))
	if err != nil {
		if errors.Is(err, ErrBookingNotFound) {
			responses.HandleError(ctx, c.logger, err)
		} else {
			responses.HandleError(ctx, c.logger, err)
		}
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[BookingResponse]{
			Message: "Booking retrieved successfully",
			Item: BookingResponse{
				GuestName:  booking.GuestName,
				GuestEmail: booking.GuestEmail,
				Date:       booking.Date,
				StartTime:  booking.StartTime,
				EndTime:    booking.EndTime,
				Status:     BookingStatus(booking.Status), // Convert models.BookingStatus to dto.BookingStatus
				CreatedAt:  booking.CreatedAt,
				UpdatedAt:  booking.UpdatedAt,
			},
		},
	)
}

// GetBookings retrieves all bookings for the current user
func (c *Controller) GetBookings(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	userID := ctx.GetUint("user_id")
	if userID == 0 {
		responses.HandleError(ctx, c.logger, errors.New("Unauthorized"))
		return
	}

	bookings, total, err := c.service.GetBookings(userID, page, limit)
	if err != nil {
		responses.HandleError(ctx, c.logger, errors.New("Failed to fetch bookings"))
		return
	}

	var bookingResponses []BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, BookingResponse{
			ID:         booking.ID,
			GuestName:  booking.GuestName,
			GuestEmail: booking.GuestEmail,
			Date:       booking.Date,
			StartTime:  booking.StartTime,
			EndTime:    booking.EndTime,
			Status:     BookingStatus(booking.Status),
			CreatedAt:  booking.CreatedAt,
			UpdatedAt:  booking.UpdatedAt,
		})
	}

	responses.ListResponse(
		ctx,
		http.StatusOK,
		responses.ListResponseType[BookingResponse]{
			Message: "Bookings retrieved successfully",
			Items:   bookingResponses,
			Pagination: responses.PaginationResponseType{
				Total:   int64(total),
				HasNext: int64(page*limit) < int64(total),
			},
		},
	)
}

// UpdateBookingStatus updates the status of a booking
func (c *Controller) UpdateBookingStatus(ctx *gin.Context) {
	bookingID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	var req UpdateBookingStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	if err := c.service.UpdateBookingStatus(uint(bookingID), models.BookingStatus(req.Status)); err != nil {
		if errors.Is(err, ErrBookingNotFound) {
			responses.HandleError(ctx, c.logger, err)
		} else {
			responses.HandleError(ctx, c.logger, err)
		}
		return
	}

	responses.DetailResponse(
		ctx,
		http.StatusOK,
		responses.DetailResponseType[string]{
			Message: "Booking status updated successfully",
		},
	)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
