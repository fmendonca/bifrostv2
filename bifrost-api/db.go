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
	ID           int    `json:"id"`
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	APIKey       string `json:"api_key"`
	RedisChannel string `json:"redis_channel"`
	Status       string `json:"status"`
	LastSeen     string `json:"last_seen"`
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
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "bifrost"),
		getEnv("DB_SSLMODE", "disable"))

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}
	log.Println("✅ Connected to PostgreSQL")
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
	log.Println("✅ Auto-migrate done")
}

func RegisterHost(name string) (*Host, error) {
	hostUUID := uuid.New().String()
	apiKey := uuid.New().String()
	channel := fmt.Sprintf("vm-actions-%s", name)
	_, err := DB.Exec(`INSERT INTO hosts (name, uuid, api_key, redis_channel)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO UPDATE SET last_seen = NOW(), status = 'active'`,
		name, hostUUID, apiKey, channel)
	if err != nil {
		return nil, err
	}
	return GetHostByName(name)
}

func GetHostByName(name string) (*Host, error) {
	var h Host
	err := DB.QueryRow(`SELECT id, name, uuid, api_key, redis_channel, status, last_seen FROM hosts WHERE name=$1`, name).
		Scan(&h.ID, &h.Name, &h.UUID, &h.APIKey, &h.RedisChannel, &h.Status, &h.LastSeen)
	return &h, err
}

func GetHostByAPIKey(apiKey string) (*Host, error) {
	var h Host
	err := DB.QueryRow(`SELECT id, name, uuid, api_key, redis_channel, status, last_seen FROM hosts WHERE api_key=$1`, apiKey).
		Scan(&h.ID, &h.Name, &h.UUID, &h.APIKey, &h.RedisChannel, &h.Status, &h.LastSeen)
	return &h, err
}

func GetOrCreateFrontendHost() (*Host, error) {
	const frontendName = "bifrost-frontend"
	host, err := GetHostByName(frontendName)
	if err == nil {
		return host, nil
	}
	newHost, err := RegisterHost(frontendName)
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend host: %w", err)
	}
	return newHost, nil
}

func InsertOrUpdateVM(vm VM) (string, error) {
	_, err := DB.Exec(`
		INSERT INTO vms 
		(name, uuid, state, cpu_allocation, memory_allocation, disks, interfaces, metadata, timestamp, pending_action, host_uuid)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (uuid) DO UPDATE SET 
			name=EXCLUDED.name, state=EXCLUDED.state, cpu_allocation=EXCLUDED.cpu_allocation,
			memory_allocation=EXCLUDED.memory_allocation, disks=EXCLUDED.disks, interfaces=EXCLUDED.interfaces,
			metadata=EXCLUDED.metadata, timestamp=EXCLUDED.timestamp, pending_action=EXCLUDED.pending_action,
			host_uuid=EXCLUDED.host_uuid
	`, vm.Name, vm.UUID, vm.State, vm.CPUAllocation, vm.MemoryAllocation,
		vm.Disks, vm.Interfaces, vm.Metadata, vm.Timestamp, vm.PendingAction, vm.HostUUID)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func GetAllVMs() ([]VM, error) {
	rows, err := DB.Query(`SELECT name, uuid, state, cpu_allocation, memory_allocation, disks, interfaces, metadata, timestamp, pending_action, host_uuid FROM vms ORDER BY timestamp DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vms []VM
	for rows.Next() {
		var vm VM
		var pendingAction sql.NullString
		var hostUUID sql.NullString

		err := rows.Scan(&vm.Name, &vm.UUID, &vm.State, &vm.CPUAllocation, &vm.MemoryAllocation,
			&vm.Disks, &vm.Interfaces, &vm.Metadata, &vm.Timestamp, &pendingAction, &hostUUID)
		if err != nil {
			return nil, err
		}

		vm.PendingAction = pendingAction.String
		vm.HostUUID = hostUUID.String

		vms = append(vms, vm)
	}
	return vms, nil
}

func UpdateVMState(uuid string, state string) error {
	_, err := DB.Exec(`UPDATE vms SET state=$1, pending_action=NULL, timestamp=NOW() WHERE uuid=$2`, state, uuid)
	return err
}
