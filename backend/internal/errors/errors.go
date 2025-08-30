package errors

import "fmt"

type APIError struct {
	StatusCode int `json:"status_code"`
	Message    any `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error with status code %d: %s", e.StatusCode, e.Message)
}

func NewApiError(statusCode int, message any) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}
