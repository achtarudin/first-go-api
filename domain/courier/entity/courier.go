package entity

type Courier struct {
	ID        int     `json:"id"`
	RoleID    int     `json:"role_id,omitempty"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"password,omitempty"`
	Token     string  `json:"token,omitempty"`
	Phone     string  `json:"phone"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}
