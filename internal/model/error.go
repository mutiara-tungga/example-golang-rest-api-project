package model

import (
	pkgErr "golang-rest-api/pkg/error"
)

var (
	ErrorExecQuery = pkgErr.NewCustomError("error execute query", "ERROR_EXECUTE_QUERY", 500)
)
