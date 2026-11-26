package httpapi

import (
	"net/http"

	httpmiddleware "golang-rest-api/pkg/http_middleware"
)

// Profile godoc
// @Summary Get current user profile
// @Description Return the profile for the user identified by the access token cookie.
// @Tags user
// @Produce json
// @Security CookieAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/user/profile [get]
func (h UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	claims, err := httpmiddleware.Claims(r.Context())
	if err != nil {
		writeError(w, mapUserError(err))
		return
	}

	profile, err := h.useCases.Profile(r.Context(), claims.Subject)
	if err != nil {
		writeError(w, mapUserError(err))
		return
	}
	if err := writeSuccess(w, http.StatusOK, "success get profile", profile); err != nil {
		writeError(w, err)
	}
}
