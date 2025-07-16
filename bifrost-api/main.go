package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	InitRedis()
	defer DB.Close()

	http.HandleFunc("/api/v1/vms", VMsHandler)
	http.HandleFunc("/api/v1/vms/", StartVMHandler)
	http.HandleFunc("/api/v1/vms/", StopVMHandler)

	log.Println("Bifrost API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
