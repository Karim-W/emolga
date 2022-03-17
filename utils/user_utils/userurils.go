package user_utils

import "github.com/karim-w/emolga/models"

func MapUserByUserByStates(users *[]models.RedisUserEntry) *map[string][]models.RedisUserEntry {
	var usersByState = make(map[string][]models.RedisUserEntry)
	for _, user := range *users {
		if _, ok := usersByState[user.State]; !ok {
			usersByState[user.State] = []models.RedisUserEntry{}
		}
		usersByState[user.State] = append(usersByState[user.State], user)
	}
	return &usersByState
}
