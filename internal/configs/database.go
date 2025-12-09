package config

import (
	"fmt"
	"log"
	"mmgrapp/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("moneyapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal konek ke database:", err)
	}

	DB = db
	fmt.Println("✅ Database terkoneksi!")
}

// Migrate migrasi semua tabel dan tampilkan detail success/error per tabel
func Migrate() {
	tables := []struct {
		name  string
		model interface{}
	}{
		{"User", &models.User{}},
		{"Profile", &models.Profile{}},
		{"Account", &models.Account{}},
		{"Period", &models.Period{}},
		{"Income", &models.Income{}},
		{"Expense", &models.Expense{}},
		{"UserOTP", &models.UserOTP{}},
		{"RefreshToken", &models.RefreshToken{}},
	}

	for _, table := range tables {
		err := DB.AutoMigrate(table.model)
		if err != nil {
			fmt.Printf("❌ Gagal migrasi tabel %s: %v\n", table.name, err)
		} else {
			fmt.Printf("✅ Tabel %s berhasil dimigrasi\n", table.name)
		}
	}

	fmt.Println("✅ Semua migrasi selesai!")
}
