package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes mendaftarkan semua route ke server
func SetupRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	// api := r.Group("/api")
	// {
	// 	api.GET("/users", handlers.GetUsers)
	// 	api.POST("/users", handlers.CreateUser)

	// 	api.GET("/profiles", handlers.GetProfiles)
	// 	api.POST("/profiles", handlers.CreateProfile)
	// }
}
