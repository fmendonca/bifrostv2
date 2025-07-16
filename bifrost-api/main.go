package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	InitDB()
	defer DB.Close()

	// Endpoint principal de listagem e inserção/atualização
	http.HandleFunc("/api/v1/vms", VMsHandler)

	// Endpoint para ações de start/stop por UUID
	http.HandleFunc("/api/v1/vms/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/start") && r.Method == http.MethodPost {
			StartVMHandler(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/stop") && r.Method == http.MethodPost {
			StopVMHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	log.Println("Bifrost API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func vmActionHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if len(path) < len("/api/v1/vms/uuid/stop") {
		http.NotFound(w, r)
		return
	}

	if pathHasSuffix(path, "/start") {
		StartVMHandler(w, r)
		return
	}

	if pathHasSuffix(path, "/stop") {
		StopVMHandler(w, r)
		return
	}

	http.NotFound(w, r)
}

// Helper simples para checar sufixo (evita erro em path estranho)
func pathHasSuffix(path, suffix string) bool {
	if len(path) < len(suffix) {
		return false
	}
	return path[len(path)-len(suffix):] == suffix
}
