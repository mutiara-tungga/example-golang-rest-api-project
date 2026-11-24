package httpmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	pkgErr "golang-rest-api/pkg/error"
	httpserver "golang-rest-api/pkg/http_server"
	"golang-rest-api/pkg/jwt"
	"golang-rest-api/pkg/log"
)

type claimsKey struct{}

var errUnauthorized = pkgErr.NewCustomError("unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)

func JWTAuth(parser jwt.JWTParser, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ""
			if cookieName != "" {
				cookie, _ := r.Cookie(cookieName)
				if cookie != nil {
					token = cookie.Value
				}
			}
			if token == "" {
				parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					token = parts[1]
				}
			}
			if token == "" {
				log.Error(r.Context(), "token not found", fmt.Errorf("token not found"))
				writeUnauthorized(w, errUnauthorized)
				return
			}
			claims, err := parser.ParseAndValidate(r.Context(), token)
			if err != nil {
				writeUnauthorized(w, pkgErr.NewCustomErrWithOriginalErr(errUnauthorized, err))
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Claims(ctx context.Context) (jwt.JWTClaims, error) {
	claims, ok := ctx.Value(claimsKey{}).(jwt.JWTClaims)
	if !ok || claims.Subject == "" {
		return jwt.JWTClaims{}, errUnauthorized
	}
	return claims, nil
}

func writeUnauthorized(w http.ResponseWriter, err error) {
	httpserver.WriteError(w, err)
}
