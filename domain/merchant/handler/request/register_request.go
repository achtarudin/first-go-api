package request

import "cutbray/first_api/domain/merchant/entity"

// LoginRequest represents the expected structure of a login request
type RegisterRequest struct {
	Name                  string `json:"name" binding:"required" validate:"required" example:"John Doe"`
	Email                 string `json:"email" binding:"required" validate:"required,email" example:"user@example.com"`
	Password              string `json:"password" binding:"required" validate:"required,min=6,max=100" example:"password"`
	Password_Confirmation string `json:"password_confirmation" binding:"required" validate:"required,min=6,max=100,eqfield=Password" example:"password"`
}

func (r RegisterRequest) ToMerchantRegister() entity.UserMerchant {
	return entity.UserMerchant{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}
