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

	http.HandleFunc("/api/v1/vms", VMsHandler)
	http.HandleFunc("/api/v1/vms/", vmActionRouter)

	log.Println("Bifrost API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func vmActionRouter(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/start") {
		StartVMHandler(w, r)
	} else if strings.HasSuffix(r.URL.Path, "/stop") {
		StopVMHandler(w, r)
	} else {
		http.NotFound(w, r)
	}
}
