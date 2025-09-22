package cache

import (
	"log"
	"os"
	"sonit_server/constant/env"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func GetRedisClient(logger *log.Logger) *redis.Client {
	// Client created
	if redisClient != nil {
		return redisClient
	}

	// Initialize new client
	return newRedisClient(os.Getenv(env.REDIS_PORT), logger)
}

func newRedisClient(address string, logger *log.Logger) *redis.Client {
	// Initialize new client
	var client = redis.NewClient(&redis.Options{
		Addr: address,
		DB:   0, // Default
	})

	// Check if properly created
	if err := client.Ping().Err(); err != nil {
		logger.Println(err)
		return nil
	}

	// Pass value to the global client
	redisClient = client

	return client
}
