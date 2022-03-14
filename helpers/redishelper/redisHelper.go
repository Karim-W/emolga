package redishelper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
	return r.trx.SMembers(r.ctx, key).Result()
}

func (r *RedisManager) MapUsersInRoom(roomId string, users []models.RedisUserEntry, c *commands.AdminCommand) *map[string]commands.AdminCommand {
	if userList, err := r.GetFromSet(roomId); err == nil {
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
	} else {
		r.logger.Error(err)
		return nil
	}
}

func (r *RedisManager) SubToPikaEvents() {
	subscriber := r.sub.Subscribe(r.ctx, "pika_events")
	for {
		msg, err := subscriber.ReceiveMessage(r.ctx)
		if err != nil {
			r.logger.Error(err)
		}
		podUpdate := models.PresenceUpdate{}
		err = json.Unmarshal([]byte(msg.Payload), &podUpdate)
		if err != nil {
			r.logger.Error(err)
		}
		fmt.Println("Got message: " + msg.Payload)
		r.handlePresenceUpdate(&podUpdate)
	}
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\:::: PodHandlers :::://\\//\\//\\//\\//\\//\\//\\//\\//\\//\\

func (r *RedisManager) handlePresenceUpdate(update *models.PresenceUpdate) {

}

//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\:::: DI :::://\\//\\//\\//\\//\\//\\//\\//\\//\\//\\
func NewRedisManager(logger *zap.SugaredLogger) *RedisManager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	sub := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisManager{logger, rdb, sub, context.Background()}
}

var RedisModule = fx.Option(fx.Provide(NewRedisManager))
