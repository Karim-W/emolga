package redishelper

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RedisManager struct {
	logger *zap.SugaredLogger
	client *redis.Client
	ctx    context.Context
}

func (r *RedisManager) AddKeyValuePair(key string, value string) *redis.StatusCmd {
	return r.client.Set(r.ctx, key, value, time.Hour*24)
}
func (r *RedisManager) GetValueFromKVPair(key string) *redis.StringCmd {
	return r.client.Get(r.ctx, key)
}
func (r *RedisManager) AddToHash(key string, field string, value string) *redis.IntCmd {
	return r.client.HSet(r.ctx, key, field, value)
}
func (r *RedisManager) GetFromHash(key string, field string) *redis.StringCmd {
	return r.client.HGet(r.ctx, key, field)
}
func (r *RedisManager) AddToSet(key string, value string) *redis.IntCmd {
	return r.client.SAdd(r.ctx, key, value)
}
func (r *RedisManager) GetFromSet(key string) *redis.StringSliceCmd {
	return r.client.SMembers(r.ctx, key)
}

func NewRedisManager(logger *zap.SugaredLogger) *RedisManager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()
	rdb.Set(ctx, "foo", "bar", time.Hour*24)
	return &RedisManager{logger, rdb, context.Background()}
}

var RedisModule = fx.Option(fx.Provide(NewRedisManager))
