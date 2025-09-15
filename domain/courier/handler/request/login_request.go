package request

import "cutbray/first_api/domain/courier/entity"

// LoginRequest represents the expected structure of a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" validate:"required,min=6,max=100" example:"password"`
}

func (r LoginRequest) ToCourierLogin() entity.Courier {
	return entity.Courier{
		Email:    r.Email,
		Password: r.Password,
	}
}
