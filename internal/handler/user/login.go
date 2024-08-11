package user

import (
	"encoding/json"
	"golang-rest-api/internal/model"
	modelUser "golang-rest-api/internal/model/user"
	pkgErr "golang-rest-api/pkg/error"
	httpserver "golang-rest-api/pkg/http_server"
	"golang-rest-api/pkg/log"
	"net/http"
)

func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := modelUser.UserLoginReq{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error(ctx, "error decode json", err)
		return pkgErr.NewCustomErrWithOriginalErr(model.ErrorInvalidJson, err)
	}

	jwtToken, err := h.userService.UserLogin(ctx, req)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     modelUser.AccessTokenCookieName,
		Value:    jwtToken.AccessToken,
		Expires:  jwtToken.ExpiresAt,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     modelUser.RefreshTokenCookieName,
		Value:    jwtToken.RefreshToken,
		Expires:  jwtToken.RefreshTokenExpiresAt,
		HttpOnly: true,
	})

	httpserver.WriteJsonMsgOnly(ctx, w, http.StatusOK, "login success")
	return nil
}
