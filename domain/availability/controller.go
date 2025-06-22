package availability

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

	ctxUserID := ctx.MustGet(framework.UID)
	userID, ok := utils.AnyToUint(ctxUserID)
	c.logger.Debug("User ID from context:", userID)
	if !ok {
		c.logger.Error("Failed to get user ID from context")
		responses.HandleError(ctx, c.logger, common.ErrInvalidUserID)
		return
	}

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
	userID, ok := utils.AnyToUint(ctx.MustGet(framework.UID))
	if !ok {
		c.logger.Error("Failed to get user ID from context")
		responses.HandleError(ctx, c.logger, common.ErrInvalidUserID)
		return
	}

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
