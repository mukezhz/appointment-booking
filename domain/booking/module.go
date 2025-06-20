package booking

import (
	"clean-architecture/domain/availability"
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ModuleParams struct {
	fx.In

	DB               *infrastructure.Database
	AuthMiddleware   middlewares.AuthMiddleware
	Router           *gin.Engine
	Logger           framework.Logger
	AvailabilityRepo *availability.Repository
}

var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewService,
		NewController,
		NewRouter,
	),
	fx.Invoke(RegisterRoutes),
)
