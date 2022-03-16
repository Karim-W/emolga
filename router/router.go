package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/karim-w/emolga/controllers"
	"github.com/karim-w/emolga/helpers/redishelper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func SetupRoutes(log *zap.SugaredLogger, redis *redishelper.RedisManager,
	aac *controllers.AdminActionsController,
	p *controllers.PresenceController,
	u *controllers.UserStateController,
	s *controllers.SessionController,
	h *controllers.HearingController,
) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:      false,
		ServerHeader: "Fiber",
		AppName:      "Emolga",
	})
	app.Use(recover.New())
	app.Use(cors.New())
	base := app.Group("/api/v1")
	adminActionsGroup := base.Group("/Actions")
	aac.SetupRoutes(&adminActionsGroup)
	presenceGroup := base.Group("/Presence")
	p.SetupRoutes(&presenceGroup)
	userStateGroup := base.Group("/UserStates")
	u.SetupRoutes(&userStateGroup)
	sessionGroup := base.Group("/Sessions")
	s.SetupRoutes(&sessionGroup)
	hearingGroup := base.Group("/Hearings")
	h.SetupRoutes(&hearingGroup)
	app.Listen(":4000")
	return app
}

func registerHooks(lifecycle fx.Lifecycle, ginRouter *fiber.App, logger *zap.SugaredLogger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Initializing server")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Terminating server")
			logger.Sync()
			return nil
		},
	})
}

var Module = fx.Options(fx.Provide(SetupRoutes), fx.Invoke(registerHooks))
