package httpapi

import (
	"errors"
	"net/http"

	service "golang-rest-api/internal/service/user"
	pkgErr "golang-rest-api/pkg/error"
	httpserver "golang-rest-api/pkg/http_server"
)

func writeSuccess(w http.ResponseWriter, status int, message string, data any) error {
	return httpserver.WriteSuccess(w, status, message, data)
}

func writeMessage(w http.ResponseWriter, status int, message string) error {
	return httpserver.WriteMessage(w, status, message)
}

func writeError(w http.ResponseWriter, err error) {
	httpserver.WriteError(w, err)
}

func decodeJSON(r *http.Request, dst any) error {
	return httpserver.DecodeJSON(r, dst)
}

func mapUserError(err error) error {
	switch {
	case errors.Is(err, service.ErrNotFound):
		return pkgErr.NewCustomErrWithOriginalErr(pkgErr.NewCustomError("error user not found", "USER_NOT_FOUND", http.StatusNotFound), err)
	case errors.Is(err, service.ErrDuplicateUsername):
		return pkgErr.NewCustomErrWithOriginalErr(pkgErr.NewCustomError("error duplicate username", "USER_ERROR_DUPLICATE_USERNAME", http.StatusBadRequest), err)
	case errors.Is(err, service.ErrWrongPassword):
		return pkgErr.NewCustomErrWithOriginalErr(pkgErr.NewCustomError("error password", "LOGIN_ERROR_WRONG_PASSWORD", http.StatusBadRequest), err)
	default:
		return err
	}
}
