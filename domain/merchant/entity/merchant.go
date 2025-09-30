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

type UserMerchantLogin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserMerchantRegister struct {
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Password string   `json:"password,omitempty"`
	Merchant Merchant `json:"merchant,omitempty"`
}

type AddMerchant struct {
	Merchant Merchant `json:"merchant,omitempty"`
}
