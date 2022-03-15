package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/karim-w/emolga/models"
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
	t := ctx.GetReqHeaders()
	pr := models.PresenceUpdate{}
	text := ctx.Body()
	err := json.Unmarshal(text, &pr)
	if err != nil {
		p.logger.Errorf("Error while unmarshalling presence update: %v", err)
		return ctx.Status(500).SendString("Error while unmarshalling presence update")
	}
	p.service.PublishPresence(&pr, t["Transactionid"])
	ctx.Status(202)
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
