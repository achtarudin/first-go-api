package usecase

import (
	"context"
	"cutbray/first_api/domain/auth/entity"
	"cutbray/first_api/domain/auth/repository"
	"cutbray/first_api/utils"
	"errors"
)

type hashPasswordFunc func(password string) (string, error)
type verifyPasswordFunc func(password string, hash string) bool
type AuthUsecase interface {
	Login(ctx context.Context, user *entity.User, verifyPassword verifyPasswordFunc) error
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
func (a *authUsecase) Login(ctx context.Context, user *entity.User, verifyPassword verifyPasswordFunc) error {
	inputPassword := user.Password

	foundUser, err := a.repo.FindByEmail(ctx, user)
	if err != nil {
		return errors.New("user not found")
	}

	if verifyPassword(inputPassword, foundUser.Password) == false {
		return errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user)

	if err != nil {
		return errors.New("failed to generate token")
	}

	user.Token = token
	user.Password = ""

	return nil
}

// Register implements AuthUsecase.
func (a *authUsecase) Register(ctx context.Context, user *entity.User, hashPassword hashPasswordFunc) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	return a.repo.Save(ctx, user, hash)
}
