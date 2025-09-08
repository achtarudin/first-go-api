package repository

import (
	"cutbray/first_api/repository/model"

	"gorm.io/gorm"
)

// CourierRepository interface defines methods for courier operations
type CourierRepository interface {
	Create(courier *model.Courier) error
	GetByID(id uint) (*model.Courier, error)
	Update(courier *model.Courier) error
	Delete(id uint) error
	GetAll(limit, offset int) ([]model.Courier, error)
	GetNearby(latitude, longitude float64, radius float64) ([]model.Courier, error)
}

// courierRepository implements CourierRepository interface
type courierRepository struct {
	db *gorm.DB
}

// NewCourierRepository creates a new courier repository
func NewCourierRepository(db *gorm.DB) CourierRepository {
	return &courierRepository{db: db}
}

func (r *courierRepository) Create(courier *model.Courier) error {
	return r.db.Create(courier).Error
}

func (r *courierRepository) GetByID(id uint) (*model.Courier, error) {
	var courier model.Courier
	err := r.db.First(&courier, id).Error
	if err != nil {
		return nil, err
	}
	return &courier, nil
}

func (r *courierRepository) Update(courier *model.Courier) error {
	return r.db.Save(courier).Error
}

func (r *courierRepository) Delete(id uint) error {
	return r.db.Delete(&model.Courier{}, id).Error
}

func (r *courierRepository) GetAll(limit, offset int) ([]model.Courier, error) {
	var couriers []model.Courier
	err := r.db.Limit(limit).Offset(offset).Find(&couriers).Error
	return couriers, err
}

// GetNearby finds couriers within a certain radius (simplified implementation)
func (r *courierRepository) GetNearby(latitude, longitude float64, radius float64) ([]model.Courier, error) {
	var couriers []model.Courier
	// Simple box query - in production, you'd want to use proper geospatial queries
	err := r.db.Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
		latitude-radius, latitude+radius,
		longitude-radius, longitude+radius).Find(&couriers).Error
	return couriers, err
}
