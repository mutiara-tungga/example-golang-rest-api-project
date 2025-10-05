package user

import (
	userMock "golang-rest-api/internal/service/user/mock"
	httpserver "golang-rest-api/pkg/http_server"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_UserHandler_Login(t *testing.T) {
	type expectedRes struct {
		httpStatus int
		respBody   string
	}
	testCases := []struct {
		name        string
		mockFunc    func(userSrv *userMock.MockIUserService)
		req         string
		expectedRes expectedRes
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := userMock.NewMockIUserService(ctrl)

			tc.mockFunc(userSvc)

			hndlr := NewUserHandler(userSvc)
			r := chi.NewRouter()
			r.Method(http.MethodPost, "/api/v1/user/login", httpserver.HandlerWithError(hndlr.Login))

			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", strings.NewReader(tc.req))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			raw := w.Result()
			defer raw.Body.Close()
			body, err := io.ReadAll(raw.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.expectedRes.httpStatus, w.Result().StatusCode)
			assert.JSONEq(t, tc.expectedRes.respBody, string(body))
		})
	}
}
