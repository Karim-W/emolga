package redishelper

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/models/commands"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RedisManager struct {
	logger *zap.SugaredLogger
	trx    *redis.Client
	sub    *redis.Client
	ctx    context.Context
}

func (r *RedisManager) AddKeyValuePair(key string, value string) *redis.StatusCmd {
	return r.trx.Set(r.ctx, key, value, time.Hour*24)
}
func (r *RedisManager) GetValueFromKVPair(key string) (string, error) {
	return r.trx.Get(r.ctx, key).Result()
}
func (r *RedisManager) AddToHash(key string, field string, value string) *redis.IntCmd {
	return r.trx.HSet(r.ctx, key, field, value)
}
func (r *RedisManager) GetFromHash(key string, field string) *redis.StringCmd {
	return r.trx.HGet(r.ctx, key, field)
}
func (r *RedisManager) AddToSet(key string, value string) *redis.IntCmd {
	return r.trx.SAdd(r.ctx, key, value)
}
func (r *RedisManager) GetUserEntry(uID string) (*models.RedisUserEntry, error) {
	user := models.RedisUserEntry{}
	if userText, err := r.GetValueFromKVPair(uID); err == nil {
		if err = json.Unmarshal([]byte(userText), &user); err != nil {
			r.logger.Error(err)
			return nil, err
		} else {
			return &user, nil
		}
	}
	return &user, nil
}
func (r *RedisManager) FetchUserPod(userId string) (string, error) {
	if user, err := r.GetUserEntry(userId); err == nil {
		return user.ServerInstance, nil
	} else {
		return "", err
	}
}
func (r *RedisManager) GetFromSet(key string) ([]string, error) {
	return r.sub.SMembers(r.ctx, key).Result()
}

func (r *RedisManager) MapUsersInRoom(userList []string, c *commands.AdminCommand) *map[string]commands.AdminCommand {
	userMap := map[string]commands.AdminCommand{}
	wg := sync.WaitGroup{}
	wg.Add(len(userList))
	for _, user := range userList {
		go func(user string) {
			defer wg.Done()
			if podName, err := r.FetchUserPod(user); err == nil {
				if val, ok := userMap[podName]; ok {
					val.Audience = append(val.Audience, user)
				} else {
					userMap[podName] = commands.AdminCommand{
						Audience:     []string{user},
						AudienceType: "user",
						Command:      c.Command,
						Data:         c.Data,
					}
				}
			}
		}(user)
	}
	wg.Wait()
	return &userMap
}

func (r *RedisManager) AddToZSet(key string, value string) (int64, error) {
	return r.trx.ZAdd(r.ctx, key, &redis.Z{
		Score:  0,
		Member: value,
	}).Result()
}

func (r *RedisManager) UpdateUserState(userId string, state string) error {
	if user, err := r.GetUserEntry(userId); err == nil {
		user.State = state
		if stringifedUser, err := json.Marshal(user); err == nil {
			r.trx.Set(r.ctx, userId, stringifedUser, time.Hour)
			return nil
		} else {
			r.logger.Error(err)
			return err
		}
	} else {
		r.logger.Error(err)
		return err
	}
}

func (r *RedisManager) AddToCommandsHash(key string, identifier string, value commands.AdminCommand) (int64, error) {
	if stringifedCommand, err := json.Marshal(value); err == nil {
		return r.trx.HSet(r.ctx, key+"-Commands", identifier, stringifedCommand).Result()
	} else {
		return 0, err
	}
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\:::: PodHandlers :::://\\//\\//\\//\\//\\//\\//\\//\\//\\//\\
func (r *RedisManager) HandlePublishCommand(c *commands.AdminCommand, tid string) {
	r.logger.Infow("Handling publish command", c)
	switch c.AudienceType {
	case "user":
		r.logger.Infow("Publishing to user")
		mappedUsers := r.MapUsersInRoom(c.Audience, c)
		for k, v := range *mappedUsers {
			go func(k string, v commands.AdminCommand) {
				r.AddToCommandsHash(k, tid, v)
				r.AddToZSet(k, tid)
			}(k, v)
		}
	case "session":
		r.logger.Infow("Publishing to session(s):", c.Audience)
		for _, session := range c.Audience {
			go func(session string) {
				if userList, err := r.GetFromSet(session); err == nil {
					mappedUsers := r.MapUsersInRoom(userList, c)
					for k, v := range *mappedUsers {
						go func(k string, v commands.AdminCommand) {
							r.AddToCommandsHash(k, tid, v)
							r.AddToZSet(k, tid)
						}(k, v)
					}
				}
			}(session)

		}
	case "hearing":
		r.logger.Infow("Publishing to hearing(s):", c.Audience)
		for _, hearing := range c.Audience {
			go func(hearing string) {
				if userList, err := r.GetFromSet(hearing); err == nil {
					mappedUsers := r.MapUsersInRoom(userList, c)
					for k, v := range *mappedUsers {
						go func(k string, v commands.AdminCommand) {
							r.AddToCommandsHash(k, tid, v)
							r.AddToZSet(k, tid)
						}(k, v)
					}
				}
			}(hearing)
		}
	}
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\:::: DI :::://\\//\\//\\//\\//\\//\\//\\//\\//\\//\\
func NewRedisManager(logger *zap.SugaredLogger) *RedisManager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	sub := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdb.Set(context.Background(), "foo", "bar", time.Hour)
	return &RedisManager{logger, rdb, sub, context.Background()}
}

var RedisModule = fx.Option(fx.Provide(NewRedisManager))
