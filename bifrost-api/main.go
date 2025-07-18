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

	// 🔓 Rotas públicas
	http.HandleFunc("/api/v1/agent/register", RegisterHostHandler)
	http.HandleFunc("/api/v1/agent/frontend-key", FrontendKeyHandler)

	// 🔒 Rotas autenticadas
	http.HandleFunc("/api/v1/vms", AuthMiddleware(VMsHandler))
	http.HandleFunc("/api/v1/vms/update", AuthMiddleware(UpdateVMHandler))
	http.HandleFunc("/api/v1/vms/", AuthMiddleware(vmActionRouter))

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
