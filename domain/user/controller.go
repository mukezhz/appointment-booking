package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"clean-architecture/domain/constants"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/responses"
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

func (c *Controller) Register(ctx *gin.Context) {
	var dto RegisterUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	user, err := c.service.Register(
		ctx.Request.Context(),
		dto.Email,
		dto.FirstName,
		dto.LastName,
		dto.FirstNameJa,
		dto.LastNameJa,
		constants.UserRole(dto.Role),
	)
	if err != nil {
		c.logger.Error("failed to register user: ", err)
		responses.HandleError(ctx, c.logger, err)
		return
	}

	response := responses.DetailResponseType[UserResponseDTO]{
		Item:    toUserResponseDTO(user),
		Message: "User registered successfully",
	}
	responses.DetailResponse(ctx, http.StatusCreated, response)
}

func (c *Controller) Login(ctx *gin.Context) {
	var dto LoginUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	c.logger.Info("attempting to verify user credentials for email: ", dto.Email)

	user, err := c.service.VerifyCredentials(ctx.Request.Context(), dto.Email)
	if err != nil {
		c.logger.Error("failed to verify user credentials: ", err)
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[UserResponseDTO]{
		Item:    toUserResponseDTO(user),
		Message: "User logged in successfully",
	})
}

func (c *Controller) GetProfile(ctx *gin.Context) {
	id, err := getUserIDFromContext(ctx)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}
	user, err := c.service.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error("failed to get user profile: ", err)
		responses.HandleError(ctx, c.logger, err)
		return
	}
	responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[UserResponseDTO]{
		Item:    toUserResponseDTO(user),
		Message: "User profile retrieved successfully",
	})
}

func (c *Controller) UpdateProfile(ctx *gin.Context) {
	id, err := getUserIDFromContext(ctx)
	if err != nil {
		responses.HandleError(ctx, c.logger, err)
		return
	}

	var dto UpdateUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		responses.HandleValidationError(ctx, c.logger, err)
		return
	}

	user, err := c.service.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error("failed to get user profile: ", err)
		responses.HandleError(ctx, c.logger, err)
		return
	}

	// Update user fields
	if dto.FirstName != "" {
		user.FirstName = dto.FirstName
	}
	if dto.LastName != "" {
		user.LastName = dto.LastName
	}
	if dto.FirstNameJa != "" {
		user.FirstNameJa = dto.FirstNameJa
	}
	if dto.LastNameJa != "" {
		user.LastNameJa = dto.LastNameJa
	}

	if err := c.service.UpdateUser(ctx.Request.Context(), user); err != nil {
		c.logger.Error("failed to update user profile: ", err)
		responses.HandleError(ctx, c.logger, err)
		return
	}

	responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[UserResponseDTO]{
		Item:    toUserResponseDTO(user),
		Message: "User profile updated successfully",
	})
}
