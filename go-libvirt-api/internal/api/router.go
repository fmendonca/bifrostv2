package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupRouter(database *gorm.DB) *gin.Engine {
	db = database
	r := gin.Default()
	SetupRoutes(r)
	return r
}

func SetupRoutes(r *gin.Engine) {
	// Public route
	r.POST("/login", Login)

	// Protected routes (require JWT)
	auth := r.Group("/")
	auth.Use(JWTAuthMiddleware())

	auth.POST("/hosts", AddHost)
	auth.GET("/hosts", ListHosts)
	auth.GET("/hosts/:id/ping", PingHost)
	auth.GET("/hosts/:id/vms", ListVMs)
	auth.POST("/hosts/:id/vms", CreateVM)
	auth.POST("/vms/:vm_id/start", StartVM)
	auth.POST("/vms/:vm_id/stop", StopVM)
	auth.POST("/vms/:vm_id/reboot", RebootVM)
	auth.POST("/vms/:vm_id/pause", PauseVM)
	auth.POST("/vms/:vm_id/resume", ResumeVM)
	auth.DELETE("/vms/:vm_id", DeleteVM)
}
