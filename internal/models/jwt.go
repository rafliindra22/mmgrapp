package models

import "time"

type RefreshToken struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int       `gorm:"index;not null"`  // index untuk query cepat
	Token     string    `gorm:"unique;not null"` // pastikan token unik
	ExpiresAt time.Time `gorm:"not null"`        // expiry token
	IsRevoked bool      `gorm:"default:false"`   // untuk revoke / logout
	CreatedAt time.Time `gorm:"autoCreateTime"`  // timestamp otomatis
}
