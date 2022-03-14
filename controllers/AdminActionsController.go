package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karim-w/emolga/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AdminActionsController struct {
	logger  *zap.SugaredLogger
	service *services.AdminActionsService
}

func (a *AdminActionsController) PublishAdminAction(ctx *fiber.Ctx) error {
	a.logger.Infof("Publish Admin Action Request")
	ctx.JSON(map[string]string{
		"message": "Publish Admin Action Request",
	})
}

func (a *AdminActionsController) SetupRoutes(rg fiber.Router) {
	rg.Post("/publishAdminAction", a.PublishAdminAction)
}

func NewAdminActionController(logger *zap.SugaredLogger, service *services.AdminActionsService) *AdminActionsController {
	return &AdminActionsController{
		logger:  logger,
		service: service,
	}
}

var AdminActionsControllerModule = fx.Option(fx.Provide(NewAdminActionController))
