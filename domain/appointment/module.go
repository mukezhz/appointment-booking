package appointment

import (
	"go.uber.org/fx"
)

var Module = fx.Module("appointment",
	fx.Options(
		fx.Provide(
			NewRepository,
			NewService,
			NewController,
			NewRoute,
		),
		fx.Invoke(
			RegisterRoutes,
		),
	))
