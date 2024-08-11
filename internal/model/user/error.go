package user

import (
	pkgErr "golang-rest-api/pkg/error"
	"net/http"
)

var (
	ErrorDuplicateUsername       = pkgErr.NewCustomError("error duplicate username", "USER_ERROR_DUPLICATE_USERNAME", http.StatusBadRequest)
	ErrorUserNotFound            = pkgErr.NewCustomError("error user not found", "USER_NOT_FOUND", http.StatusNotFound)
	ErrorLoginErrorWrongPassword = pkgErr.NewCustomError("error password", "LOGIN_ERROR_WRONG_PASSWORD", http.StatusBadRequest)
)
