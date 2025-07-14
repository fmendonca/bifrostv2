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

// Roteador principal
func VMsHandler(w http.ResponseWriter, r *http.Request) {
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

	duration := time.Since(start)
	log.Printf("Handled %s %s in %v", r.Method, r.URL.Path, duration)
}

// POST /api/v1/vms → atualiza inventário
func handlePostVMs(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	count := 0
	for _, vm := range payload.VMs {
		vm.Timestamp = payload.Timestamp
		action, err := InsertOrUpdateVM(vm)
		if err != nil {
			log.Printf("Failed to insert/update VM %s: %v", vm.Name, err)
			continue
		}
		log.Printf("%s VM: %s (UUID: %s) with disks: %s", action, vm.Name, vm.UUID, string(vm.Disks))
		count++
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Processed %d VMs", count)
	log.Printf("POST /api/v1/vms: Processed %d VMs", count)
}

// GET /api/v1/vms → lista VMs
func handleGetVMs(w http.ResponseWriter, r *http.Request) {
	vms, err := GetAllVMs()
	if err != nil {
		log.Printf("Database error on GET: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(vms)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Response encoding error", http.StatusInternalServerError)
		return
	}

	log.Printf("GET /api/v1/vms: Returned %d VMs", len(vms))
}

// POST /api/v1/vms/:uuid/start → atualiza status no banco
func StartVMHandler(w http.ResponseWriter, r *http.Request) {
	uuid := extractUUID(r.URL.Path, "/start")
	if uuid == "" {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	err := UpdateVMState(uuid, "running")
	if err != nil {
		log.Printf("Failed to start VM %s: %v", uuid, err)
		http.Error(w, "Failed to start VM", http.StatusInternalServerError)
		return
	}

	log.Printf("Started VM %s", uuid)
	fmt.Fprintf(w, "VM %s started", uuid)
}

// POST /api/v1/vms/:uuid/stop → atualiza status no banco
func StopVMHandler(w http.ResponseWriter, r *http.Request) {
	uuid := extractUUID(r.URL.Path, "/stop")
	if uuid == "" {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	err := UpdateVMState(uuid, "stopped")
	if err != nil {
		log.Printf("Failed to stop VM %s: %v", uuid, err)
		http.Error(w, "Failed to stop VM", http.StatusInternalServerError)
		return
	}

	log.Printf("Stopped VM %s", uuid)
	fmt.Fprintf(w, "VM %s stopped", uuid)
}

// Helper para extrair UUID da URL
func extractUUID(path, suffix string) string {
	if !strings.HasSuffix(path, suffix) {
		return ""
	}
	uuid := strings.TrimSuffix(strings.TrimPrefix(path, "/api/v1/vms/"), suffix)
	return strings.Trim(uuid, "/")
}
