package config

import (
	"go-libvirt-api/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=libvirtdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// AutoMigrate
	db.AutoMigrate(&models.Host{}, &models.VM{})

	return db
}
