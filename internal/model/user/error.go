package user

import (
	pkgErr "golang-rest-api/pkg/error"
	"net/http"
)

var (
	ErrorDuplicateUsername = pkgErr.NewCustomError("error duplicate username", "DUPLICATE_USERNAME", http.StatusBadRequest)
)
