package services

import (
	"errors"

	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/models/commands"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UsersService struct {
	logger         *zap.SugaredLogger
	hearingService *HearingService
	sessionService *SessionService
	db             *redishelper.RedisManager
}

func (u *UsersService) GetUsersInSessionMappedBstates(session string, tid string) (*map[string][]models.RedisUserEntry, error) {
	if userMap, err := u.sessionService.GetUsersMappedByState(session, tid); userMap != nil {
		return userMap, nil
	} else {
		return nil, err
	}
}

func (u *UsersService) GetUsersInHearingMappedBstates(hearingId string, tid string) (*map[string][]models.RedisUserEntry, error) {
	if userMap, err := u.hearingService.GetUsersMappedByState(hearingId, tid); userMap != nil {
		return userMap, nil
	} else {
		return nil, err
	}
}

func (u *UsersService) SetStates(c *commands.AdminCommand, tid string) error { //haha like react get it
	go func() {
		u.db.HandlePublishCommand(c, tid)
	}()
	switch c.AudienceType {
	case "user":
		u.handleGuestandUserTypes(c, tid)
		return nil
	case "session":
		u.handleSessionUserTypeForStateUpdate(c, tid)
		return nil
	case "hearing":
		u.handleHearingUserTypeForStateUpdate(c, tid)
		return nil
	case "guest":
		u.handleGuestandUserTypes(c, tid)
		return nil
	default:
		u.logger.Errorw("Invalid audience type", "audienceType", c.AudienceType)
		return errors.New("invalid audience type")
	}
}

func (u *UsersService) handleGuestandUserTypes(command *commands.AdminCommand, tid string) {
	for _, user := range command.Audience {
		go func(user string) {
			u.db.UpdateUserState(user, command.Data["state"].(string))
		}(user)
	}
}

func (u *UsersService) handleSessionUserTypeForStateUpdate(command *commands.AdminCommand, tid string) {
	for _, session := range command.Audience {
		go func(session string) {
			if list, err := u.db.GetFromSet(session); err == nil {
				for _, user := range list {
					go func(user string) {
						u.db.UpdateUserState(user, command.Data["state"].(string))
					}(user)
				}
			}
		}(session)
	}
}
func (u *UsersService) handleHearingUserTypeForStateUpdate(command *commands.AdminCommand, tid string) {
	for _, hearing := range command.Audience {
		go func(hearing string) {
			if list, err := u.db.GetFromSet(hearing); err == nil {
				for _, user := range list {
					go func(user string) {
						u.db.UpdateUserState(user, command.Data["state"].(string))
					}(user)
				}
			}
		}(hearing)
	}
}

func UserServiceProvider(log *zap.SugaredLogger, db *redishelper.RedisManager, hearingService *HearingService, sessionService *SessionService) *UsersService {
	return &UsersService{
		logger:         log,
		hearingService: hearingService,
		sessionService: sessionService,
		db:             db,
	}
}

var UserServiceModule = fx.Option(fx.Provide(UserServiceProvider))
