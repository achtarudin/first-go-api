package model

import (
	"time"

	"gorm.io/gorm"
)

// Courier represents the couriers table
type Courier struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
	Latitude  *float64       `json:"latitude" gorm:"type:double"`
	Longitude *float64       `json:"longitude" gorm:"type:double"`
	CreatedAt *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:CourierID"`
}

// TableName specifies the table name for the Courier model
func (Courier) TableName() string {
	return "couriers"
}
