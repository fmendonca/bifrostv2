package api

import (
	"net/http"

	"go-libvirt-api/internal/libvirtclient"
	"go-libvirt-api/internal/models"
	"go-libvirt-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitHandler(database *gorm.DB) {
	db = database
}

func StartVM(c *gin.Context) {
	controlVM(c, service.StartVM, "VM started")
}

func StopVM(c *gin.Context) {
	controlVM(c, service.StopVM, "VM stopped")
}

func RebootVM(c *gin.Context) {
	controlVM(c, service.RebootVM, "VM rebooted")
}

func PauseVM(c *gin.Context) {
	controlVM(c, service.PauseVM, "VM paused")
}

func ResumeVM(c *gin.Context) {
	controlVM(c, service.ResumeVM, "VM resumed")
}

func DeleteVM(c *gin.Context) {
	controlVM(c, service.DeleteVM, "VM deleted")
}

func controlVM(c *gin.Context, action func(*libvirtclient.Client, string) error, successMsg string) {
	vmID := c.Param("vm_id")
	var vm models.VM
	if err := db.First(&vm, vmID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	var host models.Host
	if err := db.First(&host, vm.HostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	conn, err := libvirtclient.Connect("qemu+tcp://" + host.Address + "/system")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	if err := action(conn, vm.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": successMsg})
	}
}

func CreateVM(c *gin.Context) {
	hostID := c.Param("id")
	var host models.Host
	if err := db.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	var vm models.VM
	if err := c.ShouldBindJSON(&vm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vm.HostID = host.ID

	conn, err := libvirtclient.Connect("qemu+tcp://" + host.Address + "/system")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	xml := service.GenerateVMXML(vm.Name, vm.Memory, vm.CPU, "/var/lib/libvirt/images/"+vm.Name+".qcow2", vm.Network)
	if err := service.CreateVM(conn.Conn(), xml); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&vm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM created", "vm": vm})
}

func AddHost(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&host).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Host added", "host": host})
}

func ListHosts(c *gin.Context) {
	var hosts []models.Host
	if err := db.Find(&hosts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hosts)
}

func ListVMs(c *gin.Context) {
	hostID := c.Param("id")
	var vms []models.VM
	if err := db.Where("host_id = ?", hostID).Find(&vms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vms)
}

func PingHost(c *gin.Context) {
	hostID := c.Param("id")
	var host models.Host
	if err := db.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	conn, err := libvirtclient.Connect("qemu+tcp://" + host.Address + "/system")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "offline", "error": err.Error()})
		return
	}
	defer conn.Close()

	c.JSON(http.StatusOK, gin.H{"status": "online"})
}

func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if credentials.Username != "admin" || credentials.Password != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := service.GenerateJWT(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
