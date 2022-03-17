package services

import (
	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/repo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SessionService struct {
	logger *zap.SugaredLogger
	redis  *redishelper.RedisManager
	repo   *repo.SessionRepo
}

func (s *SessionService) GetUsersInSession(session string, tid string) (*[]string, error) {
	return s.repo.GetUsersInSession(session, tid)
}

func (s *SessionService) GetExpandedUsersInSession(sId string, tid string) (*[]models.RedisUserEntry, error) {
	if usersList, err := s.repo.GetUsersInSession(sId, tid); err == nil {
		return s.repo.ExpandUserDetails(usersList)
	} else {
		return nil, err
	}
}

func (s *SessionService) GetUsersMappedByState(session string, tid string) (*map[string][]models.RedisUserEntry, error) {
	if usersList, err := s.GetExpandedUsersInSession(session, tid); err == nil {
		return s.repo.MapUserByState(usersList), nil
	} else {
		return nil, err
	}
}

func SessionServiceProvider(log *zap.SugaredLogger, redis *redishelper.RedisManager) *SessionService {
	return &SessionService{
		logger: log,
		redis:  redis,
	}
}

var SessionServiceModule = fx.Option(fx.Provide(SessionServiceProvider))
