package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer DB.Close()

	http.HandleFunc("/api/v1/vms", VMsHandler)

	log.Println("Bifrost API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
