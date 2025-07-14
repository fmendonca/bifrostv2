package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Payload struct {
	Timestamp string `json:"timestamp"`
	VMs       []VM   `json:"vms"`
}

func VMsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w) // <- adiciona CORS antes de tudo

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

	duration := time.Since(start)
	log.Printf("Handled %s %s in %v", r.Method, r.URL.Path, duration)
}

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

// Função utilitária para habilitar CORS
func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
