package services

import (
	"github.com/karim-w/emolga/helpers/redishelper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AdminActionsService struct {
	logger *zap.SugaredLogger
	redis  *redishelper.RedisManager
}

func ProvideAdminActionService(log *zap.SugaredLogger, redis *redishelper.RedisManager) *AdminActionsService {
	return &AdminActionsService{
		logger: log,
		redis:  redis,
	}
}

var AdminActionsServiceModule = fx.Option(fx.Provide(ProvideAdminActionService))
