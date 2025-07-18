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

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-KEY")
}

func FrontendKeyHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	secret := r.URL.Query().Get("secret")
	if secret == "" || secret != getEnv("FRONTEND_BOOT_SECRET", "meuSegredoForte") {
		http.Error(w, "Invalid secret", http.StatusUnauthorized)
		return
	}

	host, err := GetOrCreateFrontendHost()
	if err != nil {
		http.Error(w, "Failed to get frontend key", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"api_key": host.APIKey})
}

func RegisterHostHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	host, err := RegisterHost(req.Name)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(host)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

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
		for _, vm := range payload.VMs {
			vm.HostUUID = hostUUID
			InsertOrUpdateVM(vm)
		}
		fmt.Fprintf(w, "Processed %d VMs", len(payload.VMs))

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

func StartVMHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	uuid := extractUUID(r.URL.Path, "/start")
	if uuid == "" {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	channel := r.Header.Get("X-REDIS-CHANNEL")
	if err := PublishAction(channel, uuid, "start"); err != nil {
		http.Error(w, "Failed to publish", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Start action sent to %s", uuid)
}

func StopVMHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	uuid := extractUUID(r.URL.Path, "/stop")
	if uuid == "" {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	channel := r.Header.Get("X-REDIS-CHANNEL")
	if err := PublishAction(channel, uuid, "stop"); err != nil {
		http.Error(w, "Failed to publish", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Stop action sent to %s", uuid)
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

	err := UpdateVMState(update.UUID, update.Result)
	if err != nil {
		http.Error(w, "Failed to update VM", http.StatusInternalServerError)
		return
	}

	log.Printf("Update recebido: UUID=%s, Action=%s, Result=%s, Timestamp=%s",
		update.UUID, update.Action, update.Result, update.Timestamp)

	fmt.Fprint(w, "Update recebido")
}

func extractUUID(path, suffix string) string {
	if !strings.HasSuffix(path, suffix) {
		return ""
	}
	uuid := strings.TrimSuffix(strings.TrimPrefix(path, "/api/v1/vms/"), suffix)
	return strings.Trim(uuid, "/")
}
