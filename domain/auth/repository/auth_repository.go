package repository

import (
	"context"
	"cutbray/first_api/domain/auth/entity"
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/customerror"
	"cutbray/first_api/pkg/model"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Trx(ctx context.Context, fn func(tx *gorm.DB) error) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindById(ctx context.Context, id string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User, hashedPassword string) error
	SaveWithTrx(ctx context.Context, user *entity.User, hashedPassword string, tx *gorm.DB) error

	ReadAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}

type authRepository struct {
	db *infra.Database
}

func NewAuthRepository(db *infra.Database) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (a *authRepository) Trx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return a.db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// FindByEmail implements AuthRepository.
func (a *authRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.User

	result := a.db.DB.WithContext(ctx).Select("id", "name", "email", "password").First(&userModel, "email = ?", email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customerror.New(customerror.CodeNotFound, "user not found",
				map[string]any{"email": "user with this email not found"}, result.Error,
			)
		}
		return nil, customerror.New(customerror.CodeDBQuery, "failed to query user by email",
			nil, result.Error,
		)
	}

	user := &entity.User{
		ID:       int(userModel.ID),
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: userModel.Password,
	}
	return user, nil
}

// FindById implements AuthRepository.
func (a *authRepository) FindById(ctx context.Context, id string) (*entity.User, error) {
	user := model.User{}

	result := a.db.DB.WithContext(ctx).First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customerror.New(customerror.CodeNotFound, "user not found",
				map[string]any{"id": "user with this id not found"}, result.Error,
			)
		}
		return nil, customerror.New(customerror.CodeDBQuery, "failed to query user by email",
			nil, result.Error,
		)
	}

	return &entity.User{
		ID:       int(user.ID),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

// Save implements AuthRepository.
func (a *authRepository) Save(ctx context.Context, user *entity.User, hashedPassword string) error {

	userIsExists, err := a.FindByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if userIsExists != nil {
		return customerror.New(customerror.CodeConflict, "email already exists",
			map[string]any{"email": "email already exists"}, gorm.ErrDuplicatedKey,
		)
	}

	m := &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	// Create a single record
	result := a.db.DB.WithContext(ctx).Create(m)

	if result.Error != nil {
		return customerror.New(customerror.CodeInternal, "failed to save the user",
			nil, result.Error,
		)
	}

	user.ID = int(m.ID)
	user.Password = ""
	return nil
}

// Save implements AuthRepository.
func (a *authRepository) SaveWithTrx(ctx context.Context, user *entity.User, hashedPassword string, tx *gorm.DB) error {

	if tx == nil {
		tx = a.db.DB.WithContext(ctx)
	}

	userIsExists, err := a.FindByEmail(ctx, user.Email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if userIsExists != nil {
		return customerror.New(customerror.CodeConflict, "email already exists",
			map[string]any{"email": "email already exists"}, gorm.ErrDuplicatedKey,
		)
	}

	m := &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	// Create a single record
	result := tx.Create(m)

	if result.Error != nil {
		return customerror.New(customerror.CodeInternal, "failed to save the user",
			nil, result.Error,
		)
	}

	user.ID = int(m.ID)
	user.Password = ""
	return nil
}

// ReadAll implements AuthRepository.
func (a *authRepository) ReadAll(ctx context.Context) ([]entity.User, error) {

	var users []model.User

	result := a.db.DB.WithContext(ctx).Select("id", "name", "email").Find(&users)
	if result.Error != nil {
		return nil, customerror.New(customerror.CodeDBQuery, "failed to query all users",
			nil, result.Error,
		)
	}

	var entityUsers = make([]entity.User, 0, len(users))
	for _, u := range users {
		entityUsers = append(entityUsers, entity.User{
			ID:       int(u.ID),
			Name:     u.Name,
			Email:    u.Email,
			Password: u.Password,
		})
	}

	return entityUsers, nil
}

// Update implements AuthRepository.
func (a *authRepository) Update(ctx context.Context, user *entity.User) error {
	return nil
}

// Delete implements AuthRepository.
func (a *authRepository) Delete(ctx context.Context, id string) error {
	return nil
}
