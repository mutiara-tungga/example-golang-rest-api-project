package httpapi

import (
	"fmt"
	"net/http"

	service "golang-rest-api/internal/service/user"
	pkgErr "golang-rest-api/pkg/error"
	httpserver "golang-rest-api/pkg/http_server"
	"golang-rest-api/pkg/validator"
)

type SuccessResponse = httpserver.SuccessResponse
type ErrorResponse = httpserver.ErrorResponse

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login godoc
// @Summary Login
// @Description Authenticate a user and set access and refresh token cookies.
// @Tags user
// @Accept json
// @Produce json
// @Param request body loginRequest true "Request body"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/user/login [post]
func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, mapUserError(pkgErr.NewCustomErrWithOriginalErr(errorInvalidJSON, err)))
		return
	}
	if err := validator.Validate.StructCtx(ctx, req); err != nil {
		writeError(w, mapUserError(pkgErr.NewCustomError(fmt.Sprintf("payload not valid: %s", err), "PAYLOAD_NOT_VALID", http.StatusBadRequest)))
		return
	}
	tokens, err := h.useCases.Login(ctx, service.LoginCommand{Username: req.Username, Password: req.Password})
	if err != nil {
		writeError(w, mapUserError(err))
		return
	}

	http.SetCookie(w, &http.Cookie{Name: AccessTokenCookieName, Value: tokens.AccessToken, Expires: tokens.ExpiresAt, HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: RefreshTokenCookieName, Value: tokens.RefreshToken, Expires: tokens.RefreshTokenExpiresAt, HttpOnly: true})
	if err := writeMessage(w, http.StatusOK, "login success"); err != nil {
		writeError(w, err)
	}
}
