package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	infinispanURL  = getEnv("INFINISPAN_URL", "http://localhost:11222/rest/v2/caches/vm-actions")
	infinispanUser = getEnv("INFINISPAN_USER", "user")
	infinispanPass = getEnv("INFINISPAN_PASS", "pass")
)

// Publica ação no Infinispan via REST API
func publishActionToInfinispan(uuid, action string) error {
	jsonStr := []byte(fmt.Sprintf(`{"uuid":"%s","action":"%s"}`, uuid, action))
	req, err := http.NewRequest("POST", infinispanURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.SetBasicAuth(infinispanUser, infinispanPass)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ Published action '%s' for VM %s to Infinispan", action, uuid)
		return nil
	}

	return fmt.Errorf("infinispan responded with status %d", resp.StatusCode)
}
