package request

import "cutbray/first_api/domain/courier/entity"

// LoginRequest represents the expected structure of a login request
type RegisterRequest struct {
	Name                  string `json:"name" binding:"required" validate:"required" example:"John Doe"`
	Email                 string `json:"email" binding:"required" validate:"required,email" example:"user@example.com"`
	Phone                 string `json:"phone" binding:"required" validate:"required,min=6,max=15" example:"+6282118302438"`
	Password              string `json:"password" binding:"required" validate:"required,min=6,max=100" example:"password"`
	Password_Confirmation string `json:"password_confirmation" binding:"required" validate:"required,min=6,max=100,eqfield=Password" example:"password"`
}

func (r RegisterRequest) ToCourierRegister() entity.Courier {
	return entity.Courier{
		Name:     r.Name,
		Email:    r.Email,
		Phone:    r.Phone,
		Password: r.Password,
	}
}
