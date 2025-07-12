package models

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Name       string
	Address    string
	Port       int
	User       string
	AuthMethod string
	Password   string
	SSHKeyPath string
}

type VM struct {
	gorm.Model
	HostID  uint
	Name    string
	CPU     int
	Memory  int
	Disk    int
	Network string
}
