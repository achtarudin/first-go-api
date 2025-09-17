package repository

import (
	"context"
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/model"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
		Where("role_id = ?", courier.RoleId)

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

	roleId, err := c.FindRoleCourier(ctx, "courier", nil)
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
		"users.id",
		"users.name",
		"users.email",
		"couriers.phone",
		"couriers.latitude",
		"couriers.longitude",
	}

	// Argumen untuk placeholders '?'
	var selectArgs []interface{}

	// Filter by name
	name, exists := searchParams["name"]
	if exists && name != "" {
		query = query.Where("users.name LIKE ?", "%"+name+"%")
	}

	// Filter by email
	email, exists := searchParams["email"]
	if exists && email != "" {
		query = query.Where("users.email LIKE ?", "%"+email+"%")
	}

	// Filter by latitude and longitude to calculate distance
	lonStr, lonExists := searchParams["longitude"]
	latStr, latExists := searchParams["latitude"]

	if lonExists && latExists && lonStr != "" && latStr != "" {
		lon, errLon := strconv.ParseFloat(lonStr, 64)
		lat, errLat := strconv.ParseFloat(latStr, 64)

		if errLat == nil && errLon == nil {
			selectParts = append(selectParts, "ST_Distance_Sphere(POINT(?, ?), couriers.location) AS distance_in_meters")
			selectArgs = append(selectArgs, lon, lat)
		} else {
			// Jika parsing gagal (input tidak valid), beri nilai NULL
			selectParts = append(selectParts, "NULL AS distance_in_meters")
		}
	} else {
		selectParts = append(selectParts, "NULL AS distance_in_meters")
	}

	// Default values for pagination
	paginate := 10
	if perPageStr, exists := searchParams["perPage"]; exists && perPageStr != "" {
		perPage, err := strconv.Atoi(perPageStr)
		if err == nil {
			paginate = perPage
		}
	}

	// Default values for current page
	currentPage := 1
	if pageStr, exists := searchParams["page"]; exists && pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil {
			currentPage = page
		}
	}

	// Hitung total data
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (currentPage - 1) * paginate

	// Siapkan variabel untuk menampung hasil
	var results []entity.Courier

	// Gabungkan bagian-bagian SELECT menjadi satu string
	selectStatement := strings.Join(selectParts, ", ")

	// Eksekusi query dengan SELECT, LIMIT, OFFSET, dan ORDER BY
	err = query.Select(selectStatement, selectArgs...).
		Limit(paginate).
		Offset(offset).
		Order("distance_in_meters ASC").
		// Having("distance_in_meters <= 50").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	response := &entity.CourierWithPaginate[entity.Courier]{
		CurrentPage: currentPage,
		Data:        results,
		PerPage:     paginate,
		Total:       total,
	}
	return response, nil
}
