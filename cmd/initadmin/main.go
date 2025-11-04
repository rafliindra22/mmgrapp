package main

import (
	"bufio"
	"fmt"
	"log"
	config "mmgrapp/internal/configs"
	"mmgrapp/internal/models"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	// Koneksi ke database
	config.ConnectDB()
	db := config.DB

	reader := bufio.NewReader(os.Stdin)

	// Input Username
	fmt.Print("Masukkan username admin: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username) // hapus whitespace, termasuk \r\n

	// Input Email
	fmt.Print("Masukkan email admin: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email) // hapus whitespace, termasuk \r\n

	// Input Password (tidak terlihat di console)
	fmt.Print("Masukkan password admin: ")
	bytePassword, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	password := string(bytePassword)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Gagal hash password:", err)
	}

	admin := models.User{
		Username:   username,
		Email:      email,
		Password:   string(hashedPassword),
		IsAdmin:    true,
		IsVerified: true,
	}

	// Buat admin jika belum ada (idempotent)
	db.FirstOrCreate(&admin, models.User{Username: admin.Username})

	fmt.Println("✅ Admin berhasil dibuat:", admin.Username)
}
