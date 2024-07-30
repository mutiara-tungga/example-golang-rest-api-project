package httpserver

import (
	"context"
	"encoding/json"
	pkgErr "golang-rest-api/pkg/error"
	"golang-rest-api/pkg/log"
	"net/http"
)

type HttpSuccessResponse struct {
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Meta       any    `json:"meta,omitempty"`
	HTTPStatus int    `json:"http_status"`
}

type HttpErrorResponse struct {
	Message    string `json:"message,omitempty"`
	HTTPStatus int    `json:"http_status"`
	ErrorCode  string `json:"error_code"`
}

type PaginationMetaInfo struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

const (
	HeaderKeyContentType  = "Content-Type"
	HeaderApplicationJson = "application/json"
)

func WriteJsonMsgWithData(ctx context.Context, w http.ResponseWriter, statusCode int, msg string, data any) {
	hr := HttpSuccessResponse{
		HTTPStatus: statusCode,
		Message:    msg,
		Data:       data,
	}

	w.Header().Set(HeaderKeyContentType, HeaderApplicationJson)
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(hr)
	if err != nil {
		log.Error(ctx, "error http write reponse", err)
	}
}

func WriteJsonError(ctx context.Context, w http.ResponseWriter, err error) {
	hr := HttpErrorResponse{
		HTTPStatus: http.StatusInternalServerError,
		ErrorCode:  "INTERNAL_SERVER_ERROR",
		Message:    err.Error(),
	}

	if err, ok := err.(*pkgErr.CustomError); ok {
		hr.ErrorCode = err.GetErrorCode()
		hr.HTTPStatus = err.GetStatusCode()
	}

	w.Header().Set(HeaderKeyContentType, HeaderApplicationJson)
	w.WriteHeader(hr.HTTPStatus)

	err = json.NewEncoder(w).Encode(hr)
	if err != nil {
		log.Error(ctx, "error http write reponse", err)
	}
}
