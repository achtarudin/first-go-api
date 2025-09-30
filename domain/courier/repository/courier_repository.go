package repository

import (
	"context"
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/customerror"
	"cutbray/first_api/pkg/model"
	"cutbray/first_api/pkg/utils"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CourierRepository interface {
	Trx(ctx context.Context, clouseure func(tx *gorm.DB) error) error
	Create(ctx context.Context, courier *entity.Courier, tx *gorm.DB) (*entity.Courier, error)
	FindRoleCourier(ctx context.Context, roleName model.RoleStatus, tx *gorm.DB) (uint, error)
	FindByEmail(ctx context.Context, email string, roleId int, tx *gorm.DB) (*entity.Courier, error)
	ReadAll(ctx context.Context, searchParams *entity.SearchCourier, tx *gorm.DB) (*entity.CourierWithPaginate[entity.Courier], error)
	ReadByLongLat(ctx context.Context, searchParams *entity.SearchByLongLatCourier, tx *gorm.DB) (*entity.CourierWithPaginate[entity.Courier], error)
	ReadNearest(ctx context.Context, searchParams *entity.SearchNearestCourier, tx *gorm.DB) (*[]entity.Courier, error)
	Update(ctx context.Context, userId int, courier *entity.UpdateCourier, tx *gorm.DB) (*entity.Courier, error)
	Delete(ctx context.Context, userId int, tx *gorm.DB) (*entity.Courier, error)
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

	defaultLocation := model.Point{
		Lng: courier.Longitude,
		Lat: courier.Latitude,
	}

	userModel := model.User{
		Name:     courier.Name,
		Email:    courier.Email,
		Password: courier.Password,
		Roles:    []*model.Role{{ID: uint(courier.RoleId)}},
		Courier: model.Courier{Phone: courier.Phone,
			Longitude: &defaultLocation.Lng,
			Latitude:  &defaultLocation.Lat,
			Location:  defaultLocation,
		},
	}

	result := tx.Create(&userModel)

	if result.Error != nil {

		isDuplicateKey := errors.Is(result.Error, gorm.ErrDuplicatedKey)

		return nil, utils.IfElse(isDuplicateKey,
			customerror.New(
				customerror.CodeConflict,
				"email already exists",
				map[string]any{"email": "email already exists"},
				result.Error,
			),
			customerror.New(
				customerror.CodeInternal,
				"failed to save the user",
				nil,
				result.Error,
			),
		)
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

func (c *courierRepository) FindByEmail(ctx context.Context, email string, roleId int, tx *gorm.DB) (*entity.Courier, error) {

	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	var userModel model.User

	// Subquery untuk mendapatkan user_id dari user_roles dengan role courier
	subQuery := tx.
		Table("user_roles").
		Select("user_roles.user_id").
		Where("role_id = ?", roleId)

	// Cari user dengan email dan role courier
	result := tx.Where("email = ?", email).
		Where("id IN (?)", subQuery).
		Where("id IN (?)", tx.Model(&model.Courier{}).Select("user_id")).
		Preload("Courier").Preload("Roles").First(&userModel)

	if result.Error != nil {
		isNotfound := errors.Is(result.Error, gorm.ErrRecordNotFound)

		return nil, utils.IfElse(isNotfound,
			customerror.New(
				customerror.CodeNotFound,
				"User not found",
				map[string]any{"email": "user not found"},
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

	foundCourier := &entity.Courier{
		ID:       int(userModel.ID),
		Name:     userModel.Name,
		Email:    userModel.Email,
		Phone:    userModel.Courier.Phone,
		Password: userModel.Password,
	}

	return foundCourier, nil
}

func (c *courierRepository) ReadAll(ctx context.Context, searchParams *entity.SearchCourier, tx *gorm.DB) (*entity.CourierWithPaginate[entity.Courier], error) {
	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	roleId, err := c.FindRoleCourier(ctx, model.RoleCourier, nil)
	if err != nil {
		return nil, err
	}

	// Bangun query dasar yang efisien dengan JOIN
	query := tx.Model(&model.User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN couriers ON couriers.user_id = users.id").
		Where("user_roles.role_id = ?", roleId).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL")

	// Kolom-kolom dasar yang selalu dipilih
	selectParts := []string{
		"users.name",
		"users.email",
		"couriers.user_id as id",
		"couriers.phone",
		"couriers.latitude",
		"couriers.longitude",
	}

	// Argumen untuk placeholders '?'
	var selectArgs []interface{}

	// Filter by name
	if searchParams.Name != "" {
		query = query.Where("users.name LIKE ?", "%"+searchParams.Name+"%")
	}

	// Filter by email
	if searchParams.Email != "" {
		query = query.Where("users.email LIKE ?", "%"+searchParams.Email+"%")
	}

	// Filter by latitude and longitude to calculate distance
	hasLatLong := searchParams.Latitude != 0.0 && searchParams.Longitude != 0.0
	if hasLatLong {
		selectParts = append(selectParts, "ST_Distance_Sphere(POINT(?, ?), couriers.location) AS distance_in_meters")
		selectArgs = append(selectArgs, searchParams.Longitude, searchParams.Latitude)
	} else {
		selectParts = append(selectParts, "NULL AS distance_in_meters")
	}

	// Variabel untuk menampung hasil total
	var total int64

	// Variabel untuk menampung hasil results
	var results []entity.Courier

	selectStatement := strings.Join(selectParts, ", ")
	outerQuery := tx.Table("(?) AS couriers_with_distance", query.Select(selectStatement, selectArgs...))

	if hasLatLong && searchParams.Radius > 0 {
		outerQuery = outerQuery.Where("distance_in_meters <= ?", searchParams.Radius)
	}

	if err := outerQuery.Count(&total).Error; err != nil {
		return nil, customerror.New(
			customerror.CodeInternal,
			"Internal server error",
			nil,
			err,
		)
	}

	if total > 0 {
		// Calculate offset
		offset := (searchParams.Page - 1) * searchParams.PerPage

		// Eksekusi query dengan SELECT, LIMIT, OFFSET, dan ORDER BY
		err = outerQuery.
			Order(fmt.Sprintf("%s %s", searchParams.SortBy, searchParams.OrderBy)).
			Limit(searchParams.PerPage).
			Offset(offset).
			Scan(&results).Error

		if err != nil {
			return nil, customerror.New(
				customerror.CodeInternal,
				"Internal server error",
				nil,
				err,
			)
		}
	}

	response := &entity.CourierWithPaginate[entity.Courier]{
		CurrentPage: searchParams.Page,
		Data:        results,
		PerPage:     searchParams.PerPage,
		Total:       total,
	}
	return response, nil
}

func (c *courierRepository) ReadByLongLat(ctx context.Context, searchParams *entity.SearchByLongLatCourier, tx *gorm.DB) (*entity.CourierWithPaginate[entity.Courier], error) {
	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	roleId, err := c.FindRoleCourier(ctx, model.RoleCourier, nil)
	if err != nil {
		return nil, err
	}

	query := tx.Model(model.Courier{}).
		Joins("JOIN users ON users.id = couriers.user_id").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleId).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL")

	// Kolom-kolom dasar yang selalu dipilih
	selectParts := []string{
		"users.name",
		"users.email",
		"couriers.user_id as id",
		"couriers.phone",
		"couriers.latitude",
		"couriers.longitude",
		"ST_Distance_Sphere(POINT(?, ?), couriers.location) AS distance_in_meters",
	}

	var selectArgs []interface{}

	selectStatement := strings.Join(selectParts, ", ")

	selectArgs = append(selectArgs, searchParams.Longitude, searchParams.Latitude)

	outerQuery := tx.Table("(?) AS couriers_with_distance", query.Select(selectStatement, selectArgs...))

	var total int64
	var results []entity.Courier

	if err := outerQuery.Count(&total).Error; err != nil {
		return nil, customerror.New(
			customerror.CodeInternal,
			"Internal server error",
			nil,
			err,
		)
	}

	if total > 0 {
		offset := (searchParams.Page - 1) * searchParams.PerPage

		result := outerQuery.
			Order(fmt.Sprintf("%s %s", searchParams.SortBy, searchParams.OrderBy)).
			Limit(searchParams.PerPage).
			Offset(offset).
			Scan(&results)

		if result.Error != nil {
			return nil, customerror.New(
				customerror.CodeInternal,
				"Internal server error",
				nil,
				result.Error,
			)
		}
	}

	response := &entity.CourierWithPaginate[entity.Courier]{
		CurrentPage: searchParams.Page,
		PerPage:     searchParams.PerPage,
		Data:        results,
		Total:       total,
	}
	return response, nil
}

func (c *courierRepository) ReadNearest(ctx context.Context, searchParams *entity.SearchNearestCourier, tx *gorm.DB) (*[]entity.Courier, error) {
	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	// Find role ID for courier
	roleId, err := c.FindRoleCourier(ctx, model.RoleCourier, nil)
	if err != nil {
		return nil, err
	}

	query := tx.Model(model.Courier{}).
		Joins("JOIN users ON users.id = couriers.user_id").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleId).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL")

	// Kolom-kolom dasar yang selalu dipilih
	selectArgs := []interface{}{}
	selectParts := []string{
		"users.name",
		"users.email",
		"couriers.user_id as id",
		"couriers.phone",
		"couriers.latitude",
		"couriers.longitude",
		"ST_Distance_Sphere(POINT(?, ?), couriers.location) AS distance_in_meters",
	}

	selectStatement := strings.Join(selectParts, ", ")

	selectArgs = append(selectArgs, searchParams.Longitude, searchParams.Latitude)

	outerQuery := tx.Table("(?) AS couriers_with_distance", query.Select(selectStatement, searchParams.Longitude, searchParams.Latitude))

	var results []entity.Courier

	result := outerQuery.Where("distance_in_meters <= ?", searchParams.Radius).
		Order(fmt.Sprintf("%s %s", "distance_in_meters", "ASC")).
		Scan(&results)

	if result.Error != nil {
		return nil, customerror.New(
			customerror.CodeInternal,
			"Internal server error",
			nil,
			result.Error,
		)
	}
	return &results, nil
}

func (c *courierRepository) Update(ctx context.Context, userId int, courier *entity.UpdateCourier, tx *gorm.DB) (*entity.Courier, error) {
	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	// Find existing user
	var existingUser model.User
	err := tx.Preload("Courier").First(&existingUser, uint(userId)).Error

	// Handle error if is not existing user
	if err != nil {
		isNotfound := errors.Is(err, gorm.ErrRecordNotFound)
		return nil, utils.IfElse(isNotfound,
			customerror.New(
				customerror.CodeNotFound,
				"Courier with the given user id not found",
				nil,
				err,
			),
			customerror.New(
				customerror.CodeInternal,
				"Internal server error",
				nil,
				err,
			),
		)
	}

	// Update user
	userUpdate := tx.Model(&existingUser).Updates(model.User{
		Name:     courier.Name,
		Password: courier.Password,
	})

	// handle error if update user failed
	if userUpdate.Error != nil {
		return nil, customerror.New(
			customerror.CodeInternal,
			"Internal server error",
			nil,
			userUpdate.Error,
		)
	}

	// Update courier
	if courier.Longitude == 0 {
		courier.Longitude = *existingUser.Courier.Longitude
	}

	if courier.Latitude == 0 {
		courier.Latitude = *existingUser.Courier.Latitude
	}

	courierUpdate := tx.Model(&existingUser.Courier).Updates(model.Courier{
		Phone:     courier.Phone,
		Longitude: &courier.Longitude,
		Latitude:  &courier.Latitude,
		Location: model.Point{
			Lng: courier.Longitude,
			Lat: courier.Latitude,
		},
	})

	// Handle error if update courier failed
	if courierUpdate.Error != nil {
		return nil, customerror.New(
			customerror.CodeInternal,
			"Internal server error",
			nil,
			courierUpdate.Error,
		)
	}

	// Reload the updated user and courier
	err = tx.Preload("Courier").First(&existingUser, uint(userId)).Error
	if err != nil {
		return nil, customerror.New(
			customerror.CodeInternal,
			"Internal server error",
			nil,
			err,
		)
	}

	updatedCourier := &entity.Courier{
		ID:        int(existingUser.ID),
		Name:      existingUser.Name,
		Email:     existingUser.Email,
		Phone:     existingUser.Courier.Phone,
		Longitude: *existingUser.Courier.Longitude,
		Latitude:  *existingUser.Courier.Latitude,
	}

	return updatedCourier, nil

}

func (c *courierRepository) Delete(ctx context.Context, userId int, tx *gorm.DB) (*entity.Courier, error) {
	if tx == nil {
		tx = c.db.DB.WithContext(ctx)
	}

	// Find existing user
	var existingUser model.User
	err := tx.Preload("Courier").First(&existingUser, uint(userId)).Error

	// Handle error if is not existing user

	if err != nil {
		isNotfound := errors.Is(err, gorm.ErrRecordNotFound)
		return nil, utils.IfElse(isNotfound,
			customerror.New(
				customerror.CodeNotFound,
				"User not found",
				nil,
				err,
			),
			customerror.New(
				customerror.CodeInternal,
				"Internal server error",
				nil,
				err,
			),
		)
	}

	// Delete courier record first if present
	if existingUser.Courier.ID != 0 {
		result := tx.Delete(&existingUser.Courier)
		if result.Error != nil {
			return nil, customerror.New(
				customerror.CodeInternal,
				"Failed to delete courier record",
				nil,
				result.Error,
			)
		}
	}

	// Delete user record
	result := tx.Delete(&existingUser)
	if result.Error != nil {
		return nil, customerror.New(
			customerror.CodeInternal,
			"Failed to delete user record",
			nil,
			result.Error,
		)
	}

	deletedCourier := &entity.Courier{
		ID:    int(existingUser.ID),
		Name:  existingUser.Name,
		Email: existingUser.Email,
		Phone: existingUser.Courier.Phone,
	}

	return deletedCourier, nil
}
