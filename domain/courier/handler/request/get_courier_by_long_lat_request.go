package request

import (
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/pkg/utils"
)

type GetCourierByLongLatRequest struct {
	Longitude *string `form:"longitude" validate:"required,longitude"`
	Latitude  *string `form:"latitude" validate:"required,latitude"`
	PerPage   *string `form:"per_page" validate:"omitempty,number,min=1"`
	Page      *string `form:"page" validate:"omitempty,number,min=1"`
	SortBy    *string `form:"sort_by" validate:"omitempty,oneof=id distance_in_meters"`
	OrderBy   *string `form:"order_by" validate:"omitempty,oneof=ASC DESC asc desc"`
}

func (r *GetCourierByLongLatRequest) ToEntity() *entity.SearchByLongLatCourier {
	return &entity.SearchByLongLatCourier{
		Longitude: utils.ParseFloat64(r.Longitude),
		Latitude:  utils.ParseFloat64(r.Latitude),
		Page:      utils.DerefOrDefault(utils.ParseIntPointer(r.Page), 1),
		PerPage:   utils.DerefOrDefault(utils.ParseIntPointer(r.PerPage), 10),
		SortBy:    utils.DerefOrDefault(r.SortBy, "id"),
		OrderBy:   utils.DerefOrDefault(r.OrderBy, "asc"),
	}
}
