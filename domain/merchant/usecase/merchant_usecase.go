package usecase

import "cutbray/first_api/domain/merchant/repository"

type MerchantUsecase interface {
}

type merchantUsecase struct {
	repo repository.MerchantRepository
}

func NewMerchantUsecase(repo repository.MerchantRepository) MerchantUsecase {
	return &merchantUsecase{
		repo: repo,
	}
}
