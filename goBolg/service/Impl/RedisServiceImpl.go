package Impl

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"time"
)

// RedisServiceImpl 实现了 RedisService 接口。
type RedisServiceImpl struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisServiceImpl 创建一个实现 RedisService 接口的新实例。
func NewRedisServiceImpl(client *redis.Client) *RedisServiceImpl {
	return &RedisServiceImpl{client: client} // 注意这里使用正确的结构体名
}

// Set 在 Redis 中存储一个键值对，并设置过期时间
func (r *RedisServiceImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 从 Redis 中获取键对应的值
func (r *RedisServiceImpl) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// Delete one or more keys
func (r *RedisServiceImpl) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

// Increment a key
func (r *RedisServiceImpl) Incr(ctx context.Context, key string, delta int64) (int64, error) {
	return r.client.IncrBy(ctx, key, delta).Result()
}

// Decrement a key
func (r *RedisServiceImpl) Decr(ctx context.Context, key string, delta int64) (int64, error) {
	return r.client.DecrBy(ctx, key, delta).Result()
}

// HSet 在 Redis 的哈希类型中设置字段
func (r *RedisServiceImpl) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	// 执行 HSet 操作，并返回结果
	result, err := r.client.HSet(ctx, key, values...).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// HGet 从 Redis 的哈希类型中获取指定字段的值
func (r *RedisServiceImpl) HGet(ctx context.Context, key string, field string) (string, error) {
	// 使用传入的 contxt 执行 HGet 操作，并返回结果
	result, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// HIncr 在哈希表中对指定字段的值执行递增操作
func (r *RedisServiceImpl) HIncr(ctx context.Context, key string, hashKey string, delta int64) (int64, error) {
	return r.client.HIncrBy(ctx, key, hashKey, delta).Result()
}

// HGetAll 从 Redis 的哈希类型中获取所有字段和值
func (r *RedisServiceImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	// 使用传入的 contxt 执行 HGetAll 操作，并返回结果
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RedisServiceImpl) ZAdd(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return r.client.ZAdd(ctx, key, members...).Result()
}

func (r *RedisServiceImpl) ZIncr(ctx context.Context, key string, value interface{}, score float64) (float64, error) {
	return r.client.ZIncrBy(ctx, key, score, value.(string)).Result()
}

func (r *RedisServiceImpl) ZDecr(ctx context.Context, key string, value interface{}, score float64) (float64, error) {
	return r.client.ZIncrBy(ctx, key, -score, value.(string)).Result()
}

func (r *RedisServiceImpl) ZReverseRangeWithScore(ctx context.Context, key string, start, end int64) (map[interface{}]float64, error) {
	results, err := r.client.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	resultMap := make(map[interface{}]float64)
	for _, z := range results {
		resultMap[z.Member] = z.Score
	}
	return resultMap, nil
}

func (r *RedisServiceImpl) ZScore(ctx context.Context, key string, member string) (float64, error) {
	score, err := r.client.ZScore(ctx, key, member).Result()
	if err != nil {
		return 0, err
	}
	return score, nil
}

func (r *RedisServiceImpl) ZAllScore(ctx context.Context, key string) (map[interface{}]float64, error) {
	return r.ZReverseRangeWithScore(ctx, key, 0, -1)
}

func (r *RedisServiceImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	// Retrieve members from Redis
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *RedisServiceImpl) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.SAdd(ctx, key, values...).Result()
}

func (r *RedisServiceImpl) SIsMember(ctx context.Context, key string, value interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, value).Result()
}

func (r *RedisServiceImpl) SSize(ctx context.Context, key string) (int64, error) {
	return r.client.SCard(ctx, key).Result()
}

func (r *RedisServiceImpl) SRemove(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.SRem(ctx, key, values...).Result()
}

