package responses

import (
	"errors"
	"net/http"

	"github.com/mukezhz/appointment-booking/pkg/errorz"
	"github.com/mukezhz/appointment-booking/pkg/framework"
	"github.com/mukezhz/appointment-booking/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleValidationError(
	ctx *gin.Context,
	logger framework.Logger,
	err error,
) {
	logger.Error(err)
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func HandleErrorWithStatus(
	ctx *gin.Context,
	logger framework.Logger,
	statusCode int,
	err error,
) {
	logger.Error(err)
	ctx.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

func HandleError(
	ctx *gin.Context,
	logger framework.Logger,
	err error,
) {
	msgForUnhandledError := "An error occurred while processing your request. Please try again later."

	var apiErr *errorz.APIError
	msg := err.Error()
	if ok := errors.As(err, &apiErr); ok {
		if msg == "" {
			msg = apiErr.Message
		}
		ctx.JSON(apiErr.StatusCode, gin.H{
			"error": msg,
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": gorm.ErrRecordNotFound.Error(),
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": msgForUnhandledError,
	})

	utils.CurrentSentryService.CaptureException(err)
}
