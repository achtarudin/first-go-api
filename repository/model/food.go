package model

import (
	"time"

	"gorm.io/gorm"
)

// Food represents the foods table
type Food struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null"`
	Price       uint           `json:"price" gorm:"not null"`
	MerchantID  *uint          `json:"merchant_id" gorm:"index"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Merchant     *Merchant     `json:"merchant,omitempty" gorm:"foreignKey:MerchantID"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:FoodID"`
}

// TableName specifies the table name for the Food model
func (Food) TableName() string {
	return "foods"
}
