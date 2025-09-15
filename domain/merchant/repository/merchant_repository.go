package repository

import (
	"context"
	"cutbray/first_api/domain/merchant/entity"
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
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

// Ini bisa jadi method di BaseRepository atau struct DB wrapper Anda
func (c *merchantRepository) Trx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	// 1. Mulai transaksi
	tx := c.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// 2. Siapkan defer untuk panic recovery
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := fn(tx)
	if err != nil {

		if roolbackError := tx.Rollback().Error; roolbackError != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, roolbackError)
		}

		return err
	}

	// 4. Jika closure berhasil, Commit
	return tx.Commit().Error
}

func (c *merchantRepository) Create(context context.Context, merchant *entity.UserMerchant, tx *gorm.DB) error {

	if tx == nil {
		tx = c.db.WithContext(context)
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
		tx = c.db.WithContext(ctx)
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
