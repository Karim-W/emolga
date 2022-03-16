package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HearingController struct {
	logger  *zap.SugaredLogger
	service *services.HearingService
}

func (h *HearingController) getUsersInHearing(ctx *fiber.Ctx) error {
	expanded := ctx.Query("expanded")
	tid := ctx.GetReqHeaders()["Transactionid"]
	switch expanded {
	case "true":
		if list, err := h.service.GetExpandedUsersInHearing(ctx.Params("Hearing"), tid); err != nil {
			h.logger.Errorf("Error while getting expanded users in Hearing: %v", err)
			return ctx.Status(500).JSON(map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(list)
		}
	case "userMapping":
		if list, err := h.service.GetUsersMappedByState(ctx.Params("Hearing"), tid); err != nil {
			h.logger.Errorf("Error while getting users mapped by state: %v", err)
			return ctx.Status(500).JSON(map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(list)
		}
	default:
		if list, err := h.service.GetUsersInHearing(ctx.Params("Hearing"), tid); err != nil {
			h.logger.Errorf("Error while getting users in Hearing: %v", err)
			return ctx.Status(500).JSON(map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(list)
		}
	}
	return nil
}

func (h *HearingController) AddPSTNUser(ctx *fiber.Ctx) error {
	tid := ctx.GetReqHeaders()["Transactionid"]
	m := models.PstnUser{}
	if err := ctx.BodyParser(&m); err != nil {
		h.logger.Errorf("Error while parsing body: %v", err)
		return ctx.Status(500).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}
	if len(m.HearingIds) == 0 {
		m.HearingIds = []string{ctx.Params("Hearing")}
	}
	if err := h.service.AddPSTNUser(m, tid); err != nil {
		h.logger.Errorf("Error while adding PSTN user: %v", err)
		return ctx.Status(500).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}
	return nil
}

func (h *HearingController) SetupRoutes(rg *fiber.Router) {
	(*rg).Get("/Hearing/:Hearing", h.getUsersInHearing)
	(*rg).Post("/Hearing/:Hearing/users/pstn", h.AddPSTNUser)

}
func HearingControllerProvider(log *zap.SugaredLogger, h *services.HearingService) *HearingController {
	return &HearingController{
		logger:  log,
		service: h,
	}
}

var HearingControllerModule = fx.Option(fx.Provide(HearingControllerProvider))
