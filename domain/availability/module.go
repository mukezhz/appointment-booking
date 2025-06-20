package availability

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"availability",
	fx.Provide(
		NewRepository,
		NewService,
		NewController,
		NewRoute,
	),
	fx.Invoke(RegisterRoutes),
)
