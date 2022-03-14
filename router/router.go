package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/karim-w/emolga/controllers"
	"github.com/karim-w/emolga/helpers/redishelper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func SetupRoutes(log *zap.SugaredLogger, redis *redishelper.RedisManager, aac *controllers.AdminActionsController) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:      false,
		ServerHeader: "Fiber",
		AppName:      "Emolga",
	})
	go redis.SubToPikaEvents()
	base := app.Group("/api/v1")
	adminActionsGroup := base.Group("/Actions")
	aac.SetupRoutes(adminActionsGroup)
	app.Get("/", func(c *fiber.Ctx) error {
		log.Info("Hello world")
		return c.SendString("Hello, World 👋!")
	})
	app.Listen(":3000")
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
