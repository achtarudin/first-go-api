package response

// BindErrorResponse represents a standard error response for binding errors
type BindErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad Request"`
	Errors  any    `json:"errors,omitempty" example:"{\"field\":\"error message\"}"`
} //@name BindErrorResponse
