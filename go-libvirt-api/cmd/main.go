package main

import (
	"go-libvirt-api/internal/api"
	"go-libvirt-api/internal/config"
)

func main() {
	db := config.InitDB()
	r := api.SetupRouter(db)
	r.Run(":8080")
}
