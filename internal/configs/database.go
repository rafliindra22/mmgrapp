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

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Account{},
		&models.Period{},
		&models.Income{},
		&models.Expense{},
	)
	if err != nil {
		log.Fatal("❌ Gagal migrate database:", err)
	}

	fmt.Println("✅ Migrasi database sukses!")
}
