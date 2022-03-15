package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karim-w/emolga/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SessionController struct {
	logger  *zap.SugaredLogger
	service *services.SessionService
}

func (s *SessionController) getUsersInSession(ctx *fiber.Ctx) error {
	expanded := ctx.Query("expanded")
	tid := ctx.GetReqHeaders()["Transactionid"]
	switch expanded {
	case "true":
		if list, err := s.service.GetExpandedUsersInSession(ctx.Params("session"), tid); err != nil {
			s.logger.Errorf("Error while getting expanded users in session: %v", err)
			return ctx.Status(500).JSON(map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(list)
		}
	case "userMapping":
		if list, err := s.service.GetUsersMappedByState(ctx.Params("session"), tid); err != nil {
			s.logger.Errorf("Error while getting users mapped by state: %v", err)
			return ctx.Status(500).JSON(map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(list)
		}
	default:
		if list, err := s.service.GetUsersInSession(ctx.Params("session"), tid); err != nil {
			s.logger.Errorf("Error while getting users in session: %v", err)
			return ctx.Status(500).JSON(map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(list)
		}
	}
	return nil
}

func (s *SessionController) SetupRoutes(rg *fiber.Router) {
	(*rg).Get("/session/:session", s.getUsersInSession)

}
func SessionControllerProvider(log *zap.SugaredLogger, s *services.SessionService) *SessionController {
	return &SessionController{
		logger:  log,
		service: s,
	}
}

var SessionControllerModule = fx.Option(fx.Provide(SessionControllerProvider))
