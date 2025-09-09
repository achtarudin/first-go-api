package model

import (
	"time"

	"gorm.io/gorm"
)

// Merchant represents the merchants table
type Merchant struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	Address   string         `json:"address" gorm:"type:text"`
	UserID    *uint          `json:"user_id" gorm:"index"`
	CreatedAt *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User         *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Foods        []Food        `json:"foods,omitempty" gorm:"foreignKey:MerchantID"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:MerchantID"`
}

// TableName specifies the table name for the Merchant model
func (Merchant) TableName() string {
	return "merchants"
}
