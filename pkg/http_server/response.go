package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	pkgErr "golang-rest-api/pkg/error"
)

type SuccessResponse struct {
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Meta       any    `json:"meta,omitempty"`
	HTTPStatus int    `json:"http_status"`
}

type ErrorResponse struct {
	Message    string `json:"message,omitempty"`
	HTTPStatus int    `json:"http_status"`
	ErrorCode  string `json:"error_code"`
}

func WriteSuccess(w http.ResponseWriter, status int, message string, data any) error {
	return WriteJSON(w, status, SuccessResponse{HTTPStatus: status, Message: message, Data: data})
}

func WriteMessage(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, SuccessResponse{HTTPStatus: status, Message: message})
}

func WriteError(w http.ResponseWriter, err error) {
	resp := ErrorResponse{HTTPStatus: http.StatusInternalServerError, ErrorCode: "INTERNAL_SERVER_ERROR", Message: err.Error()}
	var custom pkgErr.CustomError
	if errors.As(err, &custom) {
		resp.ErrorCode, resp.HTTPStatus = custom.GetErrorCode(), custom.GetStatusCode()
	}
	_ = WriteJSON(w, resp.HTTPStatus, resp)
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

func DecodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
