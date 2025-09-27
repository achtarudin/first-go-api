package usecase

import (
	"context"
	"cutbray/first_api/domain/auth/entity"
	"cutbray/first_api/domain/auth/repository"
	"cutbray/first_api/pkg/customerror"
	"errors"
)

type hashPasswordFunc func(password string) (string, error)
type verifyPasswordFunc func(password string, hash string) bool
type generateTokenFunc func(user *entity.User) (string, error)

type AuthUsecase interface {
	Login(ctx context.Context, user *entity.User, verifyPassword verifyPasswordFunc, generateToken generateTokenFunc) error
	Register(ctx context.Context, user *entity.User, hashPassword hashPasswordFunc) error
}

type authUsecase struct {
	repo repository.AuthRepository
}

func NewAuthUsecase(repo repository.AuthRepository) AuthUsecase {
	return &authUsecase{
		repo: repo,
	}
}

// Login implements AuthUsecase.
func (a *authUsecase) Login(ctx context.Context, user *entity.User, verifyPassword verifyPasswordFunc, generateToken generateTokenFunc) error {
	inputPassword := user.Password

	foundUser, err := a.repo.FindByEmail(ctx, user.Email)

	if err != nil {
		return err
	}

	if verifyPassword(inputPassword, foundUser.Password) == false {
		return customerror.New(customerror.CodeUnauthorized, "invalid credentials",
			map[string]any{"email": "invalid email or password"}, errors.New("password does not match"),
		)
	}

	token, err := generateToken(user)

	if err != nil {
		return customerror.New(customerror.CodeHashFailed, "failed to generate token",
			nil, err,
		)
	}

	user.ID = foundUser.ID
	user.Token = token
	user.Password = ""

	return nil
}

// Register implements AuthUsecase.
func (a *authUsecase) Register(ctx context.Context, user *entity.User, hashPassword hashPasswordFunc) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return customerror.New(customerror.CodeHashFailed, "failed to hash password",
			nil, err,
		)
	}
	return a.repo.Save(ctx, user, hash)
}
