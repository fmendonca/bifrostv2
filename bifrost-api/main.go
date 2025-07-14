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
