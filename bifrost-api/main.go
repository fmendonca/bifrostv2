package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer DB.Close()

	InitRedis()
	defer redisClient.Close()

	http.HandleFunc("/api/v1/vms", VMsHandler)
	http.HandleFunc("/api/v1/vms/action", VMsActionHandler)

	log.Println("Bifrost API running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
