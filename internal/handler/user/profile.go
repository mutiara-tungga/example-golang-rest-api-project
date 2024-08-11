package user

import (
	httpmiddleware "golang-rest-api/pkg/http_middleware"
	httpserver "golang-rest-api/pkg/http_server"
	"net/http"
)

func (h UserHandler) UserProfile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	u, err := httpmiddleware.GetUserClaims(ctx)
	if err != nil {
		return err
	}

	resp, err := h.userService.UserProfile(ctx, u.Subject)
	if err != nil {
		return err
	}

	httpserver.WriteJsonMsgWithData(ctx, w, http.StatusOK, "success get progile", resp)
	return nil
}
