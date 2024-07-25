package error

import (
	"fmt"
)

type CustomError struct {
	message       string
	errorCode     string
	statusCode    int
	originalError error
}

func (err CustomError) Error() string {
	return err.message
}

func (err CustomError) OriginalError() error {
	return err.originalError
}

func NewCustomError(
	message string,
	errorCode string,
	statusCode int,
) CustomError {
	return CustomError{
		message:       message,
		statusCode:    statusCode,
		errorCode:     errorCode,
		originalError: fmt.Errorf("error: %s, code: %s, status code: %d", message, errorCode, statusCode),
	}
}

func NewCustomErrWithOriginalErr(
	customErr CustomError,
	originalErr error,
) CustomError {
	return CustomError{
		message:       customErr.message,
		statusCode:    customErr.statusCode,
		errorCode:     customErr.errorCode,
		originalError: originalErr,
	}
}
