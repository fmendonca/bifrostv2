package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

// Publica ação (start/stop) no Redis com TTL (opcional)
func publishActionToRedis(uuid string, action string) error {
	key := fmt.Sprintf("vm:%s:action", uuid)
	value := action

	// Define TTL de 5 minutos para limpar automaticamente (opcional, ajuste se quiser)
	err := RedisClient.Set(RedisCtx, key, value, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to publish action to Redis: %w", err)
	}

	log.Printf("Published action to Redis: %s -> %s", key, value)
	return nil
}
