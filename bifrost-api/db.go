package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var DB *sql.DB

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
}

func InitDB() {
	var err error

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "bifrost")
	sslmode := getEnv("DB_SSLMODE", "disable")

	adminConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		host, port, user, password, sslmode)

	adminDB, err := sql.Open("postgres", adminConnStr)
	if err != nil {
		log.Fatal("Admin DB connection failed:", err)
	}
	defer adminDB.Close()

	_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "42P04" {
			log.Printf("Database %s already exists.", dbname)
		} else {
			log.Fatalf("Failed to create database %s: %v", dbname, err)
		}
	} else {
		log.Printf("Database %s created successfully.", dbname)
	}

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

	err = createTableIfNotExists()
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	log.Println("Table checked/created.")
}

func createTableIfNotExists() error {
	query := `
    CREATE TABLE IF NOT EXISTS vms (
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
        pending_action TEXT
    );`
	_, err := DB.Exec(query)
	return err
}

func InsertOrUpdateVM(vm VM) (string, error) {
	res, err := DB.Exec(`
        INSERT INTO vms 
        (name, uuid, state, cpu_allocation, memory_allocation, disks, interfaces, metadata, timestamp, pending_action)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (uuid) DO UPDATE SET 
            name = EXCLUDED.name,
            state = EXCLUDED.state,
            cpu_allocation = EXCLUDED.cpu_allocation,
            memory_allocation = EXCLUDED.memory_allocation,
            disks = EXCLUDED.disks,
            interfaces = EXCLUDED.interfaces,
            metadata = EXCLUDED.metadata,
            timestamp = EXCLUDED.timestamp
    `, vm.Name, vm.UUID, vm.State, vm.CPUAllocation, vm.MemoryAllocation,
		vm.Disks, vm.Interfaces, vm.Metadata, vm.Timestamp, vm.PendingAction)

	if err != nil {
		return "", err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if rows == 1 {
		return "inserted/updated", nil
	}
	return "unchanged", nil
}

func GetAllVMs() ([]VM, error) {
	rows, err := DB.Query(`
        SELECT name, uuid, state, cpu_allocation, memory_allocation, disks, interfaces, metadata, timestamp, pending_action
        FROM vms ORDER BY timestamp DESC LIMIT 100
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vms []VM
	for rows.Next() {
		var vm VM
		err := rows.Scan(&vm.Name, &vm.UUID, &vm.State, &vm.CPUAllocation, &vm.MemoryAllocation,
			&vm.Disks, &vm.Interfaces, &vm.Metadata, &vm.Timestamp, &vm.PendingAction)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

// ✅ Busca apenas VMs com pending_action definido
func GetPendingActions() ([]VM, error) {
	rows, err := DB.Query(`
        SELECT name, uuid, state, cpu_allocation, memory_allocation, disks, interfaces, metadata, timestamp, pending_action
        FROM vms WHERE pending_action IS NOT NULL
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vms []VM
	for rows.Next() {
		var vm VM
		err := rows.Scan(&vm.Name, &vm.UUID, &vm.State, &vm.CPUAllocation, &vm.MemoryAllocation,
			&vm.Disks, &vm.Interfaces, &vm.Metadata, &vm.Timestamp, &vm.PendingAction)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

// ✅ Atualiza estado + limpa pending_action
func UpdateVMState(uuid string, state string) error {
	_, err := DB.Exec(`
        UPDATE vms 
        SET state = $1, pending_action = NULL, timestamp = NOW()
        WHERE uuid = $2
    `, state, uuid)
	return err
}

// ✅ Marca pending_action (ex.: start, stop)
func MarkPendingAction(uuid string, action string) error {
	_, err := DB.Exec(`
        UPDATE vms 
        SET pending_action = $1
        WHERE uuid = $2
    `, action, uuid)
	return err
}
