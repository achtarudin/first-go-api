package model

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents the transactions table
type Transaction struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     *uint          `json:"user_id" gorm:"index"`
	MerchantID *uint          `json:"merchant_id" gorm:"index"`
	CourierID  *uint          `json:"courier_id" gorm:"index"`
	FoodID     *uint          `json:"food_id" gorm:"index"`
	Quantity   uint           `json:"quantity" gorm:"not null"`
	TotalPrice uint           `json:"total_price" gorm:"not null"`
	Status     string         `json:"status" gorm:"type:varchar(20);default:pending"`
	CreatedAt  *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Merchant *Merchant `json:"merchant,omitempty" gorm:"foreignKey:MerchantID"`
	Courier  *Courier  `json:"courier,omitempty" gorm:"foreignKey:CourierID"`
	Food     *Food     `json:"food,omitempty" gorm:"foreignKey:FoodID"`
}

// TableName specifies the table name for the Transaction model
func (Transaction) TableName() string {
	return "transactions"
}
