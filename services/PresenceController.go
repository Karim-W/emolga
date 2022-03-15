package services

import (
	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/models/commands"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PresenceService struct {
	logger *zap.SugaredLogger
	redis  *redishelper.RedisManager
}

func (p *PresenceService) PublishPresence(presence *models.PresenceUpdate, tid string) {
	p.logger.Infow("Publishing presence update", presence, "Transaction ID", tid)
	dat := make(map[string]interface{})
	if presence.NotificationType == "Added" {
		dat["addedUser"] = []string{presence.UserId}
	} else if presence.NotificationType == "Removed" {
		dat["removedUser"] = []string{presence.UserId}
	}
	command := &commands.AdminCommand{
		Command: "sessionRoasterUpdate",
		Data:    dat,
	}
	sessionsCommand := *command
	hearingCommand := *command
	for _, e := range presence.NotifiedEntities {
		if e.EntityType == "session" {
			sessionsCommand.AudienceType = "session"
			sessionsCommand.Audience = append(sessionsCommand.Audience, e.EntityId)
		} else if e.EntityType == "hearing" {
			hearingCommand.AudienceType = "hearing"
			hearingCommand.Audience = append(hearingCommand.Audience, e.EntityId)
		}
	}
	if len(sessionsCommand.Audience) > 0 {
		p.redis.HandlePublishCommand(&sessionsCommand, tid)
	}
	if len(hearingCommand.Audience) > 0 {
		p.redis.HandlePublishCommand(&hearingCommand, tid)
	}
}

func PresenceServiceProvider(l *zap.SugaredLogger, r *redishelper.RedisManager) *PresenceService {
	return &PresenceService{
		logger: l,
		redis:  r,
	}
}

var PresenceServiceModule = fx.Option(fx.Provide(PresenceServiceProvider))
