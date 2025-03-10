package user

import (
	"encoding/json"
	"fmt"
	"golang-rest-api/internal/model"
	modelUser "golang-rest-api/internal/model/user"
	pkgErr "golang-rest-api/pkg/error"
	httpserver "golang-rest-api/pkg/http_server"
	"golang-rest-api/pkg/log"
	"golang-rest-api/pkg/validator"
	"net/http"
)

// Login godoc
// @Summary      Login
// @Description  Login
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body modelUser.UserLoginReq true "Request Body"
// @Success      200  {object}  httpserver.HttpSuccessResponse
// @Failure      400  {object}  httpserver.HttpErrorResponse
// @Failure      404  {object}  httpserver.HttpErrorResponse
// @Failure      500  {object}  httpserver.HttpErrorResponse
// @Router       /api/v1/user/login [post]
func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := modelUser.UserLoginReq{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error(ctx, "error decode json", err)
		return pkgErr.NewCustomErrWithOriginalErr(model.ErrorInvalidJson, err)
	}

	err = validator.Validate.StructCtx(ctx, req)
	if err != nil {
		return pkgErr.NewCustomError(fmt.Sprintf("payload not valid: %s", err.Error()), "PAYLOAD_NOT_VALID", http.StatusBadRequest)
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
