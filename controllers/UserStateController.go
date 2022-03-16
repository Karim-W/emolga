package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karim-w/emolga/models/commands"
	"github.com/karim-w/emolga/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UserStateController struct {
	logger  *zap.SugaredLogger
	service *services.UsersService
}

func (u *UserStateController) getUserInSessionMapped(ctx *fiber.Ctx) error {
	tid := ctx.GetReqHeaders()["Transactionid"]
	if list, err := u.service.GetUsersInSessionMappedBstates(ctx.Params("session"), tid); err != nil {
		u.logger.Errorf("Error while getting users mapped by state: %v", err)
		return ctx.Status(500).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(list)
	}
	return nil
}

func (u *UserStateController) getUsersInHearingMapped(ctx *fiber.Ctx) error {
	tid := ctx.GetReqHeaders()["Transactionid"]
	var arr []string
	u.logger.Infof("getUsersInHearingMapped")
	if err := ctx.BodyParser(&arr); err != nil {
		u.logger.Infof("hearingIds: %v", arr)
		return ctx.Status(500).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}
	u.logger.Infof("hearingIds: %v", arr)
	mappedHearings := make(map[string]interface{})
	for _, hearingId := range arr {
		if list, err := u.service.GetUsersInHearingMappedBstates(hearingId, tid); err != nil {
			u.logger.Errorf("Error while getting users mapped by state: %v", err)
			return ctx.Status(500).JSON(map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			mappedHearings[hearingId] = list
		}
	}
	ctx.JSON(mappedHearings)
	return nil
}

func (u *UserStateController) setUserState(ctx *fiber.Ctx) error {
	tid := ctx.GetReqHeaders()["Transactionid"]
	var c commands.AdminCommand
	if err := ctx.BodyParser(&c); err != nil {
		u.logger.Errorf("Error while getting user state: %v", err)
		return ctx.Status(500).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}
	if err := u.service.SetStates(&c, tid); err != nil {
		u.logger.Errorf("Error while setting user state: %v", err)
		return ctx.Status(500).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		ctx.Status(202)
	}
	return nil
}

func (u *UserStateController) SetupRoutes(rg *fiber.Router) {
	(*rg).Get("/session/:session", u.getUserInSessionMapped)
	(*rg).Post("/hearing/", u.getUsersInHearingMapped)
	(*rg).Post("", u.setUserState)

}
func UserStateControllerProvider(log *zap.SugaredLogger, service *services.UsersService) *UserStateController {
	return &UserStateController{
		logger:  log,
		service: service,
	}
}

var UserStateControllerModule = fx.Option(fx.Provide(UserStateControllerProvider))
