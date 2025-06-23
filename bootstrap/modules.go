package bootstrap

import (
	"github.com/mukezhz/appointment-booking/domain"
	"github.com/mukezhz/appointment-booking/pkg"
	"github.com/mukezhz/appointment-booking/seeds"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	pkg.Module,
	domain.Module,
	seeds.Module,
)
