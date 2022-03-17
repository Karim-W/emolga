package common

import (
	"sync"

	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/utils/user_utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UserManager struct {
	logger *zap.SugaredLogger
	client *redishelper.RedisManager
}

func (h *UserManager) GetUsersInSet(id string, tid string) (*[]string, error) {
	h.logger.Infow("GetUsersInSet", "Id:", id, "tid", tid)
	if usersList, err := h.client.GetFromSet(id); err == nil {
		return &usersList, nil
	} else {
		return nil, err
	}
}

func (h *UserManager) ExpandUserDetails(usersList *[]string) (*[]models.RedisUserEntry, error) {
	var users []models.RedisUserEntry
	wg := sync.WaitGroup{}
	wg.Add(len(*usersList))
	for _, user := range *usersList {
		go func(user string) {
			defer wg.Done()
			if userEntry, err := h.client.GetUserEntry(user); err == nil {
				users = append(users, *userEntry)
			}
		}(user)
	}
	wg.Wait()
	return &users, nil
}

func (h *UserManager) MapUserByState(users *[]models.RedisUserEntry) *map[string][]models.RedisUserEntry {
	return user_utils.MapUserByUserByStates(users)
}

func UserMangerProvider(logger *zap.SugaredLogger, client *redishelper.RedisManager) *UserManager {
	return &UserManager{
		logger: logger,
		client: client,
	}
}

var UserManagerModule = fx.Option(fx.Provide(UserMangerProvider))
