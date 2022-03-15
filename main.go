package main

import (
	"github.com/karim-w/emolga/controllers"
	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/router"
	"github.com/karim-w/emolga/services"
	"github.com/karim-w/emolga/utils/karimslogger"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		karimslogger.LogsModule,
		services.PodManagerModule,
		redishelper.RedisModule,
		services.AdminActionsServiceModule,
		services.PresenceServiceModule,
		controllers.AdminActionsControllerModule,
		controllers.PresenceControllerModule,
		router.Module,
	)
	defer app.Run()
}
