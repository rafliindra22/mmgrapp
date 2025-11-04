package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv memuat file .env dari folder configs
func LoadEnv() {
	err := godotenv.Load("internal/configs/.env")
	if err != nil {
		log.Println("⚠️  Tidak menemukan file .env, menggunakan environment default")
	}
}

// GetEnv mengambil variabel dari environment, jika tidak ada gunakan default
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
