package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/karim-w/emolga/models/commands"
	"github.com/karim-w/emolga/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AdminActionsController struct {
	logger  *zap.SugaredLogger
	service *services.AdminActionsService
}

func (a *AdminActionsController) PublishAdminAction(ctx *fiber.Ctx) error {
	t := ctx.GetReqHeaders()
	command := commands.AdminCommand{}
	text := ctx.Body()
	err := json.Unmarshal(text, &command)
	if err != nil {
		a.logger.Errorf("Error while unmarshalling admin command: %v", err)
		return ctx.Status(500).SendString("Error while unmarshalling admin command")
	}
	a.service.PublishCommand(&command, t["Transactionid"])
	a.logger.Infow("Publishing admin command:", command)
	a.logger.Infof("Publish Admin Action Request")
	ctx.JSON(map[string]string{
		"message": "Publish Admin Action Request",
	})
	return nil
}

func (a *AdminActionsController) SetupRoutes(rg *fiber.Router) {
	a.logger.Infof("Setting up admin actions routes")
	(*rg).Post("", a.PublishAdminAction)
	a.logger.Infof("Admin actions routes set up")
}

func NewAdminActionController(logger *zap.SugaredLogger, service *services.AdminActionsService) *AdminActionsController {
	return &AdminActionsController{
		logger:  logger,
		service: service,
	}
}

var AdminActionsControllerModule = fx.Option(fx.Provide(NewAdminActionController))
