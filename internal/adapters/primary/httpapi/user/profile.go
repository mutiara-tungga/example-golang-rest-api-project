package httpapi

import (
	"net/http"

	httpmiddleware "golang-rest-api/pkg/http_middleware"
)

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
