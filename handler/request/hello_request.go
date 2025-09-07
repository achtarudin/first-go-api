package request

type HelloRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Age   int    `json:"age" example:"27"`
	Email string `json:"email" example:"john.doe@example.com"`
} //@name HelloRequest