func (r *RedisServiceImpl) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	results, err := r.client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	// Convert results from []string to []interface{}
	stringResults := make([]string, len(results))
	for i, v := range results {
		stringResults[i] = v
	}
	return stringResults, nil
}

func (r *RedisServiceImpl) LSize(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

func (r *RedisServiceImpl) LIndex(ctx context.Context, key string, index int64) (interface{}, error) {
	return r.client.LIndex(ctx, key, index).Result()
}

// LPush prepends one or multiple values to a list.
func (r *RedisServiceImpl) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	// Call LPush which accepts multiple values.
	return r.client.LPush(ctx, key, values...).Result()
}

func (r *RedisServiceImpl) LPushAll(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.LPush(ctx, key, values...).Result()
}

func (r *RedisServiceImpl) LPushAllWithExpire(ctx context.Context, key string, duration time.Duration, values ...interface{}) (int64, error) {
	len, err := r.LPushAll(ctx, key, values...)
	if err != nil {
		return 0, err
	}
	_, err = r.Expire(ctx, key, duration)
	return len, err
}

func (r *RedisServiceImpl) LRemove(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return r.client.LRem(ctx, key, count, value).Result()
}

func (r *RedisServiceImpl) BitGet(ctx context.Context, key string, offset int64) (int64, error) {
	return r.client.GetBit(ctx, key, offset).Result()
}

func (r *RedisServiceImpl) BitCount(ctx context.Context, key string) (int64, error) {
	return r.client.BitCount(ctx, key, nil).Result()
}

func (r *RedisServiceImpl) BitField(ctx context.Context, key string, limit, offset int64) ([]int64, error) {
	return r.client.BitField(ctx, key, "GET", "u1", offset).Result()
}

func (r *RedisServiceImpl) BitGetAll(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *RedisServiceImpl) HyperAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.PFAdd(ctx, key, values...).Result()
}

func (r *RedisServiceImpl) HyperGet(ctx context.Context, keys ...string) (int64, error) {
	return r.client.PFCount(ctx, keys...).Result()
}

func (r *RedisServiceImpl) HyperDel(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisServiceImpl) GeoAdd(ctx context.Context, key string, x, y float64, name string) (int64, error) {
	return r.client.GeoAdd(ctx, key, &redis.GeoLocation{
		Longitude: x,
		Latitude:  y,
		Name:      name,
	}).Result()
}

func (r *RedisServiceImpl) GeoCalculationDistance(ctx context.Context, key, placeOne, placeTwo string) (float64, string, error) {
	result, err := r.client.GeoDist(ctx, key, placeOne, placeTwo, "km").Result()
	if err != nil {
		return 0, "", err
	}

	return result, "km", nil
}

// ZRevRangeWithScore 从有序集合中获取指定范围内的元素及其分数
func (r *RedisServiceImpl) ZRevRangeWithScore(ctx context.Context, key string, start, stop int64) (map[interface{}]float64, error) {
	results, err := r.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	resultMap := make(map[interface{}]float64)
	for _, z := range results {
		resultMap[z.Member] = z.Score
	}
	return resultMap, nil
}
func (r *RedisServiceImpl) GeoGetHash(ctx context.Context, key string, places ...string) ([]string, error) {
	return r.client.GeoHash(ctx, key, places...).Result()
}

func (r *RedisServiceImpl) Expire(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, duration).Result()
}

// ZIncrBy 增加有序集合成员的分数
func (r *RedisServiceImpl) ZIncrBy(ctx context.Context, key string, member string, increment float64) (float64, error) {
	return r.client.ZIncrBy(ctx, key, increment, member).Result()
}

// HDecr 在 Redis 哈希类型中递减字段值
func (r *RedisServiceImpl) HDecr(ctx context.Context, key string, hashKey string, delta int64) (int64, error) {
	return r.client.HIncrBy(ctx, key, hashKey, -delta).Result()
}

type Point struct {
	X, Y float64
}

type Distance struct {
	Value float64
	Unit  string
}
