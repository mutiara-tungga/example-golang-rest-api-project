package httpmiddleware

import (
	"context"
	"fmt"
	pkgErr "golang-rest-api/pkg/error"
	httpserver "golang-rest-api/pkg/http_server"
	"golang-rest-api/pkg/jwt"
	"golang-rest-api/pkg/log"
	"net/http"
	"strings"
)

var (
	ErrorUnauthorized      = pkgErr.NewCustomError("unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
	ErrorJWTClaimsNotFound = pkgErr.NewCustomError("unauthorized: user jwt claims not found", "JWT_CLAIMS_NOT_FOUND", http.StatusUnauthorized)
)

type contexKey string

const (
	contextKeyUserClaims contexKey = "user_claims"
)

func JWTAuthUser(validator jwt.JWTValidator, cookieName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				tokenString, err := getToken(r, cookieName)
				if err != nil {
					httpserver.WriteJsonError(ctx, w, err)
					return
				}

				tokenClaims, err := validator.Validate(ctx, tokenString)
				if err != nil {
					httpserver.WriteJsonError(ctx, w, pkgErr.NewCustomErrWithOriginalErr(ErrorUnauthorized, err))
					return
				}

				r = r.WithContext(context.WithValue(ctx, contextKeyUserClaims, tokenClaims))
				next.ServeHTTP(w, r)
			})
	}
}

func getToken(r *http.Request, cookieName string) (string, error) {
	ctx := r.Context()

	token := ""
	if len(cookieName) > 0 {
		cookie, _ := r.Cookie(cookieName)
		if cookie != nil {
			token = cookie.Value
		}
	}

	if len(token) == 0 {
		bearerToken := r.Header.Get("Authorization")
		bearerTokens := strings.Split(bearerToken, " ")
		if len(bearerTokens) > 1 {
			token = bearerTokens[1]
		}
	}

	if len(token) == 0 {
		err := fmt.Errorf("token not found")
		log.Error(ctx, "token not found", err)
		return "", pkgErr.NewCustomErrWithOriginalErr(ErrorUnauthorized, err)
	}

	return token, nil
}

func GetUserClaims(ctx context.Context) (jwt.JWTClaims, error) {
	v := ctx.Value(contextKeyUserClaims)

	jwtClaims, ok := v.(jwt.JWTClaims)
	if !ok {
		return jwt.JWTClaims{}, ErrorJWTClaimsNotFound
	}

	if len(jwtClaims.Subject) == 0 {
		return jwt.JWTClaims{}, ErrorJWTClaimsNotFound
	}

	return jwtClaims, nil
}
