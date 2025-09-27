package request

import (
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/pkg/utils"
)

type UpdateRequest struct {
	Longitude            *string `json:"longitude" validate:"required_with=Latitude,omitempty,longitude" example:"106.8260"`
	Latitude             *string `json:"latitude" validate:"required_with=Longitude,omitempty,latitude" example:"-6.1790"`
	Phone                *string `json:"phone" validate:"omitempty,min=6,max=15" example:"+6282118302438"`
	Name                 *string `json:"name" validate:"required_without_all=Longitude Latitude Password PasswordConfirmation Phone,omitempty,min=3" example:"John Doe"`
	Password             *string `json:"password" validate:"required_with=PasswordConfirmation,omitempty,min=6,max=100" example:"password"`
	PasswordConfirmation *string `json:"password_confirmation" validate:"required_with=Password,omitempty,min=6,max=100,eqfield=Password" example:"password"`
}

func (r *UpdateRequest) ToEntity() *entity.UpdateCourier {
	return &entity.UpdateCourier{
		Longitude: utils.ParseFloat64(r.Longitude),
		Latitude:  utils.ParseFloat64(r.Latitude),
		Phone:     utils.DerefOrDefault(r.Phone, ""),
		Name:      utils.DerefOrDefault(r.Name, ""),
		Password:  utils.DerefOrDefault(r.Password, ""),
	}
}
