package usecase

import "cutbray/first_api/domain/auth/repository"

type AuthUsecase struct {
	repo repository.AuthRepository
}

func NewAuthUsecase(repo repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}
