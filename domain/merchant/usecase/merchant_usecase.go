package usecase

import (
	"context"
	"cutbray/first_api/domain/merchant/entity"
	"cutbray/first_api/domain/merchant/repository"
	"cutbray/first_api/pkg/customerror"
	"errors"

	"gorm.io/gorm"
)

type hashPasswordFunc func(password string) (string, error)
type verifyPasswordFunc func(password string, hash string) bool
type generateTokenFunc func(id int, email string) (string, error)

type MerchantUsecase interface {
	Login(ctx context.Context, merchant *entity.UserMerchantLogin, verifyPassword verifyPasswordFunc, generateToken generateTokenFunc) (*entity.UserMerchant, error)
	Register(ctx context.Context, courier *entity.UserMerchantRegister, hashPassword hashPasswordFunc) (*entity.UserMerchant, error)
}

type merchantUsecase struct {
	repo repository.MerchantRepository
}

func NewMerchantUsecase(repo repository.MerchantRepository) MerchantUsecase {
	return &merchantUsecase{
		repo: repo,
	}
}

func (m *merchantUsecase) Login(ctx context.Context, merchant *entity.UserMerchantLogin,
	verifyPassword verifyPasswordFunc, generateToken generateTokenFunc) (*entity.UserMerchant, error) {

	var foundUserMerchant *entity.UserMerchant

	err := m.repo.Trx(ctx, func(tx *gorm.DB) error {
		var txErr error

		// Find user merchant by email
		foundUserMerchant, txErr = m.repo.FindByEmail(ctx, merchant.Email, tx)
		if txErr != nil {
			return txErr
		}

		// Verify password
		if verifyPassword(merchant.Password, foundUserMerchant.Password) == false {
			txErr = customerror.New(customerror.CodeInvalidInput,
				"Invalid input",
				map[string]any{"email": "invalid input"},
				errors.New("Invalid input"))
			return txErr
		}

		// Generate token
		token, err := generateToken(foundUserMerchant.ID, foundUserMerchant.Email)
		if err != nil {
			txErr = customerror.New(customerror.CodeInternal, "Failed Generate token", nil, err)
			return txErr
		}

		foundUserMerchant.Password = "" // Clear password before returning
		foundUserMerchant.Token = token
		return nil

	})
	return foundUserMerchant, err
}

func (m *merchantUsecase) Register(ctx context.Context, merchant *entity.UserMerchantRegister,
	hashPassword hashPasswordFunc) (*entity.UserMerchant, error) {

	var userMerchant *entity.UserMerchant

	result, err := hashPassword(merchant.Password)

	if err != nil {
		return nil, customerror.New(customerror.CodeHashFailed, "Failed to has password", nil, err)
	}

	merchant.Password = result

	err = m.repo.Trx(ctx, func(tx *gorm.DB) error {
		var txErr error

		userMerchant, txErr = m.repo.Create(ctx, merchant, tx)

		return txErr
	})

	return userMerchant, err
}
