package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Host struct {
	ID           int
	Name         string
	UUID         string
	APIKey       string
	RedisChannel string
	Status       string
	LastSeen     string
}

type VM struct {
	Name             string          `json:"name"`
	UUID             string          `json:"uuid"`
	State            string          `json:"state"`
	CPUAllocation    int             `json:"cpu_allocation"`
	MemoryAllocation int64           `json:"memory_allocation"`
	Disks            json.RawMessage `json:"disks"`
	Interfaces       json.RawMessage `json:"interfaces"`
	Metadata         json.RawMessage `json:"metadata"`
	Timestamp        string          `json:"timestamp"`
	PendingAction    string          `json:"pending_action,omitempty"`
	HostUUID         string          `json:"host_uuid"`
}

func InitDB() {
	var err error

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "bifrost")
	sslmode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}
	log.Println("Connected to PostgreSQL database:", dbname)

	autoMigrate()
}

func autoMigrate() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS hosts (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			uuid UUID UNIQUE NOT NULL,
			api_key TEXT UNIQUE NOT NULL,
			redis_channel TEXT NOT NULL,
			status TEXT DEFAULT 'active',
			last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS vms (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			uuid UUID NOT NULL UNIQUE,
			state TEXT,
			cpu_allocation INTEGER,
			memory_allocation BIGINT,
			disks JSONB,
			interfaces JSONB,
			metadata JSONB,
			timestamp TIMESTAMP WITH TIME ZONE,
			pending_action TEXT,
			host_uuid UUID REFERENCES hosts(uuid)
		);`,
	}
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatal("Auto-migrate failed:", err)
		}
	}
	log.Println("Auto-migrate completed.")
}

func RegisterHost(name string) (*Host, error) {
	id := uuid.New().String()
	key := uuid.New().String()
	channel := fmt.Sprintf("vm-actions-%s", id)
	_, err := DB.Exec(`INSERT INTO hosts (name, uuid, api_key, redis_channel) VALUES ($1, $2, $3, $4)
	ON CONFLICT (name) DO UPDATE SET last_seen = NOW(), status = 'active'`,
		name, id, key, channel)
	if err != nil {
		return nil, err
	}
	return GetHostByName(name)
}

func GetHostByAPIKey(apiKey string) (*Host, error) {
	row := DB.QueryRow(`SELECT id, name, uuid, api_key, redis_channel, status, last_seen FROM hosts WHERE api_key = $1`, apiKey)
	var h Host
	err := row.Scan(&h.ID, &h.Name, &h.UUID, &h.APIKey, &h.RedisChannel, &h.Status, &h.LastSeen)
	return &h, err
}

func GetHostByName(name string) (*Host, error) {
	row := DB.QueryRow(`SELECT id, name, uuid, api_key, redis_channel, status, last_seen FROM hosts WHERE name = $1`, name)
	var h Host
	err := row.Scan(&h.ID, &h.Name, &h.UUID, &h.APIKey, &h.RedisChannel, &h.Status, &h.LastSeen)
	return &h, err
}
