package usecase

import (
	"context"
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/domain/courier/repository"
	"errors"

	"gorm.io/gorm"
)

type hashPasswordFunc func(password string) (string, error)
type verifyPasswordFunc func(password string, hash string) bool
type generateTokenFunc func(id int, email string) (string, error)

type CourierUsecase interface {
	Login(ctx context.Context, email string, password string, verifyPassword verifyPasswordFunc, generateToken generateTokenFunc) (*entity.Courier, error)
	Register(ctx context.Context, courier *entity.Courier, hashPassword hashPasswordFunc) (*entity.Courier, error)
	GetAllCouriers(ctx context.Context, search *entity.SearchCourier) (*entity.CourierWithPaginate[entity.Courier], error)
	GetCourierByLongLat(ctx context.Context, search *entity.SearchByLongLatCourier) (*entity.CourierWithPaginate[entity.Courier], error)
	GetNearestCourier(ctx context.Context, search *entity.SearchNearestCourier) (*[]entity.Courier, error)
	UpdateCourier(ctx context.Context, userId int, courier *entity.UpdateCourier, hashPassword hashPasswordFunc) (*entity.Courier, error)
	DeletedCourier(ctx context.Context, userId int) (*entity.Courier, error)
}

type courierUsecase struct {
	repo repository.CourierRepository
}

func NewCourierUsecase(repo repository.CourierRepository) CourierUsecase {
	return &courierUsecase{
		repo: repo,
	}
}

// Login implements AuthUsecase.
func (c *courierUsecase) Login(ctx context.Context, email string, password string,
	verifyPassword verifyPasswordFunc, generateToken generateTokenFunc) (*entity.Courier, error) {

	var foundCourier *entity.Courier

	err := c.repo.Trx(ctx, func(txCtx *gorm.DB) error {

		var txErr error

		// Find role courier
		roleId, txErr := c.repo.FindRoleCourier(ctx, "courier", txCtx)
		if txErr != nil {
			return txErr
		}
		// Find courier by email
		foundCourier, txErr = c.repo.FindByEmail(ctx, email, int(roleId), txCtx)
		if txErr != nil {
			return txErr
		}

		// Verify password
		if verifyPassword(password, foundCourier.Password) == false {
			txErr = errors.New("invalid password")
			return txErr
		}

		// Generate token
		token, txErr := generateToken(foundCourier.ID, foundCourier.Email)
		if txErr != nil {
			txErr = errors.New("failed generate token")
			return txErr
		}

		foundCourier.Password = "" // Clear password before returning
		foundCourier.Token = token
		return nil
	})
	return foundCourier, err
}

// Register implements AuthUsecase.
func (c *courierUsecase) Register(ctx context.Context, courier *entity.Courier, hashPassword hashPasswordFunc) (*entity.Courier, error) {
	// Hash password
	hashedPassword, err := hashPassword(courier.Password)
	if err != nil {
		return nil, err
	}
	courier.Password = hashedPassword

	// Variabel untuk menampung hasil akhir
	var createdCourier *entity.Courier

	// Jalankan transaksi
	err = c.repo.Trx(ctx, func(txCtx *gorm.DB) error {
		var txErr error

		// Find role courier
		roleId, txErr := c.repo.FindRoleCourier(ctx, "courier", txCtx)

		if txErr != nil {
			return txErr
		}

		// Create courier
		courier.RoleId = int(roleId)
		createdCourier, txErr = c.repo.Create(ctx, courier, txCtx)

		if txErr != nil {
			return txErr
		}

		return nil
	})
	return createdCourier, err

}

// GetAllCouriers implements CourierUsecase.
func (c *courierUsecase) GetAllCouriers(ctx context.Context, searchParams *entity.SearchCourier) (*entity.CourierWithPaginate[entity.Courier], error) {

	var couriers *entity.CourierWithPaginate[entity.Courier]

	err := c.repo.Trx(ctx, func(tx *gorm.DB) error {

		var txErr error

		couriers, txErr = c.repo.ReadAll(ctx, searchParams, tx)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	return couriers, err

}

// GetCourierByLongLat implements CourierUsecase.
func (c *courierUsecase) GetCourierByLongLat(ctx context.Context, searchParams *entity.SearchByLongLatCourier) (*entity.CourierWithPaginate[entity.Courier], error) {
	var couriers *entity.CourierWithPaginate[entity.Courier]

	err := c.repo.Trx(ctx, func(tx *gorm.DB) error {

		var txErr error

		couriers, txErr = c.repo.ReadByLongLat(ctx, searchParams, tx)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	return couriers, err
}

func (c *courierUsecase) GetNearestCourier(ctx context.Context, searchParams *entity.SearchNearestCourier) (*[]entity.Courier, error) {
	var couriers *[]entity.Courier

	err := c.repo.Trx(ctx, func(tx *gorm.DB) error {

		var txErr error

		couriers, txErr = c.repo.ReadNearest(ctx, searchParams, tx)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	return couriers, err
}

func (c *courierUsecase) UpdateCourier(ctx context.Context, userId int, courier *entity.UpdateCourier, hashPassword hashPasswordFunc) (*entity.Courier, error) {
	var updatedCourier *entity.Courier

	if courier.Password != "" {

		hashedPassword, err := hashPassword(courier.Password)

		if err != nil {
			return nil, err
		}

		courier.Password = hashedPassword
	}

	err := c.repo.Trx(ctx, func(tx *gorm.DB) error {

		var txErr error

		updatedCourier, txErr = c.repo.Update(ctx, userId, courier, tx)

		if txErr != nil {
			return txErr
		}

		return nil
	})

	return updatedCourier, err
}

func (c *courierUsecase) DeletedCourier(ctx context.Context, userId int) (*entity.Courier, error) {
	var deletedCourier *entity.Courier

	err := c.repo.Trx(ctx, func(tx *gorm.DB) error {

		var txErr error

		deletedCourier, txErr = c.repo.Delete(ctx, userId, tx)

		if txErr != nil {
			return txErr
		}

		return nil
	})

	return deletedCourier, err
}
