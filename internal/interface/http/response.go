package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error codes
const (
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeInvalidEmail    = "INVALID_EMAIL"
	ErrCodeInvalidPassword = "INVALID_PASSWORD"
)

// APIResponse represents a standard successful API response
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total      int `json:"total" example:"100"`
	Page       int `json:"page" example:"1"`
	PageSize   int `json:"page_size" example:"20"`
	TotalPages int `json:"total_pages" example:"5"`
}

// ListData represents data structure for list responses
type ListData struct {
	Items interface{}    `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

// APIListResponse represents a standard successful list API response with pagination
type APIListResponse struct {
	Success bool     `json:"success" example:"true"`
	Data    ListData `json:"data"`
}

// APIError represents error details
type APIError struct {
	Code    string `json:"code" example:"VALIDATION_ERROR"`
	Message string `json:"message" example:"Invalid input data"`
}

// APIErrorResponse represents a standard error API response
type APIErrorResponse struct {
	Success bool     `json:"success" example:"false"`
	Error   APIError `json:"error"`
}

// SendSuccess sends a successful response with data and optional message
func SendSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// SendSuccessList sends a successful list response with pagination metadata
func SendSuccessList(c *gin.Context, items interface{}, meta PaginationMeta) {
	c.JSON(http.StatusOK, APIListResponse{
		Success: true,
		Data: ListData{
			Items: items,
			Meta:  meta,
		},
	})
}

// SendError sends an error response with error code and message
func SendError(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, APIErrorResponse{
		Success: false,
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

// SendInternalError logs the error details server-side and sends a generic 500 error to client
func SendInternalError(c *gin.Context, err error) {
	// Log detailed error information server-side
	log.Printf("[ERROR] Internal server error: %v", err)

	// Send generic error message to client
	c.JSON(http.StatusInternalServerError, APIErrorResponse{
		Success: false,
		Error: APIError{
			Code:    ErrCodeInternal,
			Message: "Something went wrong",
		},
	})
}
