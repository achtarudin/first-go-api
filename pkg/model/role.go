package model

import (
	"time"

	"gorm.io/gorm"
)

type RoleStatus string

const (
	RoleAdmin    RoleStatus = "admin"
	RoleMerchant RoleStatus = "merchant"
	RoleCourier  RoleStatus = "courier"
	RoleCustomer RoleStatus = "customer"
	RoleNotFound RoleStatus = "not_found"
)

type Role struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(50);not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relasi many-to-many
	Users []*User `gorm:"many2many:user_roles;"`
}

func (Role) TableName() string {
	return "roles"
}
