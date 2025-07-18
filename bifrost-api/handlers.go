package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Payload struct {
	Timestamp string `json:"timestamp"`
	VMs       []VM   `json:"vms"`
}

func RegisterHostHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}
	host, err := RegisterHost(req.Name)
	if err != nil {
		log.Printf("DB error registering host: %v", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	log.Printf("✅ Host registered: %s (UUID: %s)", host.Name, host.UUID)
	json.NewEncoder(w).Encode(host)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}
		host, err := GetHostByAPIKey(apiKey)
		if err != nil {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-HOST-UUID", host.UUID)
		r.Header.Set("X-REDIS-CHANNEL", host.RedisChannel)
		next(w, r)
	}
}

func VMsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	start := time.Now()
	hostUUID := r.Header.Get("X-HOST-UUID")

	switch r.Method {
	case http.MethodPost:
		var payload Payload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		count := 0
		for _, vm := range payload.VMs {
			vm.HostUUID = hostUUID
			if _, err := InsertOrUpdateVM(vm); err != nil {
				log.Printf("❌ Error saving VM %s: %v", vm.Name, err)
				continue
			}
			count++
		}
		fmt.Fprintf(w, "Processed %d VMs", count)
		log.Printf("✅ POST /api/v1/vms: %d VMs processed", count)
	case http.MethodGet:
		vms, err := GetAllVMs()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(vms)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	log.Printf("Handled %s /api/v1/vms in %v", r.Method, time.Since(start))
}

func StartStopHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	uuid := parts[4]
	action := parts[5]
	channel := r.Header.Get("X-REDIS-CHANNEL")

	if err := PublishAction(channel, uuid, action); err != nil {
		http.Error(w, "Failed to publish action", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s action sent to %s", action, uuid)
	log.Printf("✅ %s action published for VM %s to channel %s", action, uuid, channel)
}

func UpdateVMHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	var update struct {
		UUID      string `json:"uuid"`
		Action    string `json:"action"`
		Result    string `json:"result"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("ℹ️ Update received: %s %s → %s", update.UUID, update.Action, update.Result)
	if err := UpdateVMState(update.UUID, update.Result); err != nil {
		http.Error(w, "Failed to update VM state", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Update processed")
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-KEY")
}
