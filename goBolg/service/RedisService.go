package service

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	//Expire(ctx contxt.Context, key string, duration time.Duration) (bool, error)
	//GetExpire(ctx contxt.Context, key string) (time.Duration, error)
	//HasKey(ctx contxt.Context, key string) (bool, error)
	Incr(ctx context.Context, key string, delta int64) (int64, error)
	Decr(ctx context.Context, key string, delta int64) (int64, error)
	HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HIncr(ctx context.Context, key string, hashKey string, delta int64) (int64, error)
	//HDel(ctx contxt.Context, key string, fields ...string) (int64, error)
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, value interface{}) (bool, error)
	SRemove(ctx context.Context, key string, values ...interface{}) (int64, error)
	ZAdd(ctx context.Context, key string, members ...*redis.Z) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
	ZRevRangeWithScore(ctx context.Context, key string, start, stop int64) (map[interface{}]float64, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	//SMembersLRem(ctx contxt.Context, key string, count int64, value interface{}) (int64, error)
	ZAllScore(ctx context.Context, key string) (map[interface{}]float64, error)
	ZIncrBy(ctx context.Context, key string, member string, increment float64) (float64, error)
	HDecr(ctx context.Context, key string, hashKey string, delta int64) (int64, error)
}
