package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	InitRedis()
	defer DB.Close()

	// Rotas principais
	http.HandleFunc("/api/v1/vms", VMsHandler)

	// Rotas específicas de ação (não registramos duplo /api/v1/vms/)
	http.HandleFunc("/api/v1/vms/", vmActionRouter)

	log.Println("Bifrost API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func vmActionRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case hasSuffix(path, "/start"):
		StartVMHandler(w, r)
	case hasSuffix(path, "/stop"):
		StopVMHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func hasSuffix(path, suffix string) bool {
	return len(path) > len("/api/v1/vms/") && path[len(path)-len(suffix):] == suffix
}
