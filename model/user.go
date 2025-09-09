package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Merchants    []Merchant    `json:"merchants,omitempty" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
