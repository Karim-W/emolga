package services

import (
	"sync"

	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HearingService struct {
	logger *zap.SugaredLogger
	redis  *redishelper.RedisManager
}

func (h *HearingService) GetUsersInHearing(hearing string, tid string) (*[]string, error) {
	if usersList, err := h.redis.GetFromSet(hearing); err == nil {
		return &usersList, nil
	} else {
		return nil, err
	}
}

func (h *HearingService) GetExpandedUsersInHearing(hId string, tid string) (*[]models.RedisUserEntry, error) {
	if usersList, err := h.GetUsersInHearing(hId, tid); err == nil {
		var users []models.RedisUserEntry
		wg := sync.WaitGroup{}
		wg.Add(len(*usersList))
		for _, user := range *usersList {
			go func(user string) {
				defer wg.Done()
				if userEntry, err := h.redis.GetUserEntry(user); err == nil {
					users = append(users, *userEntry)
				}
			}(user)
		}
		return &users, nil
	} else {
		return nil, err
	}
}

func (h *HearingService) GetUsersMappedByState(hearing string, tid string) (*map[string][]models.RedisUserEntry, error) {
	if usersMap, err := h.GetExpandedUsersInHearing(hearing, tid); err == nil {
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

// func (h *HearingService) AddPSTNUser(email string , phoneNumber string) error {
// 	return h.redis.AddPSTNUser(email, phoneNumber)
// }

func HearingServiceProvider(log *zap.SugaredLogger, redis *redishelper.RedisManager) *HearingService {
	return &HearingService{
		logger: log,
		redis:  redis,
	}
}

var HearingServiceModule = fx.Option(fx.Provide(HearingServiceProvider))
