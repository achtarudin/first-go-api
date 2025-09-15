package entity

type Merchant struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

type UserMerchant struct {
	ID        int        `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	Token     string     `json:"token,omitempty"`
	Merchants []Merchant `json:"merchants,omitempty"`
}
