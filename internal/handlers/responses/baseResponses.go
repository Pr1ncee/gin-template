package responses

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"total_count,omitempty"`
}
