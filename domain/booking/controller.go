package booking

import (
	"clean-architecture/domain/common"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
	"clean-architecture/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
	logger  framework.Logger
}

func NewController(service *Service, logger framework.Logger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

// CreateBooking handles public booking creation
func (c *Controller) CreateBooking(ctx *gin.Context) {
	var req CreateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	userID := uint(1) // TODO: Get user ID from URL parameter for public bookings

	booking, err := c.service.CreateBooking(ctx.Request.Context(), userID, &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusCreated, responses.DetailResponseType[*BookingResponse]{
		Item:    booking,
		Message: "Booking created successfully",
	})
}

// GetBookings handles fetching user's bookings (protected route)
func (c *Controller) GetBookings(ctx *gin.Context) {
	userID, ok := utils.AnyToUint(ctx.MustGet(framework.UID))
	if !ok {
		c.logger.Error("Failed to get user ID from context")
		responses.HandleError(ctx, c.logger, common.ErrInvalidUserID)
		return
	}

	bookings, err := c.service.GetUserBookings(ctx.Request.Context(), userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.ListResponse(ctx, http.StatusOK, responses.ListResponseType[*BookingResponse]{
		Items:   bookings,
		Message: "Bookings retrieved successfully",
	})
}
