package domain

import (
	"github.com/mukezhz/appointment-booking/domain/todo"
	"github.com/mukezhz/appointment-booking/domain/user"

	"go.uber.org/fx"
)

var Module = fx.Module("domain",
	fx.Options(
		todo.Module,
		user.Module,
	),
)
