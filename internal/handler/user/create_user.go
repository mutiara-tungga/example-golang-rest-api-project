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

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := modelUser.CreateUserReq{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error(ctx, "error decode json", err)
		return pkgErr.NewCustomErrWithOriginalErr(model.ErrorInvalidJson, err)
	}

	resp, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		return err
	}

	httpserver.WriteJsonMsgWithData(ctx, w, http.StatusCreated, "user created", resp)
	return nil
}
