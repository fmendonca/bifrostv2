package main

import (
	"fmt"
	"log"

	ispn "github.com/infinispan/infinispan-go-client"
)

func PublishActionToInfinispan(vmUUID, action string) error {
	host := getEnv("INFINISPAN_HOST", "127.0.0.1:11222")
	user := getEnv("INFINISPAN_USER", "admin")
	pass := getEnv("INFINISPAN_PASSWORD", "admin")
	cacheName := getEnv("INFINISPAN_CACHE", "bifrost")

	config := ispn.HotRodConfiguration{
		URI:      host,
		User:     user,
		Password: pass,
		Sasl:     "SCRAM-SHA-512", // ou SCRAM-SHA-256 se precisar
	}

	client, err := ispn.NewHotRodClient(config)
	if err != nil {
		return fmt.Errorf("failed to connect to Infinispan: %w", err)
	}
	defer client.Close()

	cache, err := client.GetCache(cacheName)
	if err != nil {
		log.Printf("Cache '%s' does not exist. Creating...", cacheName)
		err = client.CreateCache(cacheName, "<distributed-cache/>")
		if err != nil {
			return fmt.Errorf("failed to create cache: %w", err)
		}
		cache, _ = client.GetCache(cacheName)
	}

	// Publica ação no cache
	err = cache.Put(vmUUID, action)
	if err != nil {
		return fmt.Errorf("failed to put action into cache: %w", err)
	}

	log.Printf("✅ Published %s action for VM %s to Infinispan", action, vmUUID)
	return nil
}

// Se já existe um config.go com getEnv, aqui NÃO precisa redefinir.
