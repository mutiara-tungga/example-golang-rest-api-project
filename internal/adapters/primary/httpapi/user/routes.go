package httpapi

import (
	"net/http"

	httpmiddleware "golang-rest-api/pkg/http_middleware"
	"golang-rest-api/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(users UserHandler, parser jwt.JWTParser) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	r.Post("/api/v1/user/login", users.Login)
	r.Group(func(r chi.Router) {
		r.Use(httpmiddleware.JWTAuth(parser, AccessTokenCookieName))
		r.Post("/api/v1/user", users.CreateUser)
		r.Get("/api/v1/user/profile", users.Profile)
	})
	return r
}
