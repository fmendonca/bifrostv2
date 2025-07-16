package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	RedisCtx    = context.Background()
)

// Inicializa conexão Redis usando variáveis de ambiente
func InitRedis() {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")

	addr := fmt.Sprintf("%s:%s", host, port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // default DB
	})

	// Testa conexão
	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	}

	log.Printf("Connected to Redis at %s", addr)
}

// Publica ação (start/stop) no canal Redis 'vm-actions'
func PublishAction(uuid string, action string) error {
	message := fmt.Sprintf(`{"uuid":"%s", "action":"%s"}`, uuid, action)

	err := RedisClient.Publish(RedisCtx, "vm-actions", message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish action to Redis channel: %w", err)
	}

	log.Printf("Published to Redis channel 'vm-actions': %s", message)
	return nil
}
