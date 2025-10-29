package main

import (
	"fmt"
	config "mmgrapp/internal/configs"
)

func main() {
	config.ConnectDB()

	fmt.Println("🚀 Server siap dijalankan (belum ada HTTP handler)...")
	// nanti di sini kita tambahkan router Fiber / Gin
}
