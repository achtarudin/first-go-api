package usecase

import (
	"context"
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/domain/courier/repository"
	"cutbray/first_api/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

type hashPasswordFunc func(password string) (string, error)
type verifyPasswordFunc func(password string, hash string) bool

type CourierUsecase interface {
	Login(ctx context.Context, courier *entity.Courier, verifyPassword verifyPasswordFunc) (*entity.Courier, error)
	Register(ctx context.Context, courier *entity.Courier, hashPassword hashPasswordFunc) (*entity.Courier, error)
	GetAllCouriers(ctx context.Context, entity *entity.SearchCourier) (*entity.CourierWithPaginate[entity.Courier], error)
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
func (c *courierUsecase) Login(ctx context.Context, courier *entity.Courier, verifyPassword verifyPasswordFunc) (*entity.Courier, error) {

	var foundCourier *entity.Courier

	err := c.repo.Trx(ctx, func(txCtx *gorm.DB) error {

		var txErr error

		// Find role courier
		roleId, txErr := c.repo.FindRoleCourier(ctx, "courier", txCtx)
		if txErr != nil {
			return txErr
		}
		// Find courier by email
		courier.RoleId = int(roleId)
		foundCourier, txErr = c.repo.FindByEmail(ctx, courier, txCtx)
		if txErr != nil {
			return txErr
		}

		// Verify password
		if verifyPassword(courier.Password, foundCourier.Password) == false {
			txErr = errors.New("invalid password")
			return txErr
		}

		// Generate token
		token, txErr := utils.GenerateTokenFromIdAndEmail(foundCourier.ID, foundCourier.Email)
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

func (c *courierUsecase) GetAllCouriers(ctx context.Context, searchParams *entity.SearchCourier) (*entity.CourierWithPaginate[entity.Courier], error) {

	result, err := c.repo.ReadAll(ctx, searchParams, nil)

	if err != nil {
		return nil, err
	}
	return result, nil
}
