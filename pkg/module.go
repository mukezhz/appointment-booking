package pkg

import (
	"github.com/mukezhz/appointment-booking/pkg/framework"
	"github.com/mukezhz/appointment-booking/pkg/infrastructure"
	"github.com/mukezhz/appointment-booking/pkg/middlewares"
	"github.com/mukezhz/appointment-booking/pkg/services"

	"go.uber.org/fx"
)

var Module = fx.Module("pkg",
	fx.Options(
		fx.Provide(
			framework.NewEnv,
			framework.GetLogger,
		),
	),
	services.Module,
	infrastructure.Module,
	middlewares.Module,
)
