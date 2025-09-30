package repository

import (
	"context"
	"cutbray/first_api/domain/merchant/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/customerror"
	"cutbray/first_api/pkg/model"
	"cutbray/first_api/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

type MerchantRepository interface {
	Trx(ctx context.Context, clouseure func(tx *gorm.DB) error) error
	Create(ctx context.Context, merchant *entity.UserMerchantRegister, tx *gorm.DB) (*entity.UserMerchant, error)
	FindRoleMerchant(ctx context.Context, tx *gorm.DB) (uint, error)
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

func (c *merchantRepository) Create(context context.Context, merchant *entity.UserMerchantRegister, tx *gorm.DB) (*entity.UserMerchant, error) {

	if tx == nil {
		tx = c.db.DB.WithContext(context)
	}

	// count the email
	var count int64
	err := tx.Model(&model.User{}).Where("email = ?", merchant.Email).Count(&count).Error
	if err != nil {
		return nil, customerror.New(customerror.CodeInternal, "Internal Server Error",
			nil, err)
	}
	if count > 0 {
		return nil, customerror.New(customerror.CodeConflict, "Email already exists",
			map[string]any{"email": "Email already exists"}, nil)
	}

	// Find role merchant
	roleID, err := c.FindRoleMerchant(context, tx)
	if err != nil {
		return nil, err
	}

	userModel := model.User{
		Name:     merchant.Name,
		Email:    merchant.Email,
		Password: merchant.Password,
		Roles: []*model.Role{
			{ID: roleID},
		},
		Merchants: []model.Merchant{
			{
				Name:    merchant.Merchant.Name,
				Address: merchant.Merchant.Address,
			},
		},
	}

	result := tx.Create(&userModel)
	if result.Error != nil {
		return nil, customerror.New(customerror.CodeInternal, "Internal Server Error",
			nil, result.Error,
		)
	}

	userMerchant := &entity.UserMerchant{
		ID:    int(userModel.ID),
		Name:  userModel.Name,
		Email: userModel.Email,
		Merchants: []entity.Merchant{
			{
				Name:    merchant.Merchant.Name,
				Address: merchant.Merchant.Address,
			},
		},
	}
	return userMerchant, nil
}

func (c *merchantRepository) FindRoleMerchant(ctx context.Context, tx *gorm.DB) (uint, error) {

	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	var roleModel model.Role

	result := tx.First(&roleModel, "name = ?", model.RoleMerchant)

	if result.Error != nil {
		isNotfound := errors.Is(result.Error, gorm.ErrRecordNotFound)
		return 0, utils.IfElse(isNotfound,
			customerror.New(
				customerror.CodeNotFound,
				"Role not found",
				map[string]any{"role": "Role not found"},
				result.Error,
			),
			customerror.New(
				customerror.CodeInternal,
				"Internal server error",
				nil,
				result.Error,
			),
		)
	}

	return roleModel.ID, nil
}
