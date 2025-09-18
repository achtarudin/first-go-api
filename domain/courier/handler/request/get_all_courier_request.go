package request

import (
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/pkg/utils"
)

type GetAllCourierRequest struct {
	Name      *string `form:"name" binding:"omitempty,min=1"`
	Email     *string `form:"email" binding:"omitempty,min=1"`
	Longitude *string `form:"longitude" binding:"omitempty,min=1,longitude"`
	Latitude  *string `form:"latitude" binding:"omitempty,min=1,latitude"`
	Radius    *string `form:"radius" binding:"omitempty"`
	PerPage   *string `form:"per_page" binding:"omitempty,numeric"`
	Page      *string `form:"page" binding:"omitempty,numeric"`
	SortBy    *string `form:"sort_by" binding:"omitempty,oneof=id distance_in_meters"`
	OrderBy   *string `form:"order_by" binding:"omitempty,oneof=asc desc"`
}

func (r *GetAllCourierRequest) ToEntity() *entity.SearchCourier {
	return &entity.SearchCourier{
		Name:      utils.DerefOrDefault(r.Name, ""),
		Email:     utils.DerefOrDefault(r.Email, ""),
		Longitude: utils.ParseFloat64(r.Longitude),
		Latitude:  utils.ParseFloat64(r.Latitude),
		Radius:    utils.DerefOrDefault(utils.ParseIntPointer(r.Radius), 0),
		Page:      utils.DerefOrDefault(utils.ParseIntPointer(r.Page), 1),
		PerPage:   utils.DerefOrDefault(utils.ParseIntPointer(r.PerPage), 10),
		SortBy:    utils.DerefOrDefault(r.SortBy, "id"),
		OrderBy:   utils.DerefOrDefault(r.OrderBy, "asc"),
	}
}
