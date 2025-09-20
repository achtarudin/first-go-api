package entity

type Courier struct {
	ID               int     `json:"id"`
	RoleId           int     `json:"role_id,omitempty"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Password         string  `json:"password,omitempty"`
	Token            string  `json:"token,omitempty"`
	Phone            string  `json:"phone"`
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	DistanceInMeters float64 `json:"distance_in_meters"`
}

type CourierWithPaginate[T any] struct {
	CurrentPage int   `json:"current_page"`
	Data        []T   `json:"data"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
}

type SearchCourier struct {
	Name      string
	Email     string
	Longitude float64
	Latitude  float64
	Radius    float64
	Page      int
	PerPage   int
	SortBy    string
	OrderBy   string
}
