package user

import (
	"clean-architecture/pkg/middlewares"

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
