package domain

import (
	"clean-architecture/domain/todo"

	"go.uber.org/fx"
)

var Module = fx.Module("domain",
	fx.Options(
		todo.Module,
	),
)
