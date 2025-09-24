package request

import (
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/pkg/utils"
)

type GetNearestRequest struct {
	Longitude *string `form:"longitude" validate:"required,required_with=Latitude,omitempty,longitude"`
	Latitude  *string `form:"latitude" validate:"required,required_with=Longitude,omitempty,latitude"`
	Radius    *string `form:"radius" validate:"required,numeric,min=0,excludes=-"`
}

func (r *GetNearestRequest) ToEntity() *entity.SearchNearestCourier {
	return &entity.SearchNearestCourier{
		Longitude: utils.ParseFloat64(r.Longitude),
		Latitude:  utils.ParseFloat64(r.Latitude),
		Radius:    utils.DerefOrDefault(utils.ParseFloat64Pointer(r.Radius), 0.0),
	}
}
