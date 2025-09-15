package model

import (
	"time"

	"gorm.io/gorm"
)

type UserRole struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `gorm:"index"`
	RoleID    uint `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relasi "belongs to"
	User User `gorm:"foreignKey:UserID"`
	Role Role `gorm:"foreignKey:RoleID"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
