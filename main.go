package main

import (
	"github.com/karim-w/emolga/clients"
	"github.com/karim-w/emolga/controllers"
	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/repo"
	"github.com/karim-w/emolga/router"
	"github.com/karim-w/emolga/services"
	"github.com/karim-w/emolga/utils/hermes"
	"github.com/karim-w/emolga/utils/karimslogger"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		karimslogger.LogsModule,
		redishelper.RedisModule,
		hermes.Module,
		clients.PineduleClientModule,
		repo.HearingRepoModule,
		repo.SessionRepoModule,
		services.AdminActionsServiceModule,
		services.PresenceServiceModule,
		services.SessionServiceModule,
		services.HearingServiceModule,
		services.UserServiceModule,
		controllers.HearingControllerModule,
		controllers.SessionControllerModule,
		controllers.AdminActionsControllerModule,
		controllers.PresenceControllerModule,
		controllers.UserStateControllerModule,
		router.Module,
	)
	defer app.Run()
}
