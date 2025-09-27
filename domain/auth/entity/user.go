package entity

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

func NewUser(id int, name, email, password string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
}
