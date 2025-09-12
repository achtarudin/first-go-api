package repository

import (
	"context"
	"cutbray/first_api/domain/auth/entity"
	"cutbray/first_api/pkg/model"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, user *entity.User) (*entity.User, error)
	FindById(ctx context.Context, id string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User, hashPassword string) error
	ReadAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

// Delete implements AuthRepository.
func (a *authRepository) Delete(ctx context.Context, id string) error {
	return nil
}

// FindByEmail implements AuthRepository.
func (a *authRepository) FindByEmail(ctx context.Context, user *entity.User) (*entity.User, error) {
	var userModel model.User
	if err := a.db.WithContext(ctx).First(&userModel, "email = ?", user.Email).Error; err != nil {
		return nil, err
	}

	user.ID = int(userModel.ID)
	user.Password = userModel.Password
	return user, nil
}

// FindById implements AuthRepository.
func (a *authRepository) FindById(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := a.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ReadAll implements AuthRepository.
func (a *authRepository) ReadAll(ctx context.Context) ([]entity.User, error) {
	return nil, nil
}

// Save implements AuthRepository.
func (a *authRepository) Save(ctx context.Context, user *entity.User, hashPassword string) error {
	var userModel model.User
	userModel.Name = user.Name
	userModel.Email = user.Email
	userModel.Password = hashPassword

	// Create a single record
	result := a.db.WithContext(ctx).Create(&userModel)

	if result.Error != nil {
		// Cek error duplikat dari MySQL
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("email already exists")
		}
		return result.Error
	}

	user.ID = int(userModel.ID)
	user.Password = ""
	return nil
}

// Update implements AuthRepository.
func (a *authRepository) Update(ctx context.Context, user *entity.User) error {
	return nil
}
