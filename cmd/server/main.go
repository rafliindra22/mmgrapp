package main

import (
	"fmt"
	"log"
	config "mmgrapp/internal/configs"
	"mmgrapp/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Koneksi database
	config.ConnectDB()

	// 2. Buat router Gin
	r := gin.Default()

	// 3. Setup semua routes
	routes.SetupRoutes(r)

	// 4. Jalankan server di port 8080
	fmt.Println("🚀 Server siap dijalankan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
