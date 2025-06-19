package domain

import (
	"clean-architecture/domain/todo"
	"clean-architecture/domain/user"

	"go.uber.org/fx"
)

var Module = fx.Module("domain",
	fx.Options(
		todo.Module,
		user.Module,
	),
)
