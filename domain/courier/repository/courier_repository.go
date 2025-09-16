package repository

import (
	"context"
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CourierRepository interface {
	Trx(ctx context.Context, clouseure func(tx *gorm.DB) error) error
	Create(ctx context.Context, courier *entity.Courier, tx *gorm.DB) (*entity.Courier, error)
	FindRoleCourier(ctx context.Context, roleName model.RoleStatus, tx *gorm.DB) (uint, error)
	FindByEmail(ctx context.Context, courier *entity.Courier, tx *gorm.DB) (*entity.Courier, error)
}

type courierRepository struct {
	db *infra.Database
}

func NewCourierRepository(db *infra.Database) CourierRepository {
	return &courierRepository{
		db: db,
	}
}

// Ini bisa jadi method di BaseRepository atau struct DB wrapper Anda
func (c *courierRepository) Trx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return c.db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (c *courierRepository) Create(ctx context.Context, courier *entity.Courier, tx *gorm.DB) (*entity.Courier, error) {

	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	userModel := model.User{
		Name:     courier.Name,
		Email:    courier.Email,
		Password: courier.Password,
		Roles:    []*model.Role{{ID: uint(courier.RoleID)}},
		Courier:  model.Courier{Phone: courier.Phone},
	}

	result := tx.Create(&userModel)

	if result.Error != nil {
		// Cek error duplikat dari MySQL
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, errors.New("email already exists")
		}
		return nil, result.Error
	}

	createdCourier := &entity.Courier{
		ID:    int(userModel.ID),
		Name:  userModel.Name,
		Email: userModel.Email,
		Phone: userModel.Courier.Phone,
	}
	return createdCourier, nil
}

func (c *courierRepository) FindRoleCourier(ctx context.Context, roleName model.RoleStatus, tx *gorm.DB) (uint, error) {

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

func (c *courierRepository) FindByEmail(ctx context.Context, courier *entity.Courier, tx *gorm.DB) (*entity.Courier, error) {

	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	var userModel model.User

	// Subquery untuk mendapatkan user_id dari user_roles dengan role courier
	subQuery := tx.
		Table("user_roles").
		Select("user_roles.user_id").
		Where("role_id = ?", courier.RoleID)

	// Cari user dengan email dan role courier
	result := tx.Where("email = ?", courier.Email).
		Where("id IN (?)", subQuery).
		Where("id IN (?)", tx.Model(&model.Courier{}).Select("user_id")).
		Preload("Courier").Preload("Roles").First(&userModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	foundCourier := &entity.Courier{
		ID:       int(userModel.ID),
		Name:     userModel.Name,
		Email:    userModel.Email,
		Phone:    userModel.Courier.Phone,
		Password: userModel.Password,
	}

	return foundCourier, nil
}
