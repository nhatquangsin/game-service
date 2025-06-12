package main

import (
	"go.uber.org/fx"

	"github.com/nhatquangsin/game-service/app/api/impl"
	"github.com/nhatquangsin/game-service/infra/config"
	"github.com/nhatquangsin/game-service/infra/repo/database"
	"github.com/nhatquangsin/game-service/infra/repo/repoimpl"
	viperutil "github.com/nhatquangsin/game-service/infra/utils/viper"
)

// svcFXModule represents a FX module for the service.
var svcFXModule = fx.Options(
	config.FXModule,
	repoimpl.FXModule,
	database.FXModule,
	impl.FXModule,
	viperutil.FXModule,
)

func newAPIApp() *fx.App {
	app := fx.New(
		svcFXModule,
		// Add API dependencies here.
	)

	return app
}
