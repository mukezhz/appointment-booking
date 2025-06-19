package user

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/types"

	"github.com/gin-gonic/gin"
)

func getUserIDFromContext(ctx *gin.Context) (types.BinaryUUID, error) {
	userID, exists := ctx.Get(framework.UID)
	if !exists {
		return types.BinaryUUID{}, ErrUserUnauthorized
	}

	id, err := types.ShouldParseUUID(userID.(string))
	if err != nil {
		return types.BinaryUUID{}, ErrUserInvalidUserID
	}

	return id, nil
}
