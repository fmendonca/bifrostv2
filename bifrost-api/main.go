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

	// Public route: agent registration
	http.HandleFunc("/api/v1/agent/register", RegisterHostHandler)

	// Protected routes (with auth)
	http.HandleFunc("/api/v1/vms", AuthMiddleware(VMsHandler))
	http.HandleFunc("/api/v1/vms/update", AuthMiddleware(UpdateVMHandler))
	http.HandleFunc("/api/v1/vms/", AuthMiddleware(vmActionRouter))

	log.Println("🚀 Bifrost API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func vmActionRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/start") || strings.HasSuffix(path, "/stop") {
		StartStopHandler(w, r)
	} else {
		http.NotFound(w, r)
	}
}
