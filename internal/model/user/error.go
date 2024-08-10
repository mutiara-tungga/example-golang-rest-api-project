package user

import (
	pkgErr "golang-rest-api/pkg/error"
	"net/http"
)

var (
	ErrorDuplicateUsername = pkgErr.NewCustomError("error execute query", "DUPLICATE_USERNAME", http.StatusBadRequest)
)
