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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func VMsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	start := time.Now()
	log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodPost:
		handlePostVMs(w, r)
	case http.MethodGet:
		handleGetVMs(w, r)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	log.Printf("Handled %s %s in %v", r.Method, r.URL.Path, time.Since(start))
}

func handlePostVMs(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	count := 0
	for _, vm := range payload.VMs {
		vm.Timestamp = payload.Timestamp
		if action, err := InsertOrUpdateVM(vm); err != nil {
			log.Printf("Failed to insert/update VM %s: %v", vm.Name, err)
		} else {
			log.Printf("%s VM: %s (UUID: %s)", action, vm.Name, vm.UUID)
			count++
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Processed %d VMs", count)
}

func handleGetVMs(w http.ResponseWriter, r *http.Request) {
	var (
		vms []VM
		err error
	)

	if r.URL.Query().Get("pending_action") == "1" {
		vms, err = GetPendingActions()
	} else {
		vms, err = GetAllVMs()
	}

	if err != nil {
		log.Printf("Database error on GET: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(vms); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Response encoding error", http.StatusInternalServerError)
	}
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

	if err := MarkPendingAction(uuid, "start"); err != nil {
		log.Printf("Failed to mark start action for VM %s: %v", uuid, err)
		http.Error(w, "Failed to mark start action", http.StatusInternalServerError)
		return
	}

	if err := publishActionToRedis(uuid, "start"); err != nil {
		log.Printf("Failed to publish start action to Redis: %v", err)
	}

	log.Printf("Marked VM %s for start", uuid)
	fmt.Fprintf(w, "VM %s marked for start", uuid)
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

	if err := MarkPendingAction(uuid, "stop"); err != nil {
		log.Printf("Failed to mark stop action for VM %s: %v", uuid, err)
		http.Error(w, "Failed to mark stop action", http.StatusInternalServerError)
		return
	}

	if err := publishActionToRedis(uuid, "stop"); err != nil {
		log.Printf("Failed to publish stop action to Redis: %v", err)
	}

	log.Printf("Marked VM %s for stop", uuid)
	fmt.Fprintf(w, "VM %s marked for stop", uuid)
}

func extractUUID(path, suffix string) string {
	if !strings.HasSuffix(path, suffix) {
		return ""
	}
	uuid := strings.TrimSuffix(strings.TrimPrefix(path, "/api/v1/vms/"), suffix)
	return strings.Trim(uuid, "/")
}
