package models

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Name       string
	Address    string
	Port       int
	User       string
	AuthMethod string // ssh, tcp, socket
	Password   string
	SSHKeyPath string
}

type VM struct {
	gorm.Model
	HostID  uint
	Name    string
	CPU     int
	Memory  int // MB
	Disk    int // GB
	Network string
}
