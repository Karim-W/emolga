package repo

import (
	"sync"

	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/utils/user_utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SessionRepo struct {
	logger *zap.SugaredLogger
	client *redishelper.RedisManager
}

func SessionRepoProvider(logger *zap.SugaredLogger, client *redishelper.RedisManager) *SessionRepo {
	return &SessionRepo{
		logger: logger,
		client: client,
	}
}

func (s *SessionRepo) GetUsersInSession(session string, tid string) (*[]string, error) {
	if usersList, err := s.client.GetFromSet(session); err == nil {
		return &usersList, nil
	} else {
		return nil, err
	}
}

func (s *SessionRepo) ExpandUserDetails(usersList *[]string) (*[]models.RedisUserEntry, error) {
	var users []models.RedisUserEntry
	wg := sync.WaitGroup{}
	wg.Add(len(*usersList))
	for _, user := range *usersList {
		go func(user string) {
			defer wg.Done()
			if userEntry, err := s.client.GetUserEntry(user); err == nil {
				users = append(users, *userEntry)
			}
		}(user)
	}
	wg.Wait()
	return &users, nil
}

func (s *SessionRepo) MapUserByState(users *[]models.RedisUserEntry) *map[string][]models.RedisUserEntry {
	return user_utils.MapUserByUserByStates(users)
}

var SessionRepoModule = fx.Option(fx.Provide(SessionRepoProvider))
