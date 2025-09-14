package repository

import (
	"context"
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/pkg/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CourierRepository interface {
	Trx(ctx context.Context, clouseure func(tx *gorm.DB) error) error
	Create(ctx context.Context, courier *entity.Courier, tx *gorm.DB) error
	FindRoleCourier(ctx context.Context, roleName model.RoleStatus, tx *gorm.DB) (uint, error)
}

type courierRepository struct {
	db *gorm.DB
}

func NewCourierRepository(db *gorm.DB) CourierRepository {
	return &courierRepository{
		db: db,
	}
}

// Ini bisa jadi method di BaseRepository atau struct DB wrapper Anda
func (c *courierRepository) Trx(ctx context.Context, fn func(tx *gorm.DB) error) error {
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

func (c *courierRepository) Create(context context.Context, courier *entity.Courier, tx *gorm.DB) error {

	if tx == nil {
		tx = c.db.WithContext(context)
	}

	user := model.User{
		Name:     courier.Name,
		Email:    courier.Email,
		Password: courier.Password,
		Roles:    []*model.Role{{ID: uint(courier.RoleID)}},
		Courier:  model.Courier{Phone: courier.Phone},
	}

	if err := tx.Create(&user).Error; err != nil {
		return err
	}

	courier.ID = int(user.ID)

	return nil
}

func (c *courierRepository) FindRoleCourier(ctx context.Context, roleName model.RoleStatus, tx *gorm.DB) (uint, error) {

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
