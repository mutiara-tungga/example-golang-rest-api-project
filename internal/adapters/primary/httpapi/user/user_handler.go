package httpapi

import (
	"net/http"

	service "golang-rest-api/internal/service/user"
	pkgErr "golang-rest-api/pkg/error"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

var errorInvalidJSON = pkgErr.NewCustomError("invalid json payload", "ERROR_INVALID_JSON_PAYLOAD", http.StatusBadRequest)

type UserHandler struct{ useCases service.UserUseCases }

func NewUserHandler(useCases service.UserUseCases) UserHandler {
	return UserHandler{useCases: useCases}
}
