package availability

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (c *Controller) CreateAvailability(ctx *gin.Context) {
	var req CreateAvailabilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	userID := ctx.MustGet(framework.UserIDKey).(uint)

	availability, err := c.service.CreateAvailability(ctx.Request.Context(), userID, &req)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusCreated, responses.DetailResponseType[*AvailabilityResponse]{
		Item:    availability,
		Message: "Availability created successfully",
	})
}

func (c *Controller) GetAvailability(ctx *gin.Context) {
	userID := ctx.MustGet(framework.UserIDKey).(uint)

	availabilities, err := c.service.GetUserAvailability(ctx.Request.Context(), userID)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.ListResponse(ctx, http.StatusOK, responses.ListResponseType[*AvailabilityResponse]{
		Items:   availabilities,
		Message: "Availabilities retrieved successfully",
	})
}
