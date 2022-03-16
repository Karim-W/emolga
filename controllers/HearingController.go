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

// @BasePath /Hearing

// get users in a hearing
// @Summary Get users in a hearing
// @Schemes
// @Description Api to get users in a hearing
// @Tags Hearings
// @Accept json
// @Produce json
// @Param Transactionid header string true "Transactionid"
// @Param hearingId path string true "Hearing Id"
// @Param expanded query string Array "expanded: ture: for list of users +detail || mapped for list of users +detail mapped by their state , anything else for just participant id list"
// @Success 200 {object} []models.RedisUserEntry{}
// @Router /api/v1/hearing/{hearingId} [get]
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
	case "mapped":
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

// @BasePath /Hearing

// Add PSTN User
// @Summary add PSTN User to a given hearing
// @Schemes
// @Description Api to add PSTN User to a given hearing
// @Tags Hearings
// @Accept json
// @Produce json
// @Param Transactionid header string true "Transactionid"
// @Param hearingId path string true "Hearing Id"
// @Param data body models.PstnUser{} true "user entry: only email and phone needed"
// @Success 200 {object} []models.RedisUserEntry{}
// @Router /api/v1/hearing/{hearingId}/users/pstn [post]
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
	} else {
		return ctx.SendStatus(202)
	}
}

func (h *HearingController) SetupRoutes(rg *fiber.Router) {
	(*rg).Get("/:Hearing", h.getUsersInHearing)
	(*rg).Post("/:Hearing/users/pstn", h.AddPSTNUser)

}
func HearingControllerProvider(log *zap.SugaredLogger, h *services.HearingService) *HearingController {
	return &HearingController{
		logger:  log,
		service: h,
	}
}

var HearingControllerModule = fx.Option(fx.Provide(HearingControllerProvider))
