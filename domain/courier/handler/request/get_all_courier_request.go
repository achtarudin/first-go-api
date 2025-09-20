package request

import (
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/pkg/utils"
	"fmt"
)

type GetAllCourierRequest struct {
	Name      *string `form:"name" validate:"omitempty,min=1"`
	Email     *string `form:"email" validate:"omitempty,min=1"`
	Longitude *string `form:"longitude" validate:"required_with=Latitude,omitempty,longitude"`
	Latitude  *string `form:"latitude" validate:"required_with=Longitude,omitempty,latitude"`
	Radius    *string `form:"radius" validate:"omitempty,numeric,min=0,excludes=-"`
	PerPage   *string `form:"per_page" validate:"omitempty,number,min=1"`
	Page      *string `form:"page" validate:"omitempty,number,min=1"`
	SortBy    *string `form:"sort_by" validate:"omitempty,oneof=id distance_in_meters"`
	OrderBy   *string `form:"order_by" validate:"omitempty,oneof=ASC DESC asc desc"`
}

func (r *GetAllCourierRequest) ToEntity() *entity.SearchCourier {
	fmt.Println(utils.DerefOrDefault(utils.ParseFloat64Pointer(r.Radius), 0.0))
	return &entity.SearchCourier{
		Name:      utils.DerefOrDefault(r.Name, ""),
		Email:     utils.DerefOrDefault(r.Email, ""),
		Longitude: utils.ParseFloat64(r.Longitude),
		Latitude:  utils.ParseFloat64(r.Latitude),
		Radius:    utils.DerefOrDefault(utils.ParseFloat64Pointer(r.Radius), 0.0),
		Page:      utils.DerefOrDefault(utils.ParseIntPointer(r.Page), 1),
		PerPage:   utils.DerefOrDefault(utils.ParseIntPointer(r.PerPage), 10),
		SortBy:    utils.DerefOrDefault(r.SortBy, "id"),
		OrderBy:   utils.DerefOrDefault(r.OrderBy, "asc"),
	}
}
