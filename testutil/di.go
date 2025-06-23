package testutil

import (
	"context"
	"log"

	"github.com/mukezhz/appointment-booking/domain"
	"github.com/mukezhz/appointment-booking/pkg"
	"github.com/mukezhz/appointment-booking/pkg/framework"
	"github.com/mukezhz/appointment-booking/pkg/infrastructure"

	"github.com/onsi/ginkgo/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap/zaptest"
)

func DI(t ginkgo.GinkgoTInterface, opts ...fx.Option) error {
	log.Println("Setting up DI for test...")
	finalOpts := []fx.Option{
		pkg.Module,
		domain.Module,
		fx.Decorate(
			NewTestDatabase,
			func() framework.Logger {
				return framework.Logger{
					SugaredLogger: zaptest.NewLogger(t).Sugar(),
				}
			},
		),
		fx.Invoke(func(db infrastructure.Database) {
			log.Println("Running migration...")
			db.RunMigration()
		}),
	}
	if len(opts) > 0 {
		finalOpts = append(finalOpts, opts...)
	}

	app := fxtest.New(t, finalOpts...)

	ctx := context.Background()
	return app.Start(ctx)
}
