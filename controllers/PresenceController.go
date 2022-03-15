package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karim-w/emolga/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PresenceController struct {
	logger  *zap.SugaredLogger
	service *services.PresenceService
}

func (p *PresenceController) logPresence(ctx *fiber.Ctx) error {
	p.logger.Info("Presence Update")

	ctx.Status(200)
	return nil
}

func (p *PresenceController) SetupRoutes(rg *fiber.Router) {
	(*rg).Post("", p.logPresence)
}

func PresenceControllerProvider(log *zap.SugaredLogger, service *services.PresenceService) *PresenceController {
	return &PresenceController{
		logger:  log,
		service: service,
	}
}

var PresenceControllerModule = fx.Option(fx.Provide(PresenceControllerProvider))
