package repository

import (
	"cutbray/first_api/repository/model"

	"gorm.io/gorm"
)

// FoodRepository interface defines methods for food operations
type FoodRepository interface {
	Create(food *model.Food) error
	GetByID(id uint) (*model.Food, error)
	GetByMerchantID(merchantID uint) ([]model.Food, error)
	Update(food *model.Food) error
	Delete(id uint) error
	GetAll(limit, offset int) ([]model.Food, error)
}

// foodRepository implements FoodRepository interface
type foodRepository struct {
	db *gorm.DB
}

// NewFoodRepository creates a new food repository
func NewFoodRepository(db *gorm.DB) FoodRepository {
	return &foodRepository{db: db}
}

func (r *foodRepository) Create(food *model.Food) error {
	return r.db.Create(food).Error
}

func (r *foodRepository) GetByID(id uint) (*model.Food, error) {
	var food model.Food
	err := r.db.Preload("Merchant").First(&food, id).Error
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) GetByMerchantID(merchantID uint) ([]model.Food, error) {
	var foods []model.Food
	err := r.db.Where("merchant_id = ?", merchantID).Find(&foods).Error
	return foods, err
}

func (r *foodRepository) Update(food *model.Food) error {
	return r.db.Save(food).Error
}

func (r *foodRepository) Delete(id uint) error {
	return r.db.Delete(&model.Food{}, id).Error
}

func (r *foodRepository) GetAll(limit, offset int) ([]model.Food, error) {
	var foods []model.Food
	err := r.db.Preload("Merchant").Limit(limit).Offset(offset).Find(&foods).Error
	return foods, err
}
