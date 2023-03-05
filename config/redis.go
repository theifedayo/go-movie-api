package config

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// SetRedisConfig sets the configuration values for RedisClient.
func SetRedisConfig(config *Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	_, err := RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected successfully to redis server")
}

// GetCache gets the value saved set by a key.
// It takes a key and expiration as parameters and returns the cached result.
func GetCache(key string, expiration time.Duration) (string, error) {
	cacheResult, err := RedisClient.Get(RedisClient.Context(), key).Result()
	if err != nil {
		return "", err
	}

	if cacheResult == "" {
		return "", nil
	}

	return cacheResult, nil
}

// SetCache sets the value assigned to a key.
// It takes a key, value to be saved and expiration as parameters.
func SetCache(key string, value string, expiration time.Duration) error {
	err := RedisClient.Set(RedisClient.Context(), key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}
