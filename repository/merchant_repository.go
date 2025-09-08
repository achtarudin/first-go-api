package repository

import (
	"cutbray/first_api/repository/model"

	"gorm.io/gorm"
)

// MerchantRepository interface defines methods for merchant operations
type MerchantRepository interface {
	Create(merchant *model.Merchant) error
	GetByID(id uint) (*model.Merchant, error)
	GetByUserID(userID uint) ([]model.Merchant, error)
	Update(merchant *model.Merchant) error
	Delete(id uint) error
	GetAll(limit, offset int) ([]model.Merchant, error)
}

// merchantRepository implements MerchantRepository interface
type merchantRepository struct {
	db *gorm.DB
}

// NewMerchantRepository creates a new merchant repository
func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

func (r *merchantRepository) Create(merchant *model.Merchant) error {
	return r.db.Create(merchant).Error
}

func (r *merchantRepository) GetByID(id uint) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.Preload("User").First(&merchant, id).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *merchantRepository) GetByUserID(userID uint) ([]model.Merchant, error) {
	var merchants []model.Merchant
	err := r.db.Where("user_id = ?", userID).Find(&merchants).Error
	return merchants, err
}

func (r *merchantRepository) Update(merchant *model.Merchant) error {
	return r.db.Save(merchant).Error
}

func (r *merchantRepository) Delete(id uint) error {
	return r.db.Delete(&model.Merchant{}, id).Error
}

func (r *merchantRepository) GetAll(limit, offset int) ([]model.Merchant, error) {
	var merchants []model.Merchant
	err := r.db.Preload("User").Limit(limit).Offset(offset).Find(&merchants).Error
	return merchants, err
}
