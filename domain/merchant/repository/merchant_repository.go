package repository

import (
	"context"
	"cutbray/first_api/domain/merchant/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type MerchantRepository interface {
	Trx(ctx context.Context, clouseure func(tx *gorm.DB) error) error
	Create(ctx context.Context, merchant *entity.UserMerchant, tx *gorm.DB) error
	FindRoleMerchant(ctx context.Context, roleName model.RoleStatus, tx *gorm.DB) (uint, error)
}

type merchantRepository struct {
	db *infra.Database
}

func NewMerchantRepository(db *infra.Database) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

// Ini bisa jadi method di BaseRepository atau struct DB wrapper Anda
func (c *merchantRepository) Trx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return c.db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (c *merchantRepository) Create(context context.Context, merchant *entity.UserMerchant, tx *gorm.DB) error {

	if tx == nil {
		tx = c.db.DB.WithContext(context)
	}

	user := model.User{
		Name:     merchant.Name,
		Email:    merchant.Email,
		Password: merchant.Password,
	}

	if err := tx.Create(&user).Error; err != nil {
		return err
	}

	merchant.ID = int(user.ID)

	return nil
}

func (c *merchantRepository) FindRoleMerchant(ctx context.Context, roleName model.RoleStatus, tx *gorm.DB) (uint, error) {

	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	var roleModel model.Role
	result := tx.First(&roleModel, "name = ?", roleName)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("role %s not found", roleName)
	}

	if result.Error != nil {
		return 0, result.Error
	}

	return roleModel.ID, nil
}
