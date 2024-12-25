package database

import (
	"LoganXav/sori/configs"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func RedisConnect()error {

	redisDatabase, _ := strconv.Atoi(configs.GetEnv("REDIS_DB"))

	redisClient = redis.NewClient(&redis.Options{
		Addr: configs.GetEnv("REDIS_HOST"),
		Password: configs.GetEnv("REDIS_PASSWORD"),
		DB: redisDatabase,
	})

	err := redisClient.Set(ctx, "test_key", "test_value", 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Get a value from Redis
//
//	params	key	string	key to be fetched
//	return string
func RedisGet(key string) string {
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}

	return val
}

// Set a value to Redis
//
//	params	key	string	key to be implemented
//	params	val	string	value to be implemented
//	params	sec	int	duration in seconds
//	return bool
func RedisSet(key, val string, sec int) bool {
	err := redisClient.Set(ctx, key, val, time.Duration(sec)*time.Second)
	return err == nil
}

// Get a value from Redis, otherwise save new value
//
//	params	key	string	key to be fetched/saved
//	params	val	string	value to be fetched/saved
//	params	sec	int	duration in seconds
//	return string
func RedisGetOrSet(key, val string, sec int) string {
	var v string = RedisGet(key)

	if v == "" {
		RedisSet(key, val, sec)
		v = val
	}
	return v
}