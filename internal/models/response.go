package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIErrorCode string

const (
	ErrDatabase       APIErrorCode = "DATABASE_ERROR"
	ErrRecordNotFound APIErrorCode = "RECORD_NOT_FOUND"

	ErrUnauthorized APIErrorCode = "UNAUTHORIZED"
	ErrInvalidToken APIErrorCode = "INVALID_TOKEN"

	ErrValidation   APIErrorCode = "VALIDATION_ERROR"
	ErrInvalidInput APIErrorCode = "INVALID_INPUT"

	ErrFileNotFound APIErrorCode = "FILE_NOT_FOUND"
	ErrInvalidFile  APIErrorCode = "INVALID_FILE_TYPE"
	ErrFileTooLarge APIErrorCode = "FILE_TOO_LARGE"
)

type APIResponse struct {
	Status  int       `json:"status"`
	Data    any       `json:"data,omitempty"`
	Message string    `json:"message,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    APIErrorCode `json:"code"`
	Message string       `json:"message"`
	Details string       `json:"details,omitempty"`
}

func SuccessResponse(data any) APIResponse {
	return APIResponse{
		Status: http.StatusOK,
		Data:   data,
	}
}

func ErrorResponse(status int, code APIErrorCode, message string) APIResponse {
	return APIResponse{
		Status: status,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}

func (r APIResponse) Send(ctx *gin.Context) {
	ctx.JSON(r.Status, r)
}
