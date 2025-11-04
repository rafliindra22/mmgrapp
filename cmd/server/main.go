package main

import (
	"fmt"
	"log"
	config "mmgrapp/internal/configs"
	"mmgrapp/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	config.ConnectDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	host := config.GetEnv("APP_HOST", "localhost")
	port := config.GetEnv("APP_PORT", "8080")
	address := fmt.Sprintf("%s:%s", host, port)

	fmt.Printf("🚀 Server siap dijalankan di http://%s\n", address)

	if err := r.Run(address); err != nil {
		log.Fatal("❌ Gagal menjalankan server:", err)
	}
}
