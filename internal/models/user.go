package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	Username   string `gorm:"uniqueIndex;size:100" json:"username"`
	Email      string `gorm:"uniqueIndex;size:100" json:"email"`
	Password   string `json:"-"` // disembunyikan dari response JSON
	IsActive   bool   `gorm:"default:true" json:"is_active"`
	IsAdmin    bool   `gorm:"default:false" json:"is_admin"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
	DeletedBy *int `json:"deleted_by,omitempty"`

	// Relasi ke Profile
	Profile *Profile `gorm:"foreignKey:UserID;references:ID" json:"profile,omitempty"`
}

type Profile struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	FullName   string `json:"full_name"`
	UserID     int    `json:"user_id"`
	User       *User  `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
	DeletedBy *int `json:"deleted_by,omitempty"`
}

type Period struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	StartDate time.Time      `json:"start_date"`
	EndDate   time.Time      `json:"end_date"`
	IsDefault bool           `gorm:"default:true" json:"is_default"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
	DeletedBy *int `json:"deleted_by,omitempty"`
}

type Account struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`

	UserID int   `json:"user_id"`
	User   *User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
	DeletedBy *int `json:"deleted_by,omitempty"`
}
