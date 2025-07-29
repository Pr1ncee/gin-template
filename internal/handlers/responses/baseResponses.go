/*
Package responses provides structs that represent responses for different use cases.

Specifically, this file provides a base response that can be used in any endpoints
and includes error, success and pagination.
*/
package responses

// ErrorResponse represents error response with corresponding error in a string format, message and status code.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents successful response of an endpoint.
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// PaginationResponse represents a response with pagination parameters
// such as current data and page, limit of the page and total count of all objects.
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"total_count,omitempty"`
}
