package response

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Message success"`
	Data    any    `json:"data,omitempty"`
}
