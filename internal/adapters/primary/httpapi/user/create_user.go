package httpapi

import (
	"net/http"

	service "golang-rest-api/internal/service/user"
	pkgErr "golang-rest-api/pkg/error"
)

type createUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Actor    string `json:"-"`
}

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, mapUserError(pkgErr.NewCustomErrWithOriginalErr(errorInvalidJSON, err)))
		return
	}

	result, err := h.useCases.CreateUser(r.Context(), service.CreateUserCommand{
		Name: req.Name, Username: req.Username, Phone: req.Phone, Password: req.Password, Actor: req.Actor,
	})
	if err != nil {
		writeError(w, mapUserError(err))
		return
	}
	if err := writeSuccess(w, http.StatusCreated, "user created", result); err != nil {
		writeError(w, err)
	}
}
