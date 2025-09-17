package repository

import (
	"context"
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/model"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type CourierRepository interface {
	Trx(ctx context.Context, clouseure func(tx *gorm.DB) error) error
	Create(ctx context.Context, courier *entity.Courier, tx *gorm.DB) (*entity.Courier, error)
	FindRoleCourier(ctx context.Context, roleName model.RoleStatus, tx *gorm.DB) (uint, error)
	FindByEmail(ctx context.Context, courier *entity.Courier, tx *gorm.DB) (*entity.Courier, error)
	ReadAll(ctx context.Context, searchParams map[string]string, tx *gorm.DB) (*entity.CourierWithPaginate[entity.Courier], error)
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

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found") // Pesan error bisa disesuaikan
		}
		return nil, result.Error // Kembalikan error database lainnya
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

func (c *courierRepository) ReadAll(ctx context.Context, searchParams map[string]string, tx *gorm.DB) (*entity.CourierWithPaginate[entity.Courier], error) {

	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	roleId, err := c.FindRoleCourier(ctx, "courier", tx)
	if err != nil {
		return nil, err
	}

	// Subquery untuk mendapatkan user_id dari user_roles dengan role courier
	subQuery := tx.
		Table("user_roles").
		Select("user_roles.user_id").
		Where("role_id = ?", roleId)

	// Query utama untuk mendapatkan semua courier dengan filter dan pagination
	query := tx.Where("id IN (?)", tx.Model(&model.Courier{}).Select("user_id")).
		Where("id IN (?)", subQuery)

	// Filter by name
	if name, exists := searchParams["name"]; exists && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Filter by email
	if email, exists := searchParams["email"]; exists && email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	// Pagination
	perPageApp := 10
	if perPageStr, exists := searchParams["perPage"]; exists && perPageStr != "" {
		perPage, err := strconv.Atoi(perPageStr)
		if err == nil {
			perPageApp = perPage
		}
	}

	// Default page is 1
	currentPage := 1
	if pageStr, exists := searchParams["page"]; exists && pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil {
			currentPage = page
		}
	}

	// Get total count before applying limit and offset
	var total int64
	if err := query.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset for pagination
	offset := (currentPage - 1) * perPageApp

	// Execute query with pagination
	var userModel []model.User
	if err := query.Limit(perPageApp).Offset(offset).Preload("Courier").Find(&userModel).Error; err != nil {
		return nil, err
	}

	// Map userModel to entity.Courier
	couriers := make([]entity.Courier, 0, len(userModel))
	for _, user := range userModel {
		courier := entity.Courier{
			ID:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Courier.Phone,
		}
		couriers = append(couriers, courier)
	}

	result := &entity.CourierWithPaginate[entity.Courier]{
		CurrentPage: currentPage,
		Data:        couriers,
		PerPage:     perPageApp,
		Total:       total,
	}
	return result, nil
}
