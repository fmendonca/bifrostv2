package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	InitDB()
	InitRedis()
	defer DB.Close()

	// Rotas principais
	http.HandleFunc("/api/v1/vms", VMsHandler)

	// Rota para atualizar VM (tem que vir ANTES para não colidir)
	http.HandleFunc("/api/v1/vms/update", UpdateVMHandler)

	// Rotas específicas de ação (start/stop)
	http.HandleFunc("/api/v1/vms/", vmActionRouter)

	log.Println("🚀 Bifrost API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func vmActionRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case strings.HasSuffix(path, "/start"):
		StartVMHandler(w, r)
	case strings.HasSuffix(path, "/stop"):
		StopVMHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}
