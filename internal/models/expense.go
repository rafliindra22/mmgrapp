package models

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	ID        int      `gorm:"primaryKey" json:"id"`
	UserID    int      `json:"user_id"`
	User      *User    `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	PeriodID  int      `json:"period_id"`
	Period    *Period  `gorm:"foreignKey:PeriodID;references:ID" json:"period,omitempty"`
	AccountID int      `json:"account_id"`
	Account   *Account `gorm:"foreignKey:AccountID;references:ID" json:"account,omitempty"`

	Date        time.Time `json:"date"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
	DeletedBy *int `json:"deleted_by,omitempty"`
}
