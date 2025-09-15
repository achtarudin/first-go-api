package usecase

import "cutbray/first_api/domain/courier/repository"

type CourierUsecase interface {
}

type courierUsecase struct {
	repo repository.CourierRepository
}

func NewCourierUsecase(repo repository.CourierRepository) CourierUsecase {
	return &courierUsecase{
		repo: repo,
	}
}
