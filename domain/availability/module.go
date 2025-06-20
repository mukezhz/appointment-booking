package availability

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewService,
		NewController,
		NewRoute,
	),
	fx.Invoke(RegisterRoutes),
)
