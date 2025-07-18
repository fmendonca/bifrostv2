package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {
	addr := fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379"))
	RedisClient = redis.NewClient(&redis.Options{Addr: addr})
	if _, err := RedisClient.Ping(RedisCtx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	}
	log.Printf("Connected to Redis at %s", addr)
}

func PublishAction(channel, uuid, action string) error {
	msg := fmt.Sprintf(`{"uuid":"%s","action":"%s"}`, uuid, action)
	return RedisClient.Publish(RedisCtx, channel, msg).Err()
}
