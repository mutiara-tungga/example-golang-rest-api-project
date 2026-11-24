package httpapi

import (
	"context"
	"net/http"

	service "golang-rest-api/internal/service/user"
	pkgErr "golang-rest-api/pkg/error"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

var errorInvalidJSON = pkgErr.NewCustomError("invalid json payload", "ERROR_INVALID_JSON_PAYLOAD", http.StatusBadRequest)

type userUseCases interface {
	CreateUser(ctx context.Context, cmd service.CreateUserCommand) (service.CreateUserResult, error)
	Login(ctx context.Context, cmd service.LoginCommand) (service.TokenPair, error)
	Profile(ctx context.Context, id string) (service.Profile, error)
}

type UserHandler struct{ useCases userUseCases }

func NewUserHandler(useCases userUseCases) UserHandler { return UserHandler{useCases: useCases} }
