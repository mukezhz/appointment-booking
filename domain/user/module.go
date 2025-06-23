package user

import (
	"github.com/mukezhz/appointment-booking/pkg/middlewares"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewRepository,
		NewService,
		NewController,
		NewRoute,
		fx.Annotate(
			middlewares.NewMockAuthMiddleware,
			fx.As(new(middlewares.AuthMiddleware)),
		),
	),
	fx.Invoke(
		RegisterRoutes,
	),
)
