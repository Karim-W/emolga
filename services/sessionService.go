package services

import (
	"sync"

	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SessionService struct {
	logger *zap.SugaredLogger
	redis  *redishelper.RedisManager
}

func (s *SessionService) GetUsersInSession(session string, tid string) (*[]string, error) {
	if usersList, err := s.redis.GetFromSet(session); err == nil {
		return &usersList, nil
	} else {
		return nil, err
	}
}

func (s *SessionService) GetExpandedUsersInSession(sId string, tid string) (*[]models.RedisUserEntry, error) {
	usersList, err := s.GetUsersInSession(sId, tid)
	if err != nil {
		return nil, err
	}
	var users []models.RedisUserEntry
	wg := sync.WaitGroup{}
	wg.Add(len(*usersList))
	for _, user := range *usersList {
		go func(user string) {
			defer wg.Done()
			if userEntry, err := s.redis.GetUserEntry(user); err == nil {
				users = append(users, *userEntry)
			}
		}(user)
	}
	return &users, nil
}

func (s *SessionService) GetUsersMappedByState(session string, tid string) (*map[string][]models.RedisUserEntry, error) {
	if usersMap, err := s.GetExpandedUsersInSession(session, tid); err == nil {
		var usersByState = make(map[string][]models.RedisUserEntry)
		for _, user := range *usersMap {
			go func(user models.RedisUserEntry) {
				if _, ok := usersByState[user.State]; !ok {
					usersByState[user.State] = []models.RedisUserEntry{}
				}
				usersByState[user.State] = append(usersByState[user.State], user)
			}(user)
		}
		return &usersByState, nil
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
