package api

import (
	"go-libvirt-api/internal/libvirt"
	"go-libvirt-api/internal/models"
	"go-libvirt-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartVM(c *gin.Context) {
	vmID := c.Param("vm_id")
	var vm models.VM
	if err := db.First(&vm, vmID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	host, err := getHostByVM(vm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	conn, err := libvirt.ConnectTCPTLS(host.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Disconnect()

	if err := service.StartVM(conn, vm.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "VM started"})
	}
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

func controlVM(c *gin.Context, action func(*libvirt.Libvirt, string) error, successMsg string) {
	vmID := c.Param("vm_id")
	var vm models.VM
	if err := db.First(&vm, vmID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "VM not found"})
		return
	}

	host, err := getHostByVM(vm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	conn, err := libvirt.ConnectTCPTLS(host.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Disconnect()

	if err := action(conn, vm.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": successMsg})
	}
}

func getHostByVM(vm models.VM) (models.Host, error) {
	var host models.Host
	if err := db.First(&host, vm.HostID).Error; err != nil {
		return host, err
	}
	return host, nil
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

	conn, err := libvirt.ConnectTCPTLS(host.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Disconnect()

	xml := service.GenerateVMXML(vm.Name, vm.Memory, vm.CPU, "/var/lib/libvirt/images/"+vm.Name+".qcow2", vm.Network)
	if err := service.CreateVM(conn, xml); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&vm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VM created", "vm": vm})
}

func PingHost(c *gin.Context) {
	hostID := c.Param("id")
	var host models.Host
	if err := db.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	conn, err := libvirt.ConnectTCPTLS(host.Address)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "offline", "error": err.Error()})
		return
	}
	defer conn.Disconnect()

	c.JSON(http.StatusOK, gin.H{"status": "online"})
}
