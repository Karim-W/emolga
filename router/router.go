package router

import (
	"context"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/karim-w/emolga/controllers"
	_ "github.com/karim-w/emolga/docs"
	"github.com/karim-w/emolga/helpers/redishelper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:4000
// @BasePath /
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
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))
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
